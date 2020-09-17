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
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/phase"
	"opendev.org/airship/airshipui/pkg/log"
	"sigs.k8s.io/kustomize/api/types"
)

var phaseIndex map[string]PhaseObj = buildPhaseIndex()

// PhaseObj lightweight structure to hold the name and document
// entrypoint for airshipctl phases
type PhaseObj struct {
	Group      string
	Name       string
	Entrypoint string
}

func buildPhaseIndex() map[string]PhaseObj {
	client := NewDefaultClient()
	return client.buildPhaseIndex()
}

func (client *Client) buildPhaseIndex() map[string]PhaseObj {
	idx := map[string]PhaseObj{}

	// get target path from ctl settings
	tp, err := client.settings.Config.CurrentContextTargetPath()
	if err != nil {
		log.Errorf("Error building phase index: %s", err)
		return nil
	}

	cmd := phase.Cmd{AirshipCTLSettings: client.settings}

	plan, err := cmd.Plan()
	if err != nil {
		log.Errorf("Error building phase index: %s", err)
		return nil
	}

	for grp, phases := range plan {
		for _, phase := range phases {
			p, err := cmd.GetPhase(phase)
			if err != nil {
				log.Errorf("Error building phase index: %s", err)
				return nil
			}

			entrypoint := fmt.Sprintf("%s/kustomization.yaml",
				filepath.Join(tp, p.Config.DocumentEntryPoint))

			idx[uuid.New().String()] = PhaseObj{
				Group:      grp,
				Name:       phase,
				Entrypoint: entrypoint,
			}
		}
	}
	return idx
}

// GetPhaseTree builds the initial structure of the phase tree
// consisting of phase Groups and Phases. Individual phase source
// files or rendered documents will be lazy loaded as needed
func (client *Client) GetPhaseTree() ([]KustomNode, error) {
	nodes := []KustomNode{}

	grpMap := map[string][]KustomNode{}
	for id, po := range phaseIndex {
		pNode := KustomNode{
			ID:          id,
			Name:        fmt.Sprintf("Phase: %s", po.Name),
			IsPhaseNode: true,
		}

		children, err := client.GetPhaseSourceFiles(id)
		if err != nil {
			// TODO(mfuller): push an error to UI so it can be handled by
			// toastr service, pending refactor of webservice and configs pkgs
			log.Errorf("Error building tree for phase '%s': %s", po.Name, err)
			pNode.HasError = true
		} else {
			pNode.Children = children
		}

		grpMap[po.Group] = append(grpMap[po.Group], pNode)
	}

	for name, phases := range grpMap {
		gNode := KustomNode{
			ID:       uuid.New().String(),
			Name:     fmt.Sprintf("Group: %s", name),
			Children: phases,
		}
		nodes = append(nodes, gNode)
	}

	return nodes, nil
}

// GetPhaseDocuments returns a slice of KustomNodes representing
// all of the rendered documents making up a phase bundle.
// Ordering is k8s Namespace -> k8s Kind -> document name
func GetPhaseDocuments(id string) ([]KustomNode, error) {
	if index == nil {
		index = map[string]interface{}{}
	}
	nsNodes := []KustomNode{}

	if p, ok := phaseIndex[id]; ok {
		// get map of all docs associated with this bundle
		docs, err := sortDocuments(p.Entrypoint)
		if err != nil {
			return nil, err
		}
		// namespace node
		for ns, kinds := range docs {
			nsNode := KustomNode{
				ID:       uuid.New().String(),
				Name:     ns,
				Children: []KustomNode{},
			}
			// kind node
			for kind, docs := range kinds {
				kNode := KustomNode{
					ID:       uuid.New().String(),
					Name:     kind,
					Children: []KustomNode{},
				}
				// doc node
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
			nsNodes = append(nsNodes, nsNode)
		}
	}
	return nsNodes, nil
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

// GetPhaseSourceFiles returns a slice of KustomNodes representing
// all of the directories that will be traversed when kustomize
// builds the document bundle. The tree hierarchy is:
// kustomize "type" (like function) -> directory name -> file name
func (client *Client) GetPhaseSourceFiles(id string) ([]KustomNode, error) {
	if index == nil {
		index = map[string]interface{}{}
	}
	dirNodes := []KustomNode{}

	if p, ok := phaseIndex[id]; ok {
		dirs, err := getKustomizeDirs(p.Entrypoint)
		if err != nil {
			return nil, err
		}

		dm, err := client.createDirsMap(dirs)
		if err != nil {
			return nil, err
		}

		// kustomize "type" node
		for t, data := range dm {
			tNode := KustomNode{
				ID:   uuid.New().String(),
				Name: t,
			}

			// directory node
			for _, d := range data {
				name := d[0]
				abs := d[1]
				dNode := KustomNode{
					ID:       uuid.New().String(),
					Name:     name,
					Children: []KustomNode{},
				}

				files, err := ioutil.ReadDir(abs)
				if err != nil {
					return nil, err
				}
				// file (leaf) node
				for _, f := range files {
					if !f.IsDir() {
						id := uuid.New().String()
						path := filepath.Join(abs, f.Name())
						dNode.Children = append(dNode.Children,
							KustomNode{
								ID:   id,
								Name: f.Name(),
							})
						index[id] = path
					}
				}
				tNode.Children = append(tNode.Children, dNode)
			}
			dirNodes = append(dirNodes, tNode)
		}
	}
	return dirNodes, nil
}

// KustomNode structure to represent the kustomization tree for a given phase
// bundle to be consumed by the UI frontend
type KustomNode struct {
	ID          string       `json:"id"`   // UUID for backend node index
	Name        string       `json:"name"` // name used for display purposes (cli, ui)
	IsPhaseNode bool         `json:"isPhaseNode"`
	HasError    bool         `json:"hasError"`
	Children    []KustomNode `json:"children"`
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
				log.Errorf("Error following kustomize tree: %s", err)
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

// helper function to group kustomize dirs by type (i.e. function, composite, etc)
func (client *Client) createDirsMap(dirs []string) (map[string][][]string, error) {
	dm := map[string][][]string{}

	tp, err := client.settings.Config.CurrentContextTargetPath()
	if err != nil {
		return nil, err
	}

	manifestsDir := filepath.Join(tp, "manifests")

	for _, d := range dirs {
		rel, err := filepath.Rel(manifestsDir, d)
		if err != nil {
			return nil, err
		}
		split := strings.SplitN(rel, string(os.PathSeparator), 2)
		dm[split[0]] = append(dm[split[0]], []string{split[1], d})
	}

	return dm, nil
}

func kustomLoader(kfile string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(kfile)
	if err != nil {
		// annoyingly, the actual kustomization file may be nested one
		// layer deeper than what's specified in its parent's kustomization
		// file. For example,
		//
		// resources:
		// - function/capm3
		//
		// may actually refer to function/capm3/v0.3.1/kustomization.yaml,
		// so we'll try drilling down one more level to find it
		dir := filepath.Dir(kfile)
		contents, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		if len(contents) == 0 || len(contents) > 1 || !os.FileInfo.IsDir(contents[0]) {
			return nil, fmt.Errorf("no kustomization file found at %s", dir)
		}

		kfile = filepath.Join(dir, contents[0].Name(), "kustomization.yaml")
		bytes, err = ioutil.ReadFile(kfile)
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

func makeResMap(kfile string) (map[string][]string, error) {
	resMap := map[string][]string{}

	bytes, err := kustomLoader(kfile)
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
