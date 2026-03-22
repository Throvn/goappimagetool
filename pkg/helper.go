package appimagetoolgo

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// Appends the source file to the destination file.
func AppendToFile(src string, dest string) {
	srcFile, err := os.Open(src)
	Check(err)
	defer srcFile.Sync()
	defer srcFile.Close()

	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_APPEND, 0666)
	Check(err)
	defer destFile.Sync()
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	Check(err)
}
