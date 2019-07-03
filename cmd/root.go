/*
Copyright © 2019 Thibaut Tauveron <thibaut.tauveron@gmail.com>

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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ttauveron/git-follow-up/git"
	"os"
)

var cfgFile, configPath, gitPath string
var config Config

type Config struct {
	Repositories []git.Repository
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "git-follow-up",
	Long: `        _ _           __      _ _                                  
       (_) |         / _|    | | |                                 
   __ _ _| |_ ______| |_ ___ | | | _____      ________ _   _ _ __  
  / _` + "`" + ` | | __|______|  _/ _ \| | |/ _ \ \ /\ / /______| | | | '_ \
 | (_| | | |_       | || (_) | | | (_) \ V  V /       | |_| | |_) |
  \__, |_|\__|      |_| \___/|_|_|\___/ \_/\_/         \__,_| .__/
   __/ |                                                    | |
  |___/                                                     |_|          

Keeps track of contributions made on multiple git repositories described in a yaml configuration file.
Those repositories can be hosted on any platform, and accessed through ssh, https, with or without an access token.`,
	BashCompletionFunction: bash_completion_func,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-follow-up/config.yaml)")
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

		// Search config in home directory with name ".git-follow-up" (without extension).
		configPath = home + "/.git-follow-up"
		gitPath = configPath + "/git/"
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")

		// Create repositories folder if not exists
		_ = os.MkdirAll(configPath+"/git", 0700)

		//If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			err = viper.Unmarshal(&config)
			for i := 0; i < len(config.Repositories); i++ {
				config.Repositories[i].LocalPath = gitPath + config.Repositories[i].Name
			}

			if err != nil {
				panic("Unable to unmarshal config")
			}
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

}
