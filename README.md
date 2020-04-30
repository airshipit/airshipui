# Airship UI

Airship UI is an [electron](https://www.electronjs.org/) that is designed to allow you to interact with Airship components, find and connect to the kubernetes cluster and use plugins to tie together a singular dashboard to view addons without the need to go to a separate url or application for each. 

## Prerequisites

- A working [kubernetes](https://kubernetes.io/) or [airship](https://wiki.openstack.org/wiki/Airship) installation
- [Go 1.13+](https://golang.org/dl/)
- [npm](https://www.npmjs.com/)

## Getting Started

```
git clone https://opendev.org/airship/airshipui
cd airshipui
make build
cd web
npm install
npm install --save-dev electron
npm install electron-json-config
```

## Developer's Guide

Instructions on setting up a development environment and more details can be found in the [Developer's Guide](DevelopersGuide.md)
