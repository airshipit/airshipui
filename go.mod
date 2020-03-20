module opendev.org/airship/airshipui

go 1.12

require (
	github.com/go-cmd/cmd v1.2.0
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	github.com/vmware-tanzu/octant v0.10.2-0.20200320182255-15a53b6af867
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	opendev.org/airship/airshipctl v0.0.0-20200319213630-b2a602fa07e0
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
