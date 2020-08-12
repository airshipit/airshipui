# Airship UI

Airship UI is a browser based application that is designed to allow you to interact with Airship components, find and
connect to the kubernetes cluster and use plugins to tie together a singular dashboard to view addons without the need
to go to a separate url or application.

## Prerequisites

- A working [kubernetes](https://kubernetes.io/) or [airship](https://wiki.openstack.org/wiki/Airship) installation
- [Go 1.13+](https://golang.org/dl/)

## Getting Started

```
git clone https://opendev.org/airship/airshipui
cd airshipui
make # Note running behind a proxy can cause issues, notes on solving is in the Appendix of the Developer's Guide
```

## Adding Additional Functionality

Airship UI can be seamlessly integrated with service dashboards and other web-based tools by providing the necessary
configuration in $HOME/.airship/airshipui.json.

To add service dashboards, create a section at the top-level of airshipui.json as follows:

```
"dashboards": [
        {
            "name": "Ceph",
            "baseURL": "https://ceph-dash.example.domain",
            "path": ""
        },
        {
            "name": "Horizon",
            "baseURL": "http://horizon",
            "path": "dashboard/auth/login"
        }
]
```

In the above example, the configuration for Horizon specifies a service dashboard available at
'http://horizon/dashboard/auth/login'."


## Developer's Guide

Instructions on setting up a development environment and more details can be found in the
[Developer's Guide](docs/source/developers.md)