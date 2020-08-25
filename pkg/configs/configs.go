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
	"crypto/rsa"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipui/pkg/cryptography"
	"opendev.org/airship/airshipui/pkg/log"
)

// variables related to UI config
var (
	UIConfig     Config
	UIConfigFile string
	etcDir       *string
)

// Config basic structure to hold configuration params for Airship UI
type Config struct {
	WebService *WebService       `json:"webservice,omitempty"`
	AuthMethod *AuthMethod       `json:"authMethod,omitempty"`
	Dashboards []Dashboard       `json:"dashboards,omitempty"`
	Users      map[string]string `json:"users,omitempty"`
}

// AuthMethod structure to hold authentication parameters
type AuthMethod struct {
	Type  string   `json:"type,omitempty"`
	Value []string `json:"values,omitempty"`
	URL   string   `json:"url,omitempty"`
}

// WebService describes the things we need to know to start the web container
type WebService struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	PublicKey  string `json:"publicKey,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}

// Authentication structure to hold authentication parameters
type Authentication struct {
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

// Dashboard structure
type Dashboard struct {
	Name      string `json:"name,omitempty"`
	BaseURL   string `json:"baseURL,omitempty"`
	Path      string `json:"path,omitempty"`
	IsProxied bool   `json:"isProxied,omitempty"`
}

// WsRequestType is used to set the specific types allowable for WsRequests
type WsRequestType string

// WsComponentType is used to set the specific component types allowable for WsRequests
type WsComponentType string

// WsSubComponentType is used to set the specific subcomponent types allowable for WsRequests
type WsSubComponentType string

// constants related to specific request/component/subcomponent types for WsRequests
const (
	CTL   WsRequestType = "ctl"
	UI    WsRequestType = "ui"
	Alert WsRequestType = "alert"

	Authcomplete WsComponentType = "authcomplete"
	SetConfig    WsComponentType = "setConfig"
	Initialize   WsComponentType = "initialize"
	Keepalive    WsComponentType = "keepalive"
	CTLConfig    WsComponentType = "config"
	Baremetal    WsComponentType = "baremetal"
	Document     WsComponentType = "document"
	Auth         WsComponentType = "auth"

	// auth sub components
	Approved     WsSubComponentType = "approved"
	Authenticate WsSubComponentType = "authenticate"
	Denied       WsSubComponentType = "denied"
	Refresh      WsSubComponentType = "refresh"
	Validate     WsSubComponentType = "validate"

	// ctl components
	GetDefaults         WsSubComponentType = "getDefaults"
	GenerateISO         WsSubComponentType = "generateISO"
	DocPull             WsSubComponentType = "docPull"
	Yaml                WsSubComponentType = "yaml"
	YamlWrite           WsSubComponentType = "yamlWrite"
	GetYaml             WsSubComponentType = "getYaml"
	GetSource           WsSubComponentType = "getSource"
	GetRendered         WsSubComponentType = "getRendered"
	GetPhaseTree        WsSubComponentType = "getPhaseTree"
	GetPhaseSourceFiles WsSubComponentType = "getPhaseSource"
	GetPhaseDocuments   WsSubComponentType = "getPhaseDocs"
	GetTarget           WsSubComponentType = "getTarget"
)

// WsMessage is a request / return structure used for websockets
type WsMessage struct {
	// base components of a message
	SessionID    string             `json:"sessionID,omitempty"`
	Type         WsRequestType      `json:"type,omitempty"`
	Component    WsComponentType    `json:"component,omitempty"`
	SubComponent WsSubComponentType `json:"subComponent,omitempty"`
	Timestamp    int64              `json:"timestamp,omitempty"`

	// additional conditional components that may or may not be involved in the request / response
	Error           string      `json:"error,omitempty"`
	IsAuthenticated bool        `json:"isAuthenticated,omitempty"`
	Message         string      `json:"message,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	YAML            string      `json:"yaml,omitempty"`
	Name            string      `json:"name,omitempty"`
	ID              string      `json:"id,omitempty"`
	Token           *string     `json:"token,omitempty"`

	// used for auth
	Authentication *Authentication `json:"authentication,omitempty"`

	// information related to the init of the UI
	Dashboards      []Dashboard             `json:"dashboards,omitempty"`
	AuthMethod      *AuthMethod             `json:"authMethod,omitempty"`
	AuthInfoOptions *config.AuthInfoOptions `json:"authInfoOptions,omitempty"`
	ContextOptions  *config.ContextOptions  `json:"contextOptions,omitempty"`
	ClusterOptions  *config.ClusterOptions  `json:"clusterOptions,omitempty"`
}

// SetUIConfig sets the UIConfig object with values obtained from
// airshipui.json, located at 'filename'
// TODO: add watcher to the json file to reload conf on change (maybe not needed)
func SetUIConfig() error {
	f, err := os.Open(UIConfigFile)
	if err != nil {
		return checkConfigs()
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &UIConfig)
	if err != nil {
		return err
	}

	return checkConfigs()
}

func checkConfigs() error {
	writeFile := false
	if UIConfig.WebService == nil {
		writeFile = true
		log.Debug("No UI config found, generating ssl keys & host & port info")
		err := setEtcDir()
		if err != nil {
			return err
		}

		privateKeyFile := filepath.Join(*etcDir, "key.pem")
		publicKeyFile := filepath.Join(*etcDir, "cert.pem")

		err = writeTestSSL(privateKeyFile, publicKeyFile)
		if err != nil {
			return err
		}

		UIConfig.WebService = &WebService{
			Host:       "localhost",
			Port:       10443,
			PublicKey:  publicKeyFile,
			PrivateKey: privateKeyFile,
		}
		err = cryptography.TestCertValidity(publicKeyFile)
		if err != nil {
			return err
		}
	}
	if UIConfig.Users == nil {
		writeFile = true
		err := createDefaultUser()
		if err != nil {
			return err
		}
	}

	if writeFile {
		bytes, err := json.Marshal(UIConfig)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(UIConfigFile, bytes, 0600)
	}
	return nil
}

func createDefaultUser() error {
	hash := sha512.New()
	_, err := hash.Write([]byte("admin"))
	if err != nil {
		return err
	}
	UIConfig.Users = map[string]string{"admin": hex.EncodeToString(hash.Sum(nil))}
	return nil
}

func writeTestSSL(privateKeyFile string, publicKeyFile string) error {
	// get and write out private key
	log.Warnf("Generating private key %s.  DO NOT USE THIS FOR PRODUCTION", privateKeyFile)
	privateKey, err := getAndWritePrivateKey(privateKeyFile)
	if err != nil {
		return err
	}

	// get and write out public key
	log.Warnf("Generating public key %s.  DO NOT USE THIS FOR PRODUCTION", publicKeyFile)
	err = getAndWritePublicKey(publicKeyFile, privateKey)
	if err != nil {
		return err
	}

	return nil
}

func getAndWritePrivateKey(fileName string) (*rsa.PrivateKey, error) {
	privateKeyBytes, privateKey, err := cryptography.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(fileName, privateKeyBytes, 0600)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func getAndWritePublicKey(fileName string, privateKey *rsa.PrivateKey) error {
	publicKeyBytes, err := cryptography.GeneratePublicKey(privateKey)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, publicKeyBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func setEtcDir() error {
	if etcDir == nil {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}
		dir, err = filepath.Abs(filepath.Join(path.Dir(dir), "etc"))
		if err != nil {
			return err
		}
		etcDir = &dir
	}
	return nil
}
