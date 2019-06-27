package internal

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"strings"
)

type Commit struct {
	Commit     *object.Commit
	Repository *git.Repository
	Name       string
}

func NewCommit(c *object.Commit, r *git.Repository, name string) (commit *Commit) {
	return &Commit{
		Commit:     c,
		Repository: r,
		Name:       name,
	}
}

func (c Commit) String() string {
	message := strings.Split(c.Commit.Message, "\n")[0]
	if len(message) > 70 {
		message = message[:70]+"..."
	}
	hash := c.Commit.Hash.String()[:8]
	author := c.Commit.Author.Name
	date := c.Commit.Author.When.Format("2006-01-02 15:04")
	name := c.Name

	// Color reference : https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
	return fmt.Sprintf("\033[1;31m[%s\t]\033[0m\033[1;36m[%s]\t\033[0m\033[1;34m[%s]\t\033[0m %s \033[1;32m\t(%s)\033[0m", name, date, hash, message, author)
}

type ByDate []Commit

func (s ByDate) Len() int {
	return len(s)

}
func (s ByDate) Less(i, j int) bool {
	return s[i].Commit.Author.When.Before(s[j].Commit.Author.When)
}

func (s ByDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
