/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/takutakahashi/github-token-renewer/pkg/config"
	"github.com/takutakahashi/github-token-renewer/pkg/github.go"
	"github.com/takutakahashi/github-token-renewer/pkg/output"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-token-renewer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Error(err)
			return
		}
		cfg, err := config.Load(cfgPath)
		if err != nil {
			logrus.Error(err)
			return
		}
		app, err := github.NewApp(*cfg, nil)
		if err != nil {
			logrus.Error(err)
			return
		}
		tokenMap, err := app.GenerateInstallationToken()
		if err != nil {
			logrus.Error(err)
			return
		}
		for installationID, token := range tokenMap {
			for _, in := range cfg.Installations {
				if in.ID != installationID {
					continue
				}
				k, err := output.NewKubernetes(*in.Output.KubernetesSecret)
				if err != nil {
					logrus.Error(err)
					continue
				}
				if err := k.Output(token); err != nil {
					logrus.Error(err)
					continue
				}

			}
		}
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github-token-renewer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("config", "c", "./config.yaml", "Config path")
}
