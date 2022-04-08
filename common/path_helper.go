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

func GetUserNameFromGithub(urlStr string) (string, error) {
	parts, err := ParseGithubUrl(urlStr)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return (*parts)[0], nil
}

func GetProjectFromGithub(urlStr string) (string, error) {
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
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
