package git

type GitHandler interface {
	Clone(string) error
	GetBranch() (string, error)
	NewBranch(string) error
	RemoveBranch(string) error
	CheckoutBranch(string) error
}
