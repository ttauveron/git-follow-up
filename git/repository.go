package git

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"os/user"
)

type Repository struct {
	Url            string
	Labels         []string
	Name           string
	LocalPath      string
	Authentication Authentication
}

type Authentication struct {
	Type     string
	AuthFile string `mapstructure:"auth_file"`
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

	var auth transport.AuthMethod

	switch r.Authentication.Type {
	case "ssh":
		if r.Authentication.AuthFile == "" {
			currentUser, _ := user.Current()
			auth, _ = ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")
		} else {
			auth, _ = ssh.NewPublicKeysFromFile("git", r.Authentication.AuthFile, "")
		}
		break
	case "access_token":
		accessToken, err := ioutil.ReadFile(r.Authentication.AuthFile)
		if err != nil {
			return fmt.Errorf("%v : auth file error: %v\n", r.Name, err)
		}
		auth = &http.BasicAuth{
			Username: "anything",
			Password: string(accessToken),
		}
		break
	default:
		auth = nil
		break
	}

	// Cloning repository
	repo, err := git.PlainClone(r.LocalPath, true, &git.CloneOptions{
		URL:        r.Url,
		Auth:       auth,
		NoCheckout: true,
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
		Auth:     auth,
	}
	if err := remote.Fetch(fetchOptions); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("fetching %v : %v\n", r.Name, err)
	}

	return nil
}
