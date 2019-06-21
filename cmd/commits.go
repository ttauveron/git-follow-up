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
	"github.com/spf13/pflag"
	"sort"
	"sync"
)

var filter *internal.Filter

// commitsCmd represents the commits command
var commitsCmd = &cobra.Command{
	Use:   "commits",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		filter = internal.NewFilter(cmd.Flags())

		// Sync repos if update flag is provided
		doUpdate, err := cmd.Flags().GetBool("update")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if doUpdate {
			UpdateRepos(cmd.Flags())
		}

		var commits []internal.Commit

		// Listing log messages of repositories
		for _, repo := range config.Repositories {
			// Skip update on non-matching labels
			if cmd.Flags().Changed("label") && !internal.ContainsAll(repo.Labels, filter.Labels) {
				continue
			}
			cs, err := repo.ListCommits(*filter)
			commits = append(commits, cs...)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}

		sort.Sort(internal.ByDate(commits))
		for _, v := range commits {
			fmt.Println(v)
		}

	},
}

// Syncing all repositories defined in the `config.yaml` file
func UpdateRepos(flags *pflag.FlagSet) {
	var wg sync.WaitGroup
	for _, repo := range config.Repositories {

		// Skip update on non-matching labels
		if flags.Changed("label") && !internal.ContainsAll(repo.Labels, filter.Labels) {
			continue
		}
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

func init() {

	//TODO flag validation
	commitsCmd.Flags().String("from", "wtd", "ytd, mtd, wtd, yesterday, today, [dayOfWeek], [yyyy-MM-dd]")
	// todo https://github.com/spf13/cobra/issues/661
	commitsCmd.Flags().StringSlice("label", []string{}, "label")
	commitsCmd.Flags().StringSlice("author", []string{}, "author")
	//_ = commitsCmd.MarkFlagRequired("from")
	commitsCmd.Flags().BoolP("update", "u", false, "Update all git repos")
	rootCmd.AddCommand(commitsCmd)

}
