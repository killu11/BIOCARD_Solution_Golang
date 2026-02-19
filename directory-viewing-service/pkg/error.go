package pkg

import (
	"fmt"
	"strings"
)

func PackageError(pkgName, msg string, err error) error {
	return fmt.Errorf("%s: %s: %v", pkgName, msg, err)
}

func ErrDescription(err error) string {
	parts := strings.Split(err.Error(), ":")
	if len(parts) == 3 {
		return parts[2]
	}

	return err.Error()
}
