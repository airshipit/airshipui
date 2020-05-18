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
	"bytes"
	"fmt"
	"text/template"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipui/internal/configs"
)

// configPage struct is used for templated HTML
type configPage struct {
	ClusterRows    string
	ContextRows    string
	CredentialRows string
	Title          string
	Version        string
}

// GetCluster gets cluster information from the airshipctl config
func (c *client) GetCluster() []*config.Cluster {
	return c.settings.Config.GetClusters()
}

// getClusterTableRows turns an array of cluster into html table rows
func getClusterTableRows() string {
	info := c.GetCluster()

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
func (c *client) GetContext() []*config.Context {
	return c.settings.Config.GetContexts()
}

// getContextTableRows turns an array of contexts into html table rows
func getContextTableRows() string {
	info := c.GetContext()

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
func (c *client) GetCredential() []*config.AuthInfo {
	authinfo, err := c.settings.Config.GetAuthInfos()
	if err != nil {
		return []*config.AuthInfo{}
	}

	return authinfo
}

// getContextTableRows turns an array of contexts into html table rows
func getCredentialTableRows() string {
	info := c.GetCredential()

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

func getDefaultHTML() (string, error) {
	// go templates need an io writer, since we need a string this buffer can be converted
	var buff bytes.Buffer

	// TODO: make the node path dynamic or setable at compile time
	t, err := template.ParseFiles("./internal/integrations/ctl/templates/config.html")

	if err != nil {
		return "", err
	}

	// add contents to the page
	p := configPage{
		ClusterRows:    getClusterTableRows(),
		ContextRows:    getContextTableRows(),
		CredentialRows: getCredentialTableRows(),
		Title:          "Config",
		Version:        GetAirshipCTLVersion(),
	}

	// parse and merge the template
	err = template.Must(t, err).Execute(&buff, p)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

// SetConfig will flop between requests so we don't have to have them all mapped as function calls
func SetConfig(request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:      configs.AirshipCTL,
		Component: configs.SetConfig,
	}

	var err error
	var message string
	switch request.SubComponent {
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
