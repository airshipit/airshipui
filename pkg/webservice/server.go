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

package webservice

import (
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/util/utilfile"
	"opendev.org/airship/airshipui/util/utilhttp"
)

// semaphore to signal the UI to authenticate

const (
	staticContent = "client/dist/airshipui"
)

// test if path and file exists, if it does send a page, else 404 for you
func serveFile(w http.ResponseWriter, r *http.Request) {
	filePath, filePathErr := utilfile.FilePath(staticContent, r.URL.Path)
	if filePathErr != nil {
		utilhttp.HandleErr(w, errors.WithStack(filePathErr), http.StatusInternalServerError)
		return
	}

	fileExists, fileExistsErr := utilfile.Exists(filePath)
	if fileExistsErr != nil {
		utilhttp.HandleErr(w, errors.WithStack(fileExistsErr), http.StatusInternalServerError)
		return
	}

	if fileExists {
		http.ServeFile(w, r, filePath)
	} else {
		// this is in an else to prevent a: superfluous response.WriteHeader call
		// TODO (aschie): Determine if this should do this on any 404, or if it should 404 a request
		http.ServeFile(w, r, staticContent)
	}
}

// getCertificates returns the cert chain in a way that the net/http server struct expects
func getCertificates() []tls.Certificate {
	cert, err := tls.LoadX509KeyPair(configs.UIConfig.WebService.PublicKey, configs.UIConfig.WebService.PrivateKey)
	if err != nil {
		log.Fatal("Unable to load certificates, check the definition in etc/airshipui.json")
	}
	var certSlice []tls.Certificate
	certSlice = append(certSlice, cert)
	return certSlice
}

// WebServer will run the handler functions for WebSockets
func WebServer() {
	webServerMux := http.NewServeMux()

	// hand off the websocket upgrade over http
	webServerMux.HandleFunc("/ws", onOpen)

	// establish routing to static angular client
	log.Debug("Attempting to serve static content from ", staticContent)
	webServerMux.HandleFunc("/", serveFile)

	// TODO: Figureout if we need to toggle the proxies on and off
	// start proxies for web based use
	startProxies()

	// Calculate the address and start on the host and port specified in the config
	addr := configs.UIConfig.WebService.Host + ":" + strconv.Itoa(configs.UIConfig.WebService.Port)
	log.Infof("Attempting to start webservice on %s", addr)

	// configure logging & TLS for the http server
	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         configs.UIConfig.WebService.Host,
			Certificates:       getCertificates(),
			MinVersion:         tls.VersionTLS13,
		},
		Handler:  webServerMux,
		ErrorLog: log.Logger(),
	}

	// kick off the server, and good luck
	log.Fatal(server.ListenAndServeTLS("", ""))
}
