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

# Plugins
## Octant
[Octant](https://github.com/vmware-tanzu/octant) is a tool for developers to understand how applications run on a Kubernetes cluster. It aims to be part of the developer's toolkit for gaining insight and approaching complexity found in Kubernetes. Octant offers a combination of introspective tooling, cluster navigation, and object management along with a plugin system to further extend its capabilities.

Octant needs to be pointed to a Kubernetes Cluster. For development we recommend [setting up Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

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