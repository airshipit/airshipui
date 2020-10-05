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
	"testing"

	"github.com/stretchr/testify/assert"
	"opendev.org/airship/airshipui/pkg/configs"
)

func TestHandleUnknownBaremetalSubComponent(t *testing.T) {
	request := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Baremetal,
		SubComponent: "fake_subcomponent",
	}

	acp := "testdata/testairshipconfig"
	kcp := "testdata/testkubeconfig"

	AirshipConfigPath = &acp
	KubeConfigPath = &kcp

	user := "test"
	response := HandleBaremetalRequest(&user, request)

	e := "Subcomponent fake_subcomponent not found"
	expected := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.Baremetal,
		SubComponent: "fake_subcomponent",
		Error:        &e,
	}

	assert.Equal(t, expected, response)
}
