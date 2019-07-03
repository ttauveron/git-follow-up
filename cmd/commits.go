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
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

var filter *git.Filter

// commitsCmd represents the commits command
var commitsCmd = &cobra.Command{
	Use:   "commits",
	Short: "Get list of commits from your tracked repositories",
	Run: func(cmd *cobra.Command, args []string) {
		filter = git.NewFilter(cmd.Flags())

		// Sync repos if update flag is provided
		doUpdate, err := cmd.Flags().GetBool("update")
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		if doUpdate {
			updateCmd.Run(cmd, args)
		}

		var commits []git.Commit

		// Listing log messages of repositories
		for _, repo := range config.Repositories {
			// Skip update on non-matching labels
			if cmd.Flags().Changed("label") && !git.ContainsAll(repo.Labels, filter.Labels) {
				continue
			}
			cs, err := repo.ListCommits(*filter)
			commits = append(commits, cs...)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}

		sort.Sort(git.ByDate(commits))

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

func formatCommit(c git.Commit) (result string) {
	message := strings.Split(c.Commit.Message, "\n")[0]
	if len(message) > 70 {
		message = message[:70] + "..."
	}
	hash := c.Commit.Hash.String()[:8]
	author := c.Commit.Author.Name
	date := c.Commit.Author.When.Format("2006-01-02 15:04")
	name := c.Name

	// Color reference : https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
	if git.Contains(filter.Display, "repo") {
		result += "\033[1;31m" + name + "\t \033[0m"
	}

	if git.Contains(filter.Display, "date") {
		result += "\033[1;36m" + date + "\t \033[0m"
	}

	if git.Contains(filter.Display, "hash") {
		result += "\033[1;34m" + hash + "\t\033[0m"
	}

	if git.Contains(filter.Display, "message") {
		result += " " + message + " \t"
	}

	if git.Contains(filter.Display, "author") {
		result += "\033[1;32m" + author + "\033[0m"
	}

	return result
}

func init() {

	commitsCmd.Flags().StringSlice("label", []string{}, "filters by project labels")
	commitsCmd.Flags().StringSlice("author", []string{}, "filters by authors")

	commitsCmd.Flags().String("from", "wtd", "filters commit by date (ytd, mtd, wtd, yesterday, today, [yyyy-MM-dd])")
	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__from_values"}
	flag := commitsCmd.Flags().Lookup("from")
	flag.Annotations = annotation

	annotation = make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__display_values"}
	commitsCmd.Flags().StringSlice("display", []string{}, "fields to be displayed")
	flag = commitsCmd.Flags().Lookup("display")
	flag.Annotations = annotation

	commitsCmd.Flags().BoolP("update", "u", false, "synchronizes git repositories")
	rootCmd.AddCommand(commitsCmd)

}
