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

package testutil

import (
	"opendev.org/airship/airshipui/internal/configs"
)

// DummyDashboardConfig returns a populated Dashboard struct
func DummyDashboardConfig() configs.Dashboard {
	return configs.Dashboard{
		Name:     "dummy_dashboard",
		Protocol: "http",
		Hostname: "dummyhost",
		Port:     80,
		Path:     "fake/login/path",
	}
}

// DummyPluginDashboardConfig returns a populated PluginDashboard struct
func DummyPluginDashboardConfig() configs.PluginDashboard {
	return configs.PluginDashboard{
		Protocol: "http",
		FQDN:     "localhost",
		Port:     80,
		Path:     "index.html",
	}
}

// DummyExecutableConfig returns a populated Executable struct
func DummyExecutableConfig() configs.Executable {
	return configs.Executable{
		AutoStart: true,
		Filepath:  "/fake/path/to/executable",
		Args: []string{
			"--fakeflag",
			"fakevalue",
		},
	}
}

// DummyAuthMethodConfig returns a populated AuthMethod struct
func DummyAuthMethodConfig() *configs.AuthMethod {
	return &configs.AuthMethod{
		URL: "http://fake.auth.method.com/auth",
	}
}

// DummyPluginWithDashboardConfig returns a populated Plugin struct
// with a populated PluginDashboard
func DummyPluginWithDashboardConfig() configs.Plugin {
	d := DummyPluginDashboardConfig()
	e := DummyExecutableConfig()

	return configs.Plugin{
		Name:       "dummy_plugin_with_dash",
		Dashboard:  &d,
		Executable: &e,
	}
}

// DummyPluginNoDashboard returns a populated Plugin struct
// but omits the optional PluginDashboard
func DummyPluginNoDashboard() configs.Plugin {
	e := DummyExecutableConfig()

	return configs.Plugin{
		Name:       "dummy_plugin_no_dash",
		Executable: &e,
	}
}

// DummyNamespaceConfig returns a populated Namespace struct with
// a single Dashboard
func DummyNamespaceConfig() configs.Namespace {
	d := DummyDashboardConfig()

	return configs.Namespace{
		Name:       "dummy_namespace",
		Dashboards: []configs.Dashboard{d},
	}
}

// DummyClusterConfig returns a populated Cluster struct with
// a single Namespace
func DummyClusterConfig() configs.Cluster {
	n := DummyNamespaceConfig()

	return configs.Cluster{
		Name:       "dummy_cluster",
		BaseFqdn:   "dummy.cluster.local",
		Namespaces: []configs.Namespace{n},
	}
}

// DummyConfigNoAuth returns a populated Config struct but omits
// the optional AuthMethod
func DummyConfigNoAuth() configs.Config {
	p := DummyPluginWithDashboardConfig()
	pn := DummyPluginNoDashboard()
	c := DummyClusterConfig()

	return configs.Config{
		Plugins:  []configs.Plugin{p, pn},
		Clusters: []configs.Cluster{c},
	}
}

// DummyCompleteConfig returns a fully populated Config struct
func DummyCompleteConfig() configs.Config {
	a := DummyAuthMethodConfig()
	p := DummyPluginWithDashboardConfig()
	pn := DummyPluginNoDashboard()
	c := DummyClusterConfig()

	return configs.Config{
		AuthMethod: a,
		Plugins:    []configs.Plugin{p, pn},
		Clusters:   []configs.Cluster{c},
	}
}
