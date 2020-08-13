module opendev.org/airship/airshipui

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	opendev.org/airship/airshipctl v0.0.0-20200812155702-f61953bcf558
	sigs.k8s.io/kustomize/api v0.5.1
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
