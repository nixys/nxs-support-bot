package misc

import (
	"fmt"
	"os/user"
	"strings"
)

// DirNormalize normalizes the directory path
func DirNormalize(path string) (string, error) {

	var err error

	path, err = PathNormalize(path)
	if err != nil {
		return path, fmt.Errorf("dir normalize: %w", err)
	}

	for strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	path += "/"

	return path, nil
}

// PathNormalize normalizes the path
func PathNormalize(path string) (string, error) {

	if strings.HasPrefix(path, "~/") {

		usr, err := user.Current()
		if err != nil {
			return path, fmt.Errorf("path normalize: %w", err)
		}

		path = usr.HomeDir + "/" + strings.TrimPrefix(path, "~/")
	}

	return path, nil
}
