package common

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func GetPathFromFilePath(filePath string) string {
	return filepath.Dir(filePath)
}

func GetFileFromFilePath(filePath string) string {
	return filepath.Base(filePath)
}

func GetUserFromGithub(urlStr string) (string, error) {
	parts, err := ParseGithubUrl(urlStr)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return (*parts)[0], nil
}

func GetNameFromGithub(urlStr string) (string, error) {
	parts, err := ParseGithubUrl(urlStr)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return (*parts)[1], nil
}

func ParseGithubUrl(urlStr string) (*[]string, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	path := u.Path

	parts := strings.Split(path, "/")

	return &[]string{parts[len(parts)-2], parts[len(parts)-1]}, nil
}

func IsUrl(urlStr string) bool {
	return strings.Contains(urlStr, "http")
}
