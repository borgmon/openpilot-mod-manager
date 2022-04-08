package git

import (
	"fmt"
	"os"
	"time"

	"github.com/ldez/go-git-cmd-wrapper/branch"
	"github.com/ldez/go-git-cmd-wrapper/checkout"
	"github.com/ldez/go-git-cmd-wrapper/clone"
	"github.com/ldez/go-git-cmd-wrapper/commit"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/reset"
	"github.com/pkg/errors"
)

type GitHandlerImpl struct{}

func GetGitHandler() GitHandler {
	return &GitHandlerImpl{}
}

func (handler *GitHandlerImpl) Clone(path, url string) error {
	os.Chdir(path)
	out, err := git.Clone(clone.Repository(url), git.Debug)
	fmt.Println(out)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) GetBranchName(gitPath string) (string, error) {
	os.Chdir(gitPath)
	head, err := git.Raw("branch --show-current", git.Debug)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return head, nil
}

func (handler *GitHandlerImpl) NewBranch(gitPath string, name string) error {
	os.Chdir(gitPath)
	out, err := git.Checkout(checkout.NewBranch(handler.GenerateBranchName()), git.Debug)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(out)
	return nil
}

func (handler *GitHandlerImpl) RemoveBranch(gitPath string, name string) error {
	os.Chdir(gitPath)
	out, err := git.Branch(branch.DeleteForce, branch.BranchName(handler.GenerateBranchName()), git.Debug)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(out)

	return nil
}

func (handler *GitHandlerImpl) CheckoutBranch(gitPath string, name string) error {
	os.Chdir(gitPath)
	out, err := git.Checkout(checkout.Branch(name), git.Debug)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(out)
	return nil
}

func (handler *GitHandlerImpl) GenerateBranchName() string {
	now := time.Now()
	return "omm-" + now.Format("2006-01-02-3-4-5-pm")
}

func (handler *GitHandlerImpl) CommitBranch(gitPath string, name string) error {
	os.Chdir(gitPath)
	out, err := git.Commit(commit.Amend, commit.Message(name), commit.AllowEmpty, git.Debug)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(out)
	return nil

}

func (handler *GitHandlerImpl) ResetBranch(gitPath string) error {
	os.Chdir(gitPath)
	out, err := git.Reset(reset.Hard, git.Debug)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(out)
	return nil
}
