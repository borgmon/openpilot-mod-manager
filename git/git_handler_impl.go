package git

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

type GitHandlerImpl struct {
	CachePath      string
	OriginalBranch string
}

func NewGitHandler(CachePath string, OriginalBranch string) GitHandler {
	return &GitHandlerImpl{
		CachePath:      CachePath,
		OriginalBranch: OriginalBranch,
	}
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

func (handler *GitHandlerImpl) GetBranchName(gitPath string) (string, error) {
	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	head, err := repo.Head()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return head.String(), nil
}

func (handler *GitHandlerImpl) NewBranch(gitPath string, name string) error {
	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.CreateBranch(&config.Branch{Name: name})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) RemoveBranch(gitPath string, name string) error {
	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = repo.DeleteBranch(name)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) CheckoutBranch(gitPath string, name string) error {
	tree, err := handler.GetWorkTree(gitPath)
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

func (handler *GitHandlerImpl) GenerateBranchName() string {
	now := time.Now()
	return "omm-" + now.Format(time.RFC3339Nano)
}
func (handler *GitHandlerImpl) GetWorkTree(gitPath string) (*git.Worktree, error) {
	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tree, err := repo.Worktree()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tree, nil
}

func (handler *GitHandlerImpl) CommitBranch(gitPath string, name string) error {
	tree, err := handler.GetWorkTree(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = tree.Add(".")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = tree.Commit(name, &git.CommitOptions{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) ResetBranch(gitPath string) error {
	tree, err := handler.GetWorkTree(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = tree.Reset(nil)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
