package environment

var (
	// version will be overridden by ldflags supplied in Makefile
	version = "(dev-version)"
)

func Version() string {
	return version
}
