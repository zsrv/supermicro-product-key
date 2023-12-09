package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zsrv/supermicro-product-key/pkg/build"
)

var logVerbosity int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "supermicro-product-key",
	Short:   "Supermicro Product Key Utility",
	Version: build.Version(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Stops usage from printing for errors in RunE, while still allowing
		// usage to print for cobra errors (e.g. Args count mismatch)
		cmd.SilenceUsage = true
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().CountVarP(&logVerbosity, "verbose", "v", "increase log verbosity (specify up to 4 times)")
}

func initConfig() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel - zerolog.Level(logVerbosity))
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}).With().Caller().Logger()
}
