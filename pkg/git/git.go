package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/plumbing/object"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type Repository struct {
	Repo *git.Repository
}

func New() (*Repository, error) {
	var err error
	repo := &Repository{}
	repo.Repo, err = git.PlainOpen(".")

	conf, _ := repo.Repo.ConfigScoped(config.SystemScope)
	fmt.Printf("config: %+v\n", conf)

	return repo, err
}

func (r *Repository) Add(all bool) (bool, error) {
	changes := true
	w, err := r.Repo.Worktree()
	if err != nil {
		return changes, err
	}

	// get the worktree status
	stat, err := w.Status()
	if err != nil {
		return changes, err
	}

	if stat.IsClean() {
		return false, nil
	}

	// if all is true, add all files (also untracked ones)
	if all {
		err = w.AddWithOptions(
			&git.AddOptions{
				All:  all,
				Glob: "*",
			},
		)

		return changes, err
	}

	// add only modified and deleted files
	for file, status := range stat {
		if status.Worktree != git.Untracked {
			_, err = w.Add(file)
			if err != nil {
				return changes, err
			}
		}
	}

	return changes, err
}

func (r *Repository) Commit(msg string) error {
	useSSHSigner := true
	gitConfig, _ := r.Repo.Config()
	if gitConfig.User.Email == "" {
		gitConfig, _ = r.Repo.ConfigScoped(config.GlobalScope)
	}

	fmt.Println(gitConfig.Raw.Section("user").Options.Get("signingkey"))
	if gitConfig.Raw.Section("gpg").Options.Get("format") != "ssh" {
		useSSHSigner = false
	}

	// [core]
	//         repositoryformatversion = 0
	//         filemode = true
	//         bare = false
	//         logallrefupdates = true
	//         ignorecase = true
	//         precomposeunicode = true
	// [remote "origin"]
	//         url = git@github.com:golgoth31/gitcomm.git
	//         fetch = +refs/heads/*:refs/remotes/origin/*
	// [branch "master"]
	//         remote = origin
	//         merge = refs/heads/master
	// [remote "upstream"]
	//         url = git@github.com:karantin2020/gitcomm.git
	//         fetch = +refs/heads/*:refs/remotes/upstream/*
	// [branch "feat/integrate-gptscript"]
	//         github-pr-base-branch = "karantin2020#gitcomm#master"
	//         remote = origin
	//         merge = refs/heads/feat/integrate-gptscript
	//         github-pr-owner-number = "karantin2020#gitcomm#1"
	// [user]
	//         name = David Sabatie
	//         email = david.sabatie@notrenet.com
	//         signingkey = /Users/davidsabatie/.ssh/perso-agicap.pub
	// [gpg]
	//         format = ssh
	// [commit]
	//         gpgsign = true
	w, err := r.Repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitConfig.User.Name,
			Email: gitConfig.User.Email,
			When:  time.Now(),
		},
		// SignKey: gitConfig.Raw.Section("user").Options.Get("signingkey"),
	})
	return nil
}
