# Airship UI Developer's Guide

## Prerequisites
1. [Go](https://golang.org/dl/) v1.13 or newer
2. [Nodejs](https://nodejs.org/en/download/) v12 or newer

## Getting Started

Clone the Airship UI repository and build

    git clone https://opendev.org/airship/airshipui
    cd airshipui
    make
    make install-npm-modules # Note running behind a proxy can cause issues, notes on solving is in the Appendix
    make examples # (optional)
    make install-octant-plugins # (if running with octant)



Run the airshipui binary

    ./bin/airshipui

# Running on a separate client & server
For development purposes it could be advantageous to split the webservice and the UI across 2 machines.

## To start the webservice on a remote system without starting the ui

    bin/airshipui --headless

This will require you to tunnel the connection to the remote machine:

    ssh -L 8080:localhost:8080 <id>@<remote_host>

## To start the UI to attach to a remote machine

    bin/airshipui --remote

## Running the webservice on a remote machine with plugins
The plugins can run on the remote system but you will need to add additional SSH tunnels to the machine in order for it to work.
The plugins can also run locally but the remote server will still need the definition on the remote system to notify the UI that the plugins are available.

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

### Issues with npm / npx and the electron go module
It is possible that the airship ui will exit with an error code 1 when attempting to start the first time:

    $ bin/airshipui
    2020/05/28 14:54:16 Attempting to start webservice on localhost:8080
    2020/05/28 14:54:16 Exit electron 1

You can attempt a manual start of electron from the root of the airshipui tree to see a more verbose log

    $ npm start --prefix web

    > electron-poc@0.0.1 start /home/ubuntu/airshipui/web
    > electron .

    /home/ubuntu/airshipui/web/node_modules/electron/index.js:14
        throw new Error('Electron failed to install correctly, please delete node_modules/electron and try installing again')
        ^

    Error: Electron failed to install correctly, please delete node_modules/electron and try installing again
        at getElectronPath (/home/ubuntu/airshipui/web/node_modules/electron/index.js:14:11)
        at Object.<anonymous> (/home/ubuntu/airshipui/web/node_modules/electron/index.js:18:18)
        at Module._compile (internal/modules/cjs/loader.js:1133:30)
        at Object.Module._extensions..js (internal/modules/cjs/loader.js:1153:10)
        at Module.load (internal/modules/cjs/loader.js:977:32)
        at Function.Module._load (internal/modules/cjs/loader.js:877:14)
        at Module.require (internal/modules/cjs/loader.js:1019:19)
        at require (internal/modules/cjs/helpers.js:77:18)
        at Object.<anonymous> (/home/ubuntu/airshipui/web/node_modules/electron/cli.js:3:16)
        at Module._compile (internal/modules/cjs/loader.js:1133:30)
    npm ERR! code ELIFECYCLE
    npm ERR! errno 1
    npm ERR! electron-poc@0.0.1 start: `electron .`
    npm ERR! Exit status 1
    npm ERR!
    npm ERR! Failed at the electron-poc@0.0.1 start script.
    npm ERR! This is probably not a problem with npm. There is likely additional logging output above.

    npm ERR! A complete log of this run can be found in:
    npm ERR!     /home/ubuntu/.npm/_logs/2020-05-28T14_55_52_327Z-debug.log
    $

This is likely due to a problem with the electron install.  If you cd to the web directory and attempt to install you may see an error:

    npm install electron

    > electron@8.3.0 postinstall /home/ubuntu/airshipui/web/node_modules/electron
    > node install.js

    (node:19823) UnhandledPromiseRejectionWarning: RequestError: connect ETIMEDOUT 192.30.255.113:443
        at ClientRequest.<anonymous> (/home/ubuntu/airshipui/web/node_modules/got/source/request-as-event-emitter.js:178:14)
        at Object.onceWrapper (events.js:417:26)
        at ClientRequest.emit (events.js:322:22)
        at ClientRequest.EventEmitter.emit (domain.js:482:12)
        at ClientRequest.origin.emit (/home/ubuntu/airshipui/web/node_modules/@szmarczak/http-timer/source/index.js:37:11)
        at TLSSocket.socketErrorListener (_http_client.js:426:9)
        at TLSSocket.emit (events.js:310:20)
        at TLSSocket.EventEmitter.emit (domain.js:482:12)
        at emitErrorNT (internal/streams/destroy.js:92:8)
        at emitErrorAndCloseNT (internal/streams/destroy.js:60:3)
    (node:19823) UnhandledPromiseRejectionWarning: Unhandled promise rejection. This error originated either by throwing inside of an async function without a catch block, or by rejecting a promise which was not handled with .catch(). To terminate the node process on unhandled promise rejection, use the CLI flag `--unhandled-rejections=strict` (see https://nodejs.org/api/cli.html#cli_unhandled_rejections_mode). (rejection id: 1)
    (node:19823) [DEP0018] DeprecationWarning: Unhandled promise rejections are deprecated. In the future, promise rejections that are not handled will terminate the Node.js process with a non-zero exit code.
    npm WARN electron-poc@0.0.1 No repository field.

    + electron@8.3.0
    updated 1 package and audited 361 packages in 143.19s

    9 packages are looking for funding
    run `npm fund` for details

    found 1 low severity vulnerability
    run `npm audit fix` to fix them, or `npm audit` for details

If you're running behind a corporate proxy this is the workaround:

    npx cross-env ELECTRON_GET_USE_PROXY=true GLOBAL_AGENT_HTTPS_PROXY=http://<proxy_host>:<proxy_port> npm install electron

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