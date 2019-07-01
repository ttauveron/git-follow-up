/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"git-follow-up/internal"
	"github.com/spf13/cobra"
	"sync"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var repoList []internal.Repository

		// Skip update on non-matching labels
		if cmd.Flags().Changed("label") {
			for _, repo := range config.Repositories {
				filterLabels, _ := cmd.Flags().GetStringSlice("label")
				if internal.ContainsAll(repo.Labels, filterLabels) {
					repoList = append(repoList, repo)
				}
			}
		} else {
			repoList = append(repoList, config.Repositories...)
		}

		UpdateRepos(repoList)
	},
}

func init() {
	updateCmd.Flags().StringSlice("label", []string{}, "label")
	rootCmd.AddCommand(updateCmd)

}

// Syncing all repositories defined in the `config.yaml` file
func UpdateRepos(repos []internal.Repository) {
	var wg sync.WaitGroup
	for _, repo := range repos {

		wg.Add(1)

		go func(repository internal.Repository) {
			err := repository.SyncRepo()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			wg.Done()
		}(repo)
	}
	wg.Wait()
}
