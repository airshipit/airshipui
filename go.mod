module opendev.org/airship/airshipui

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	opendev.org/airship/airshipctl v0.0.0-20200929174731-d4ea864eeba8
	sigs.k8s.io/kustomize/api v0.5.1
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191219154910-1528d4eea6dd
	sigs.k8s.io/kustomize/kyaml => sigs.k8s.io/kustomize/kyaml v0.4.1
)
