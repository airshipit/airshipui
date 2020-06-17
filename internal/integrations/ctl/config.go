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

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipui/internal/configs"
)

// HandleConfigRequest will flop between requests so we don't have to have them all mapped as function calls
func HandleConfigRequest(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.CTLConfig,
		SubComponent: request.SubComponent,
	}

	var err error
	var message string
	switch request.SubComponent {
	case configs.GetDefaults:
		response.HTML, err = getConfigHTML()
	case configs.SetContext:
		message, err = setContext(request)
	case configs.SetCluster:
		message, err = setCluster(request)
	case configs.SetCredential:
		message, err = setCredential(request)
	default:
		err = fmt.Errorf("Subcomponent %s not found", request.SubComponent)
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Message = message
	}

	return response
}

// GetCluster gets cluster information from the airshipctl config
func (c *client) getCluster() []*config.Cluster {
	return c.settings.Config.GetClusters()
}

// getClusterTableRows turns an array of cluster into html table rows
func getClusterTableRows() string {
	info := c.getCluster()

	var rows string
	for _, config := range info {
		// TODO: all rows are editable, probably shouldn't be
		rows += "<tr><td><div contenteditable=true>" +
			config.Bootstrap + "</div></td><td><div contenteditable=true>" +
			config.NameInKubeconf + "</div></td><td><div contenteditable=true>" +
			config.ManagementConfiguration + "</div></td><td>" +
			config.KubeCluster().LocationOfOrigin + "</td><td><div contenteditable=true>" +
			config.KubeCluster().Server + "</div></td><td><div contenteditable=true>" +
			config.KubeCluster().CertificateAuthority + "</div></td><td>" +
			"<button type=\"button\" class=\"btn btn-success\" onclick=\"saveConfig(this)\">Save</button></td></tr>"
	}
	return rows
}

// GetContext gets cluster information from the airshipctl config
func (c *client) getContext() []*config.Context {
	return c.settings.Config.GetContexts()
}

// getContextTableRows turns an array of contexts into html table rows
func getContextTableRows() string {
	info := c.getContext()

	var rows string
	for _, context := range info {
		// TODO: all rows are editable, probably shouldn't be
		rows += "<tr><td><div contenteditable=true>" +
			context.NameInKubeconf + "</div></td><td><div contenteditable=true>" +
			context.Manifest + "</div></td><td>" +
			context.KubeContext().LocationOfOrigin + "</td><td><div contenteditable=true>" +
			context.KubeContext().Cluster + "</div></td><td><div contenteditable=true>" +
			context.KubeContext().AuthInfo + "</div></td><td>" +
			"<button type=\"button\" class=\"btn btn-success\" onclick=\"saveConfig(this)\">Save</button></td></tr>"
	}
	return rows
}

// GetCredential gets user credentials from the airshipctl config
func (c *client) getCredential() []*config.AuthInfo {
	authinfo, err := c.settings.Config.GetAuthInfos()
	if err != nil {
		return []*config.AuthInfo{}
	}

	return authinfo
}

// getContextTableRows turns an array of contexts into html table rows
func getCredentialTableRows() string {
	info := c.getCredential()

	var rows string
	for _, credential := range info {
		// TODO: all rows are editable, probably shouldn't be
		rows += "<tr><td>" +
			credential.KubeAuthInfo().LocationOfOrigin + "</td><td><div contenteditable=true>" +
			credential.KubeAuthInfo().Username + "</div></td><td>" +
			"<button type=\"button\" class=\"btn btn-success\" onclick=\"saveConfig(this)\">Save</button></td></tr>"
	}
	return rows
}

func getConfigHTML() (string, error) {
	return getHTML("./internal/integrations/ctl/templates/config.html", ctlPage{
		ClusterRows:    getClusterTableRows(),
		ContextRows:    getContextTableRows(),
		CredentialRows: getCredentialTableRows(),
		Title:          "Config",
		Version:        getAirshipCTLVersion(),
	})
}

// SetCluster will take ui cluster info, translate them into CTL commands and send a response back to the UI
func setCluster(request configs.WsMessage) (string, error) {
	modified, err := config.RunSetCluster(&request.ClusterOptions, c.settings.Config, true)

	var message string
	if modified {
		message = fmt.Sprintf("Cluster %q of type %q modified.",
			request.ClusterOptions.Name, request.ClusterOptions.ClusterType)
	} else {
		message = fmt.Sprintf("Cluster %q of type %q created.",
			request.ClusterOptions.Name, request.ClusterOptions.ClusterType)
	}

	return message, err
}

// SetContext will take ui context info, translate them into CTL commands and send a response back to the UI
func setContext(request configs.WsMessage) (string, error) {
	modified, err := config.RunSetContext(&request.ContextOptions, c.settings.Config, true)

	var message string
	if modified {
		message = fmt.Sprintf("Context %q modified.", request.ClusterOptions.Name)
	} else {
		message = fmt.Sprintf("Context %q created.", request.ClusterOptions.Name)
	}

	return message, err
}

// SetContext will take ui context info, translate them into CTL commands and send a response back to the UI
func setCredential(request configs.WsMessage) (string, error) {
	modified, err := config.RunSetAuthInfo(&request.AuthInfoOptions, c.settings.Config, true)

	var message string
	if modified {
		message = fmt.Sprintf("Credential %q modified.", request.ClusterOptions.Name)
	} else {
		message = fmt.Sprintf("Credential %q created.", request.ClusterOptions.Name)
	}

	return message, err
}
