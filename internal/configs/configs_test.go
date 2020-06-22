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
package configs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"opendev.org/airship/airshipui/internal/configs"
	"opendev.org/airship/airshipui/testutil"
)

const (
	fakeFile string = "/fake/config/path"
	testFile string = "testdata/airshipui.json"
)

func TestSetUIConfig(t *testing.T) {
	conf := configs.Config{
		Clusters: []configs.Cluster{
			testutil.DummyClusterConfig(),
		},
		Plugins: []configs.Plugin{
			testutil.DummyPluginWithDashboardConfig(),
			testutil.DummyPluginNoDashboard(),
		},
		AuthMethod: testutil.DummyAuthMethodConfig(),
	}

	err := configs.SetUIConfig(testFile)
	require.NoError(t, err)

	assert.Equal(t, conf, configs.UIConfig)
}

func TestFileNotFound(t *testing.T) {
	err := configs.SetUIConfig(fakeFile)
	assert.Error(t, err)
}
