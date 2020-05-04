# Airship UI Developer's Guide

## Prerequisites
1. [Go](https://golang.org/dl/) v1.13 or newer
2. [Nodejs](https://nodejs.org/en/download/) v12 or newer

## Getting Started

Clone the Airship UI repository and build

    git clone https://opendev.org/airship/airshipui
    cd airshipui
    make
    make install-octant-plugins # (if running with octant)
    cd web
    npm install
    npm install --save-dev electron
    npm install electron-json-config

Run the airshipui binary

    ./bin/airshipui

# Authentication
## Pluggable authentication methods
The AirshipUI is not designed to create authentication credentials but to have them supplied to it either by a configuration or by an external entity.  The expectation is that there will be an external URL that will handle authentication for the system which may need to be modified or created.  The endpoint will need to be able to forward a [bearer token](https://oauth.net/2/bearer-tokens/), [basic auth](https://en.wikipedia.org/wiki/Basic_access_authentication) or cookie data to the Airship UI backend service.

To configure the pluggable authentication the following must be added to the $HOME/.airshipui/airshipui.json file:
```
"authMethod": {
	"url": "<protocol>://<host:port>/<path>/<method>"
}
```
Note: By default the system will start correctly without any authentication urls supplied to the configuration.  The expectation is that AirshipUI will be running in a minimal least authorized configuration.

## Example Auth Server
There is an example authentication server in examples/authentication/main.go.  These endpoints can be added to the $HOME/.airshipui/airshipui.json and will allow the system to show a basic authentication test.
1. Basic auth on http://127.0.0.1:12321/basic-auth
2. Cookie based auth on http://127.0.0.1:12321/cookie
3. OAuth JWT (JSON Web Token) on http://127.0.0.1:12321/oauth

To start the system cd to the root of the AirshipUI repository and execute:
```
go run examples/authentication/main.go
```
### Example Auth Server Credentials
+ The example auth server id is: airshipui
+ The example auth server password is: Open Sesame!

# Plugins
## Octant
[Octant](https://github.com/vmware-tanzu/octant) is a tool for developers to understand how applications run on a Kubernetes cluster. It aims to be part of the developer's toolkit for gaining insight and approaching complexity found in Kubernetes. Octant offers a combination of introspective tooling, cluster navigation, and object management along with a plugin system to further extend its capabilities.

Octant needs to be pointed to a Kubernetes Cluster. For development it is recommended to use [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

### How to get and build Octant
If you are going to do serious Octant development you will need to adhere to [Octant's Hacking Guide](https://github.com/vmware-tanzu/octant/blob/master/HACKING.md) which includes information on how to build Octant and the steps to push changes to them.

### Running the example
Build the octant plugin executable
```
make install-octant-plugins
```
Run the octant binary and the plugin should show "Hello World just some text on the page" under the http://127.0.0.1:7777/#/airshipui-example-plugin url.

## Appendix

### Minikube

[Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) runs a single-node Kubernetes cluster for users looking to try out Kubernetes or develop with it day-to-day.  Installation instructions are available on the kubernetes website: https://kubernetes.io/docs/tasks/tools/install-minikube/).  If you are running behind a proxy it may be necessary to follow the steps outlined in the [How to use an HTTP/HTTPS proxy with minikube](https://minikube.sigs.k8s.io/docs/reference/networking/proxy/) website.

### Optional proxy settings

#### Environment settings for wget or curl

If your network has a proxy that prevents successful curls or wgets you may need to set the proxy environment variables.  The local ip is included in the no_proxy setting to prevent any local running process that may attempt api calls against it from being sent through the proxy for the request:

    ```
    export http_proxy=<proxy_host>:<proxy_port>
    export HTTP_PROXY=<proxy_host>:<proxy_port>
    export https_proxy=<proxy_host>:<proxy_port>
    export HTTPS_PROXY=<proxy_host>:<proxy_port>
    export no_proxy=localhost,127.0.0.1,<LOCAL_IP>
    export NO_PROXY=localhost,127.0.0.1,<LOCAL_IP>
    ```