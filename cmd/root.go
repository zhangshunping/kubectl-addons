/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile   string
	nodelabel string
	namespace string
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kubectl-addons",
	Short: "kubectl Educoder resouce operater",
	Long:  ``,
	Args: cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.SuggestionsMinimumDistance = 1

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-addons.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "-n default")
	RootCmd.PersistentFlags().StringVarP(&nodelabel, "nodeslector", "l", "", "-l type=others")
	RootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kubectl-addons" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kubectl-addons")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
