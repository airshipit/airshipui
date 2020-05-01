module opendev.org/airship/airshipui

go 1.13

require (
	github.com/gophercloud/gophercloud v0.9.0
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	github.com/vmware-tanzu/octant v0.12.0
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	opendev.org/airship/airshipctl v0.0.0-20200324160507-db6217f011b9
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33