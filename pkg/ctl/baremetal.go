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
	"strings"

	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/statistics"
	"opendev.org/airship/airshipui/pkg/webservice"

	"opendev.org/airship/airshipctl/pkg/remote"
)

type nodeInfo struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	BMCAddress string `json:"bmcAddress,omitempty"`
}

type phaseInfo struct {
	Name         string `json:"name,omitempty"`
	GenerateName string `json:"generateName,omitempty"`
	Namespace    string `json:"namespace,omitempty"`
	ClusterName  string `json:"clusterName,omitempty"`
}

type defaultData struct {
	Nodes  []nodeInfo  `json:"nodes,omitempty"`
	Phases []phaseInfo `json:"phases,omitempty"`
}

// HandleBaremetalRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandleBaremetalRequest(user *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Baremetal,
		SubComponent: request.SubComponent,
	}

	var err error
	var message *string

	subComponent := request.SubComponent

	if request.Targets != nil {
		s := fmt.Sprintf("%s action has been requested on hosts: %s", subComponent, strings.Join(*request.Targets, ", "))
		message = &s
	}

	switch subComponent {
	case configs.GetDefaults:
		response.Data, err = getDefaults(request)
	case configs.EjectMedia:
		err = doAction(user, request)
	case configs.PowerOff:
		err = doAction(user, request)
	case configs.PowerOn:
		err = doAction(user, request)
	case configs.PowerStatus:
		err = fmt.Errorf("Subcomponent %s not implemented", subComponent)
	case configs.Reboot:
		err = doAction(user, request)
	case configs.RemoteDirect:
		err = doAction(user, request)
	default:
		err = fmt.Errorf("Subcomponent %s not found", subComponent)
	}

	if err != nil {
		e := err.Error()
		response.Error = &e
	} else {
		response.Message = message
	}

	return response
}

func getDefaults(request configs.WsMessage) (defaultData, error) {
	nodeInfo, err := getNodeInfo(request)
	phaseInfo, err2 := getPhaseInfo()

	if err != nil && err2 != nil {
		err = fmt.Errorf("Node error: %v.  Phase error %v", err, err2)
	} else if err2 != nil {
		err = err2
	}

	return defaultData{
		Nodes:  nodeInfo,
		Phases: phaseInfo,
	}, err
}

// getNodeInfo gets and formats the default nodes as defined by the manifest(s)
func getNodeInfo(request configs.WsMessage) ([]nodeInfo, error) {
	client, err := NewClient(AirshipConfigPath, KubeConfigPath, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	selectors := []remote.HostSelector{remote.All()}
	// bootstrap is the default "phase" this may change as it does not accept an empty string as a default
	m, err := remote.NewManager(client.Config, "bootstrap", selectors...)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	data := []nodeInfo{}

	for _, host := range m.Hosts {
		data = append(data, nodeInfo{
			Name:       host.HostName,
			ID:         host.NodeID(),
			BMCAddress: host.BMCAddress,
		})
	}
	return data, nil
}

// getPhaseInfo gets and formats the phases as defined by the manifest(s)
func getPhaseInfo() ([]phaseInfo, error) {
	helper, err := getHelper()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	phases, err := helper.ListPhases()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	data := []phaseInfo{}
	for _, p := range phases {
		data = append(data, phaseInfo{
			Name:         p.Name,
			GenerateName: p.GenerateName,
			Namespace:    p.Namespace,
			ClusterName:  p.ClusterName,
		})
	}

	return data, nil
}

func doAction(user *string, request configs.WsMessage) error {
	actionType := request.ActionType
	if request.Targets == nil && actionType == nil {
		err := errors.New("No target nodes or phases defined.  Cannot proceed with request")
		return err
	}

	defaultPhase := "bootstrap"
	if request.Targets != nil {
		for _, target := range *request.Targets {
			if *actionType == configs.DirectAction {
				go actionHelper(user, target, defaultPhase, request)
			} else {
				go actionHelper(user, "", target, request)
			}
		}
	}

	return nil
}

func actionHelper(user *string, target string, phase string, request configs.WsMessage) {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Baremetal,
		SubComponent: configs.EjectMedia,
		SessionID:    request.SessionID,
		ActionType:   request.ActionType,
		Target:       &target,
	}

	// create a transaction for this singular request
	transaction := statistics.NewTransaction(user, response)

	client, err := NewClient(AirshipConfigPath, KubeConfigPath, response)
	if err != nil {
		log.Error(err)
		e := err.Error()
		response.Error = &e
		transaction.Complete(false)
		err = webservice.WebSocketSend(response)
		if err != nil {
			log.Error(err)
		}
		return
	}

	var selectors []remote.HostSelector
	if len(target) != 0 {
		selectors = []remote.HostSelector{remote.ByName(target)}
	} else {
		selectors = []remote.HostSelector{remote.All()}
	}
	m, err := remote.NewManager(client.Config, phase, selectors...)
	if err != nil {
		log.Error(err)
		e := err.Error()
		response.Error = &e
		transaction.Complete(false)
		err = webservice.WebSocketSend(response)
		if err != nil {
			log.Error(err)
		}
		return
	}

	action := request.SubComponent
	if len(m.Hosts) != 1 {
		e := fmt.Sprintf("More than one node found cannot complete %s on %s", action, target)
		log.Error(&e)
		response.Error = &e
		transaction.Complete(false)
		err = webservice.WebSocketSend(response)
		if err != nil {
			log.Error(err)
		}
		return
	}

	host := m.Hosts[0]
	switch action {
	case configs.EjectMedia:
		err = host.EjectVirtualMedia(host.Context)
	case configs.PowerOff:
		err = host.SystemPowerOff(host.Context)
	case configs.PowerOn:
		err = host.SystemPowerOn(host.Context)
	case configs.Reboot:
		err = host.RebootSystem(host.Context)
	}

	if err != nil {
		log.Error(err)
		e := err.Error()
		response.Error = &e
		transaction.Complete(false)
		err = webservice.WebSocketSend(response)
		if err != nil {
			log.Error(err)
		}
		return
	}

	s := fmt.Sprintf("%s on %s completed successfully", action, target)
	response.Message = &s
	transaction.Complete(true)
	err = webservice.WebSocketSend(response)
	if err != nil {
		log.Error(err)
	}
}
