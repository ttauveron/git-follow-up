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
	"github.com/spf13/cobra"
	"github.com/ttauveron/git-follow-up/git"
	"sync"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Synchronizes remote git repositories to local copies",
	Long: `Synchronizes remote git repositories to local copies
This operation may initially take some time for large repositories...
`,
	Run: func(cmd *cobra.Command, args []string) {
		var repoList []git.Repository

		// Skip update on non-matching labels
		if cmd.Flags().Changed("label") {
			for _, repo := range config.Repositories {
				filterLabels, _ := cmd.Flags().GetStringSlice("label")
				if git.ContainsAll(repo.Labels, filterLabels) {
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
	updateCmd.Flags().StringSlice("label", []string{}, "filters by project labels")
	rootCmd.AddCommand(updateCmd)

}

// Syncing all repositories defined in the `config.yaml` file
func UpdateRepos(repos []git.Repository) {
	var wg sync.WaitGroup
	for _, repo := range repos {

		wg.Add(1)

		go func(repository git.Repository) {
			err := repository.SyncRepo()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			wg.Done()
		}(repo)
	}
	wg.Wait()
}
