# Octant Airship UI Plugin

A plugin built for use with [Octant](https://github.com/vmware/octant) to manage
declarative infrastructure using Airship.

## Prerequisites

- [Go 1.12+](https://golang.org/dl/)

## Getting Started

```
git clone https://opendev.org/airship/airshipui
cd airshipui
make
```

The default make target builds and installs the plugin to
`$HOME/.config/octant/plugin/`

The next time Octant is run it will include plugins in the above directory.
Further information for running Octant can be found in the
[Octant Repo](https://github.com/vmware/octant).
