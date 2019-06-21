package internal

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"os/user"
)

type Repository struct {
	Url       string
	Labels    []string
	Name      string
	LocalPath string
}

func (r Repository) ListCommits(filter Filter) (commits []Commit, e error) {

	gitRepo, err := git.PlainOpen(r.LocalPath)
	if err != nil {
		return nil, fmt.Errorf("%v\n", err)
	}

	ref, err := gitRepo.Head()
	if err != nil {
		return nil, fmt.Errorf("%v\n", err)
	}

	commitIter, err := gitRepo.Log(&git.LogOptions{
		From:  ref.Hash(),
		All:   true,
		Order: git.LogOrderCommitterTime,
	})

	if err != nil {
		return nil, fmt.Errorf("%v\n", err)
	}

	err = commitIter.ForEach(func(c *object.Commit) error {
		if filter.Filter(c) {
			commits = append(commits, *NewCommit(c, gitRepo, r.Name))
		}
		return nil
	})

	return commits, nil
}

func (r Repository) SyncRepo() error {

	fmt.Println("Syncing " + r.Name + "...")

	currentUser, _ := user.Current()

	sshAuth, _ := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")

	// Cloning repository
	repo, err := git.PlainClone(r.LocalPath, false, &git.CloneOptions{
		URL:  r.Url,
		Auth: sshAuth,
	})

	switch err {
	case git.ErrRepositoryAlreadyExists:
		repo, err = git.PlainOpen(r.LocalPath)
		if err != nil {
			return fmt.Errorf("clone error: %v\n", err)
		}
		break
	case nil:
		break
	default:
		return fmt.Errorf("clone error: %v\n", err)
	}

	//Fetching all branches
	remote, err := repo.Remote("origin")

	if err != nil {
		return fmt.Errorf("%v\n", err)
	}
	fetchOptions := &git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
		Auth:     sshAuth,
	}
	if err := remote.Fetch(fetchOptions); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("fetch error : %v\n", err)
	}

	// Pulling all branches
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}
	w.Pull(&git.PullOptions{RemoteName: "origin"})

	return nil
}
