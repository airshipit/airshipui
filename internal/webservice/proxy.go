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
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"opendev.org/airship/airshipui/internal/configs"
)

// map of proxy targets which will be used based on the request
var proxyMap = map[string]*url.URL{}

type transport struct {
	http.RoundTripper
}

var _ http.RoundTripper = &transport{}

func (t *transport) RoundTrip(request *http.Request) (response *http.Response, err error) {
	// TODO: inject headers here for bearer token auth
	// example:
	// request.Header.Add("X-Auth-Token", "<token>")

	response, err = t.RoundTripper.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	// TODO: inject headers here for cookie auth
	// example:
	// response.Header.Add("Set-Cookie", "sessionid=<session>; expires=<date>; HttpOnly; Max-Age=3597; Path=/;")

	return response, nil
}

// handle a proxy request
// this is essentially a man in the middle attack that allows us to inject headers for single sign on
func handleProxy(response http.ResponseWriter, request *http.Request) {
	// retrieve the target URL from the proxy map
	target := proxyMap[request.Host]

	// short circuit for bad targets blowing up the backend
	if target == nil {
		response.WriteHeader(http.StatusInternalServerError)
		if _, err := response.Write([]byte("500 - Unable to locate proxy for request!")); err != nil {
			log.Println("Error writing response for proxy not found: ", err)
		}

		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &transport{http.DefaultTransport}

	// Update the headers to allow for SSL redirection
	request.URL.Host = target.Host
	request.URL.Scheme = target.Scheme

	host := request.Header.Get("Host")
	request.Header.Set("X-Forwarded-Host", host)
	request.Header.Set("X-Forwarded-For", host)

	proxy.ServeHTTP(response, request)
}

func getRandomPort() (string, error) {
	// get a random port for the proxy
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	// close the port so we can start the proxy
	defer listener.Close()

	// get the string of the port
	return listener.Addr().String(), nil
}

// proxyServer will proxy dashboard connections allowing us to inject headers
func proxyServer(port string) {
	proxyServerMux := http.NewServeMux()

	// some things may need a helping hand with the headers so we'll proxy it for them
	proxyServerMux.HandleFunc("/", handleProxy)

	if err := http.ListenAndServe(port, proxyServerMux); err != nil {
		log.Fatal("Error starting proxy: ", err)
	}
}

// helper function that kicks off all proxies prior to the start of the website
func startProxies() {
	for index, dashboard := range configs.UIConfig.Dashboards {
		port, err := getRandomPort()
		if err != nil {
			log.Fatal("Error starting proxy, unable to allocate port:", err)
		}

		// this will persuade the UI to use the proxy and not the original host
		dashboard.IsProxied = true

		// cache up the target for the proxy url
		target, err := url.Parse(dashboard.BaseURL)
		if err != nil {
			log.Println(err)
		}

		// set the target for the proxied request to the original url
		proxyMap[port] = target

		// set the target for the link in the ui to the proxy address
		dashboard.BaseURL = "http://" + port

		// kick off proxy
		log.Printf("Attempting to start proxy for %s on: %s\n", dashboard.Name, port)

		// set the dashboard from this point on to go to the proxy
		configs.UIConfig.Dashboards[index] = dashboard

		// and away we go.........
		go proxyServer(port)
	}
}
