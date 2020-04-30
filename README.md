# Airship UI

Airship UI is an [electron](https://www.electronjs.org/) that is designed to allow you to interact with Airship components, find and connect to the kubernetes cluster and use plugins to tie together a singular dashboard to view addons without the need to go to a separate url or application for each. 

## Prerequisites

- A working [kubernetes](https://kubernetes.io/) or [airship](https://wiki.openstack.org/wiki/Airship) installation
- [Go 1.13+](https://golang.org/dl/)

## Getting Started

```
git clone https://opendev.org/airship/airshipui
cd airshipui
make build
```

## Adding Additional Functionality

Airship UI can be seamlessly integrated with dashboards and other web-based tools by making the application aware of their service URLs.
To do this, create the file $HOME/.airshipui/plugins.json with the following structure:

```
{
  "external_dashboards":
    [
      {
        "name":"Ceph",
        "url":"https://127.0.0.1:51515"
      },
      {
        "name":"Octant",
        "url":"http://127.0.0.1:7777"
      }
    ]
}
```

Once the file is in place, dashboards can be accessed through the Plugins drop-down menu in the Airship UI navigation bar. Dashboards can
be added or removed while the application is running, and changes will take effect after the current view is refreshed.

## Developer's Guide

Instructions on setting up a development environment and more details can be found in the [Developer's Guide](DevelopersGuide.md)
