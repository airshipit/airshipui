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
package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fakeFile        string = "/fake/config/path"
	testFile        string = "testdata/airshipui.json"
	invalidTestFile string = "testdata/airshipui_invalid.json"
)

// DummyDashboardsConfig returns an array of populated Dashboard structs
func dummyDashboardsConfig() []Dashboard {
	return []Dashboard{
		{
			Name:    "dummy_dashboard",
			BaseURL: "http://dummyhost",
			Path:    "fake/login/path",
		},
	}
}

func dummyAuthMethodConfig() *AuthMethod {
	return &AuthMethod{
		URL: "http://fake.auth.method.com/auth",
	}
}

func TestSetUIConfig(t *testing.T) {
	conf := Config{
		Dashboards: dummyDashboardsConfig(),
		AuthMethod: dummyAuthMethodConfig(),
	}

	err := SetUIConfig(testFile)
	require.NoError(t, err)

	assert.Equal(t, conf, UIConfig)

	err = SetUIConfig(invalidTestFile)
	require.Error(t, err)
}

func TestFileNotFound(t *testing.T) {
	err := SetUIConfig(fakeFile)
	assert.Error(t, err)
}
