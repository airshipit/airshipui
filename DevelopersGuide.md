# Airship UI Developer's Guide

## Prerequisites

1. Airship UI needs to be pointed to a Kubernetes Cluster. For development we recommending [setting up Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
2. Install [Go](https://golang.org/dl/) v1.13 or newer

## Getting Started

Let's clone the Airship UI repository and build

    git clone https://opendev.org/airship/airshipui
    cd airshipui
    make
    cd web
    npm install
    npm install --save-dev electron
    npm install electron-json-config

Now that Airship is built and we have a binary we can run it

    ./bin/airshipui

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
