module opendev.org/airship/airshipui

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	github.com/vmware-tanzu/octant v0.12.0
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	opendev.org/airship/airshipctl v0.0.0-20200518155418-7276dd68d8d0
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
