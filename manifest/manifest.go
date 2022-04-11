package manifest

type Manifest struct {
	Name        string `yaml:"name"`
	DisplayName string `yaml:"displayName"`
	RepoUrl     string `yaml:"repoUrl"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Publisher   string `yaml:"publisher"`
}

const MANIFEST_FILE_NAME = "manifest.yml"
