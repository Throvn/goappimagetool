package goappimagetool

import "os"

// Sets the permission bits for the given file to 755.
func MakeExecutable(path string) error {
	return os.Chmod(path, 0o755)
}
