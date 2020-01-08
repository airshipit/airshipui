package commands

import (
	"fmt"
	"os"

	ctlenv "opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipui/internal/environment"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	disableAuto bool
	settings    *ctlenv.AirshipCTLSettings
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "airshipui",
	Short:   "airshipui is a graphical user interface for airship",
	Run:     launchOctant,
	Version: environment.Version(),
}

func init() {
	settings = &ctlenv.AirshipCTLSettings{}

	// Add options to rootCmd for configuration of airshipctl and kube config
	settings.InitFlags(rootCmd)

	// Load the airshipctl settings
	settings.InitConfig()

	// Add a 'version' command, in addition to the '--version' option that is auto created
	rootCmd.AddCommand(newVersionCmd())

	// Add flags for the underlying octant dashboard
	addDashboardFlags(rootCmd)
}

func launchOctant(cmd *cobra.Command, args []string) {

	kubeConfig := ""
	airshipKubeConfig := settings.KubeConfigPath()

	// If the kubeconfig specified on the command line (or defaulted to) does not exist,
	// then do not pass specify a kubeconfig.  This will permit the underlying octant
	// to use its own default, which can be either ~/.kube/config or grabbed from the
	// environment
	if fileExists(airshipKubeConfig) {
		kubeConfig = airshipKubeConfig
	}

	RunOctantWithOptions(cmd, kubeConfig, args)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Determine whether the given filename exists and is accessible
func fileExists(filename string) bool {
	f, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return !f.IsDir()
}
