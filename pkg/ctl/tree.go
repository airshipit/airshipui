/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package ctl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/google/uuid"
	"opendev.org/airship/airshipctl/pkg/document"
	"sigs.k8s.io/kustomize/api/types"
)

var (
	// TODO: retrieve this dynamically from airship config
	manifestsDir = "/home/ubuntu/workspace/airshipctl/manifests"
)

// recursively collect all kustomization.yaml files starting from
// targetDir
// TODO: this will almost certainly go away when we use phase plans
func collectKustomizations(targetDir string) ([]string, error) {
	var kustomizations []string
	pattern := "kustomization.yaml"

	var walkFunc filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// only interested in files
		if info.IsDir() {
			return nil
		}
		if match, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if match {
			kustomizations = append(kustomizations, path)
		}
		return nil
	}

	err := filepath.Walk(targetDir, walkFunc)
	if err != nil {
		return nil, err
	}

	return kustomizations, nil
}

func MakeRenderedTree() ([]KustomNode, error) {
	index = map[string]interface{}{}

	bundles := []KustomNode{}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kusts, err := collectKustomizations(filepath.Join(home, targetPath))
	if err != nil {
		return nil, err
	}

	for i, k := range kusts {
		kn := KustomNode{
			ID:       uuid.New().String(),
			Name:     fmt.Sprintf("Phase %d", i),
			Data:     k,
			Children: []KustomNode{},
		}

		// get map of all docs associated with this bundle
		docs, err := sortDocuments(k)
		if err != nil {
			return nil, err
		}

		for ns, kinds := range docs {
			nsNode := KustomNode{
				ID:       uuid.New().String(),
				Name:     ns,
				Data:     "",
				Children: []KustomNode{},
			}

			for kind, docs := range kinds {
				kNode := KustomNode{
					ID:       uuid.New().String(),
					Name:     kind,
					Data:     "",
					Children: []KustomNode{},
				}

				for _, d := range docs {
					id := uuid.New().String()
					dNode := KustomNode{
						ID:   id,
						Name: d.GetName(),
					}

					index[id] = d

					kNode.Children = append(kNode.Children, dNode)
				}

				nsNode.Children = append(nsNode.Children, kNode)
			}

			kn.Children = append(kn.Children, nsNode)
		}

		bundles = append(bundles, kn)
	}

	return bundles, nil
}

// sort a bundle's docs into namespace, kind
func sortDocuments(path string) (map[string]map[string][]document.Document, error) {
	docMap := map[string]map[string][]document.Document{}

	bundle, err := document.NewBundleByPath(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	docs, err := bundle.GetAllDocuments()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		ns := doc.GetNamespace()
		if ns == "" {
			ns = "[no namespace]"
		}
		kind := doc.GetKind()

		if docMap[ns] == nil {
			docMap[ns] = map[string][]document.Document{}
		}

		docMap[ns][kind] = append(docMap[ns][kind], doc)
	}

	return docMap, nil
}

func MakeSourceTree() ([]KustomNode, error) {
	index = map[string]interface{}{}

	bundles := []KustomNode{}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kusts, err := collectKustomizations(filepath.Join(home, targetPath))
	if err != nil {
		return nil, err
	}

	for i, k := range kusts {
		kn := KustomNode{
			ID:       uuid.New().String(),
			Name:     fmt.Sprintf("Phase %d", i),
			Data:     k,
			Children: []KustomNode{},
		}

		dirs, err := getKustomizeDirs(k)
		if err != nil {
			return nil, err
		}

		for _, d := range dirs {
			name, err := filepath.Rel(manifestsDir, d)
			if err != nil {
				name = d
			}

			n := KustomNode{
				ID:       uuid.New().String(),
				Name:     name,
				Data:     d,
				Children: []KustomNode{},
			}

			files, err := ioutil.ReadDir(d)
			if err != nil {
				return nil, err
			}

			for _, f := range files {
				if !f.IsDir() {
					id := uuid.New().String()
					path := filepath.Join(d, f.Name())
					n.Children = append(n.Children,
						KustomNode{
							ID:   id,
							Name: f.Name(),
							Data: path,
						})
					index[id] = path
				}
			}

			kn.Children = append(kn.Children, n)
		}
		bundles = append(bundles, kn)
	}

	return bundles, nil
}

type KustomNode struct {
	ID       string       `json:"id"`   // UUID, maybe not necessary; mainly for UI
	Name     string       `json:"name"` // name used for display purposes (cli, ui)
	Data     string       `json:"data"` // this could be a Kustomization object, or a string containing a file path
	Children []KustomNode `json:"children"`
}

func contains(dirs []string, val string) bool {
	for _, d := range dirs {
		if d == val {
			return true
		}
	}
	return false
}

func appendDirs(dirs *[]string, subDirs []string) {
	for _, d := range subDirs {
		if !contains(*dirs, d) {
			*dirs = append(*dirs, d)
		}
	}
}

// returns a list of all directories encountered by following the
// kustomization tree
func getKustomizeDirs(entrypoint string) ([]string, error) {
	dirs := []string{}

	// add entrypoint dir first
	dirs = append(dirs, filepath.Dir(entrypoint))

	resMap, err := makeResMap(entrypoint)
	if err != nil {
		return nil, err
	}

	for _, sources := range resMap {
		for _, s := range sources {
			fi, err := os.Stat(s)
			if err != nil {
				log.Println(err)
				continue
			}
			if os.FileInfo.IsDir(fi) {
				if !contains(dirs, s) {
					dirs = append(dirs, s)
				}
				s = filepath.Join(s, "kustomization.yaml")
				subDirs, err := getKustomizeDirs(s)
				if err != nil {
					return nil, err
				}
				appendDirs(&dirs, subDirs)
			}
		}
	}

	sort.Strings(dirs)
	return dirs, nil
}

func makeResMap(kfile string) (map[string][]string, error) {
	resMap := map[string][]string{}

	bytes, err := ioutil.ReadFile(kfile)
	if err != nil {
		return nil, err
	}

	k := types.Kustomization{}
	err = k.Unmarshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("error processing '%s': %s", kfile, err)
	}

	basedir := filepath.Dir(kfile)

	if len(k.Resources) > 0 {
		for _, p := range k.Resources {
			path := filepath.Join(basedir, p)
			resMap["Resources"] = append(resMap["Resources"], path)
		}
	}

	if len(k.ConfigMapGenerator) > 0 {
		for _, p := range k.ConfigMapGenerator {
			for _, s := range p.FileSources {
				path := filepath.Join(basedir, s)
				resMap["ConfigMapGenerator"] = append(resMap["ConfigMapGenerator"], path)
			}
		}
	}

	if len(k.SecretGenerator) > 0 {
		for _, p := range k.SecretGenerator {
			for _, s := range p.FileSources {
				path := filepath.Join(basedir, s)
				resMap["SecretGenerator"] = append(resMap["SecretGenerator"], path)
			}
		}
	}

	if len(k.Generators) > 0 {
		for _, p := range k.Generators {
			path := filepath.Join(basedir, p)
			resMap["Generators"] = append(resMap["Generators"], path)
		}
	}

	if len(k.Transformers) > 0 {
		for _, p := range k.Transformers {
			path := filepath.Join(basedir, p)
			resMap["Transformers"] = append(resMap["Transformers"], path)
		}
	}

	return resMap, nil
}
