module opendev.org/airship/airshipui

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/mattn/go-sqlite3 v1.14.3
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	opendev.org/airship/airshipctl v0.0.0-20201111150845-f2e305a56df7
	sigs.k8s.io/cli-utils v0.20.6
	sigs.k8s.io/kustomize/api v0.5.1
)

replace (
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191219154910-1528d4eea6dd
	sigs.k8s.io/kustomize/kyaml => sigs.k8s.io/kustomize/kyaml v0.4.1
)
