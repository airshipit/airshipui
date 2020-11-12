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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"opendev.org/airship/airshipctl/pkg/phase"
	"opendev.org/airship/airshipctl/pkg/phase/ifc"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"sigs.k8s.io/kustomize/api/types"
)

// TODO(mfuller): new helper each time? or once here?
var helper ifc.Helper

func getHelper() (ifc.Helper, error) {
	if helper != nil {
		return helper, nil
	}

	c, err := NewDefaultClient(configs.UIConfig.AirshipConfigPath)
	if err != nil {
		return nil, err
	}

	h, err := phase.NewHelper(c.Config)
	if err != nil {
		return nil, err
	}

	return h, nil
}

// GetPhaseTree builds the initial structure of the phase tree
// consisting of phase Groups and Phases. Individual phase source
// files or rendered documents will be lazy loaded as needed
func (client *Client) GetPhaseTree() ([]KustomNode, error) {
	nodes := []KustomNode{}

	helper, err := getHelper()
	if err != nil {
		return nil, err
	}

	phases, err := helper.ListPhases()
	if err != nil {
		return nil, err
	}

	for _, p := range phases {
		pNode := KustomNode{
			ID:          uuid.New().String(),
			PhaseID:     ifc.ID{Name: p.Name, Namespace: p.Namespace},
			Name:        fmt.Sprintf("Phase: %s", p.Name),
			IsPhaseNode: true,
			Children:    []KustomNode{},
		}

		// some phases don't have any associated documents, so don't look
		// for children unless a DocumentEntryPoint has been specified
		if p.Config.DocumentEntryPoint != "" {
			children, err := client.GetPhaseSourceFiles(pNode.PhaseID)
			if err != nil {
				// TODO(mfuller): push an error to UI so it can be handled by
				// toastr service, pending refactor of webservice and configs pkgs
				log.Errorf("Error building tree for phase '%s': %s", p.Name, err)
				pNode.HasError = true
			} else {
				pNode.Children = children
			}
		}
		nodes = append(nodes, pNode)
	}

	return nodes, nil
}

// GetPhaseSourceFiles returns a slice of KustomNodes representing
// all of the directories that will be traversed when kustomize
// builds the document bundle. The tree hierarchy is:
// kustomize "type" (like function) -> directory name -> file name
func (client *Client) GetPhaseSourceFiles(id ifc.ID) ([]KustomNode, error) {
	if fileIndex == nil {
		fileIndex = map[string]string{}
	}

	helper, err := getHelper()
	if err != nil {
		return nil, err
	}

	c := phase.NewClient(helper)

	phaseIfc, err := c.PhaseByID(id)
	if err != nil {
		return nil, err
	}

	docRoot, err := phaseIfc.DocumentRoot()
	if err != nil {
		return nil, err
	}

	dirs, err := getKustomizeDirs(filepath.Join(docRoot, "kustomization.yaml"))
	if err != nil {
		if errors.As(err, &phase.ErrDocumentEntrypointNotDefined{}) {
			return nil, nil
		}
		return nil, err
	}

	dm, err := client.createDirsMap(dirs)
	if err != nil {
		return nil, err
	}

	dirNodes := []KustomNode{}

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
					fileIndex[id] = path
				}
			}
			tNode.Children = append(tNode.Children, dNode)
		}
		dirNodes = append(dirNodes, tNode)
	}
	return dirNodes, nil
}

// KustomNode structure to represent the kustomization tree for a given phase
// bundle to be consumed by the UI frontend
type KustomNode struct {
	ID          string       `json:"id"` // UUID for backend node index
	PhaseID     ifc.ID       `json:"phaseId"`
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

	tp, err := client.Config.CurrentContextTargetPath()
	if err != nil {
		return nil, err
	}

	// these relative paths are making the BIG assumption that the target path
	// will point to the "workspace" directory which will have "airshipctl" and
	// "treasuremap" as subdirs at the same level, followed by a "manifests" dir
	// in each, followed by "functions", "composites", etc. in each manifests dir
	for _, d := range dirs {
		rel, err := filepath.Rel(tp, d)
		if err != nil {
			return nil, err
		}
		split := strings.SplitN(rel, string(os.PathSeparator), 3)
		dirGrp := filepath.Join(split[0], split[1])
		dm[dirGrp] = append(dm[dirGrp], []string{split[2], d})
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

	for _, p := range k.Resources {
		path := filepath.Join(basedir, p)
		resMap["Resources"] = append(resMap["Resources"], path)
	}

	for _, p := range k.ConfigMapGenerator {
		for _, s := range p.FileSources {
			path := filepath.Join(basedir, s)
			resMap["ConfigMapGenerator"] = append(resMap["ConfigMapGenerator"], path)
		}
	}

	for _, p := range k.SecretGenerator {
		for _, s := range p.FileSources {
			path := filepath.Join(basedir, s)
			resMap["SecretGenerator"] = append(resMap["SecretGenerator"], path)
		}
	}

	for _, p := range k.Generators {
		path := filepath.Join(basedir, p)
		resMap["Generators"] = append(resMap["Generators"], path)
	}

	for _, p := range k.Transformers {
		path := filepath.Join(basedir, p)
		resMap["Transformers"] = append(resMap["Transformers"], path)
	}

	return resMap, nil
}
