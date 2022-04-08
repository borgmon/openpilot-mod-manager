package git

type GitHandler interface {
	Clone(url string) error
	GetBranchName(gitPath string) (string, error)
	NewBranch(gitPath string, name string) error
	RemoveBranch(gitPath string, name string) error
	CheckoutBranch(gitPath string, name string) error
	GenerateBranchName() string
	CommitBranch(gitPath string, name string) error
	ResetBranch(gitPath string) error
}
