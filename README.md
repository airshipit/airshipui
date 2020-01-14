# Octant Airship UI Plugin

Airship UI is a wrapper around [Octant](https://github.com/vmware/octant) together with Octant plugin(s) that allows you to view your kubernetes cluster.  The airshipui command uses airshipctl's configuration in order to find and connect to the kubernetes cluster, launches the Octant server process, and directs your browser to the user interface.

Several plugins will be delivered with airshipui. The first, argoui, is a plugin that embeds the [Argo UI](https://github.com/argoproj/argo-ui) interface within octant, and it requires that argo be installed on your kubernetes cluster (this should be the case by default with Airship 2.0, which uses argo as its workflow engine)


## Prerequisites

- A working [kubernetes](https://kubernetes.io/) or [airship](https://wiki.openstack.org/wiki/Airship) installation
- [Go 1.12+](https://golang.org/dl/)
- [Octant](https://github.com/vmware-tanzu/octant)
- [Argo](https://github.com/argoproj/argo/blob/master/README.md)


## Getting Started

```
git clone https://opendev.org/airship/airshipui
cd airshipui
make build install-plugins
```

`make install-plugins` builds and installs the plugin to
`$HOME/.config/octant/plugin/`.

The next time Octant is run it will include plugins in the above directory.
Further information for running Octant can be found in the
[Octant Repo](https://github.com/vmware/octant).

If you would like to just build the plugin use `make`.

## Architecture

airshipui is an executable that wraps Octant.  When it is launched, it processes its own set of command-line options, performs any
necessary custom startup tasks such as reading the airshipctl config file, then normally calls the function to instantiate Octant.
This repository also contains airship plugins that will be generated as standard octant plugins, which are separate binaries.

## Developer's Guide

Step by step sample installation and more details can be found in the [Developer's Guide](DevelopersGuide.md).