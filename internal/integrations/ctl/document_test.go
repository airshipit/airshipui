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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"opendev.org/airship/airshipui/internal/configs"
)

const (
	testDocumentHTML string = "testdata/document.html"
)

func TestHandleDefaultDocumentRequest(t *testing.T) {
	html, err := ioutil.ReadFile(testDocumentHTML)
	require.NoError(t, err)

	request := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Document,
		SubComponent: configs.GetDefaults,
	}

	response := HandleDocumentRequest(request)

	expected := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Document,
		SubComponent: configs.GetDefaults,
		HTML:         string(html),
	}

	assert.Equal(t, expected, response)
}

func TestHandleUnknownDocumentSubComponent(t *testing.T) {
	request := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Document,
		SubComponent: "fake_subcomponent",
	}

	response := HandleDocumentRequest(request)

	expected := configs.WsMessage{
		Type:         configs.AirshipCTL,
		Component:    configs.Document,
		SubComponent: "fake_subcomponent",
		Error:        "Subcomponent fake_subcomponent not found",
	}

	assert.Equal(t, expected, response)
}
