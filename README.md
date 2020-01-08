# Octant Airship UI Plugin

A plugin built for use with [Octant](https://github.com/vmware/octant) to manage
declarative infrastructure using Airship.

## Prerequisites

- [Go 1.12+](https://golang.org/dl/)

## Getting Started

```
git clone https://opendev.org/airship/airshipui
cd airshipui
make install-plugins
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
