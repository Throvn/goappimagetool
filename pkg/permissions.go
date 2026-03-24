package goappimagetool

import "os"

func MakeExecutable(path string) error {
	return os.Chmod(path, 0o755)
}
