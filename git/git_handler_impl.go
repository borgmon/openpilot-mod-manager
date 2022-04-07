package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

type GitHandlerImpl struct {
	OriginalBranch string
	CachePath      string
	CurrentPath    string
}

func (handler *GitHandlerImpl) Clone(url string) error {
	_, err := git.PlainClone(handler.CachePath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) GetBranch() (string, error) {
	repo, err := git.PlainOpen(handler.CurrentPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	head, err := repo.Head()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return head.String(), nil
}

func (handler *GitHandlerImpl) NewBranch(name string) error {
	repo, err := git.PlainOpen(handler.CurrentPath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.CreateBranch(&config.Branch{Name: name})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) RemoveBranch(name string) error {
	repo, err := git.PlainOpen(handler.CurrentPath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DeleteBranch(name)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) CheckoutBranch(name string) error {
	repo, err := git.PlainOpen(handler.CurrentPath)
	if err != nil {
		return errors.WithStack(err)
	}
	tree, err := repo.Worktree()
	if err != nil {
		return errors.WithStack(err)
	}
	err = tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
