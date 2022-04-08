package git

import (
	"os"
	"strings"
	"time"

	"github.com/borgmon/openpilot-mod-manager/common"
	"github.com/ldez/go-git-cmd-wrapper/branch"
	"github.com/ldez/go-git-cmd-wrapper/checkout"
	"github.com/ldez/go-git-cmd-wrapper/clone"
	"github.com/ldez/go-git-cmd-wrapper/commit"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/reset"
	"github.com/ldez/go-git-cmd-wrapper/status"
	"github.com/pkg/errors"
)

type GitHandlerImpl struct{}

var gitHandlerInstance GitHandler

func GetGitHandler() GitHandler {
	if gitHandlerInstance != nil {
		return gitHandlerInstance
	}
	gitHandlerInstance = &GitHandlerImpl{}
	return gitHandlerInstance
}

func (handler *GitHandlerImpl) Clone(path, url string) error {
	err := os.Chdir(path)
	if err != nil {
		return errors.WithStack(err)
	}
	out, err := git.Clone(clone.Repository(url), git.Debug)
	common.LogIfVarbose(out)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) GetBranchName(gitPath string) (string, error) {
	err := os.Chdir(gitPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	out, err := git.Status(status.Short, status.Branch, git.Debug)
	common.LogIfVarbose(out)
	if err != nil {
		return "", errors.WithStack(err)
	}
	out = strings.Split(out, "\n")[0]
	return out[3:], nil
}

func (handler *GitHandlerImpl) NewBranch(gitPath string, name string) error {
	err := os.Chdir(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	out, err := git.Checkout(checkout.NewBranch(handler.GenerateBranchName()), git.Debug)
	common.LogIfVarbose(out)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) RemoveBranch(gitPath string, name string) error {
	err := os.Chdir(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	out, err := git.Branch(branch.DeleteForce, branch.BranchName(handler.GenerateBranchName()), git.Debug)
	common.LogIfVarbose(out)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (handler *GitHandlerImpl) CheckoutBranch(gitPath string, name string) error {
	err := os.Chdir(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	out, err := git.Checkout(checkout.Branch(name), git.Debug)
	common.LogIfVarbose(out)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (handler *GitHandlerImpl) GenerateBranchName() string {
	now := time.Now()
	return "omm-" + now.Format("2006-01-02-3-4-5-pm")
}

func (handler *GitHandlerImpl) CommitBranch(gitPath string, name string) error {
	err := os.Chdir(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	out, err := git.Commit(commit.Amend, commit.Message(name), commit.AllowEmpty, git.Debug)
	common.LogIfVarbose(out)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil

}

func (handler *GitHandlerImpl) ResetBranch(gitPath string) error {
	err := os.Chdir(gitPath)
	if err != nil {
		return errors.WithStack(err)
	}
	out, err := git.Reset(reset.Hard, git.Debug)
	common.LogIfVarbose(out)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
