package git

type GitHandler interface {
	Clone(path, url string) error
	GetBranchName(gitPath string) (string, error)
	NewBranch(gitPath string, name string) error
	RemoveBranch(gitPath string, name string) error
	CheckoutBranch(gitPath string, name string) error
	GenerateBranchName() string
	CommitBranch(gitPath string, name string) error
	AddBranch(gitPath string) error
	ResetBranch(gitPath string) error
	ListBranch(gitPath string) (string, error)
}
