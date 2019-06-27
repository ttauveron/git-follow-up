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
	"git-follow-up/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
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

		// initialize tabwriter
		w := new(tabwriter.Writer)
		defer w.Flush()

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, ' ', 0)
		for _, commit := range commits {
			fmt.Fprintf(w, formatCommit(commit)+"\n")
		}

	},
}

func formatCommit(c internal.Commit) (result string) {
	message := strings.Split(c.Commit.Message, "\n")[0]
	if len(message) > 70 {
		message = message[:70] + "..."
	}
	hash := c.Commit.Hash.String()[:8]
	author := c.Commit.Author.Name
	date := c.Commit.Author.When.Format("2006-01-02 15:04")
	name := c.Name

	// Color reference : https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
	if internal.Contains(filter.Display, "repo") {
		result += "\033[1;31m[" + name + "\t]\033[0m"
	}

	if internal.Contains(filter.Display, "date") {
		result += "\033[1;36m[" + date + "]\t\033[0m"
	}

	if internal.Contains(filter.Display, "hash") {
		result += "\033[1;34m[" + hash + "]\t\033[0m"
	}

	if internal.Contains(filter.Display, "message") {
		result += " " + message + " "
	}

	if internal.Contains(filter.Display, "author") {
		result += "\033[1;32m\t(" + author + ")\033[0m"
	}

	return result
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
	commitsCmd.Flags().StringSlice("label", []string{}, "label")
	commitsCmd.Flags().StringSlice("author", []string{}, "author")
	commitsCmd.Flags().StringSlice("display", []string{}, "display")
	//_ = commitsCmd.MarkFlagRequired("from")
	commitsCmd.Flags().BoolP("update", "u", false, "Update all git repos")
	rootCmd.AddCommand(commitsCmd)

}
