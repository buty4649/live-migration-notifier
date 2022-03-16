package cmd

import (
	"os"

	"buty4649/live-migration-notifier/cmdutil"
	"buty4649/live-migration-notifier/notifier"
	"buty4649/live-migration-notifier/version"

	"github.com/spf13/cobra"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:          "live-migration-notifier",
		Version:      version.Version,
		Short:        "Notify Slack of OpenStack live-migration.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRootCmd()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yml", "config file")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	return rootCmd.Execute()
}

func runRootCmd() error {
	config, err := cmdutil.ReadConfigFile(cfgFile)
	if err != nil {
		return err
	}

	n := notifier.Init(config.Uri(), config.SlackWebHookUrl)
	err = n.Start()
	if err != nil {
		return err
	}

	return nil
}
