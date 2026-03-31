package goappimagetool

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// If an error supplied it panics printing the error.
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Calculates the size in bytes for a given directory path including all contents.
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

// Used to overwrite sections of elf files.
//
// Path is the location of the elf file.
// Section is the name of the section you want to overwrite (e.g. .sign_key)
// Content is what you want to write into the section
func OverwriteSection(path string, section string, content []byte) error {
	header, err := getSectionHeaderByName(path, section)
	if err != nil {
		return err
	}

	if size := header.GetSize(); size < uint64(len(content)) {
		return fmt.Errorf("%s section has length %d instead of %d", section, size, len(content))
	}
	offset := header.GetFileOffset()

	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Sync()
	defer file.Close()

	bytesWritten, err := file.WriteAt(content, int64(offset))
	if err != nil {
		return err
	}

	if bytesWritten != len(content) {
		return fmt.Errorf("%s was not correctly written", section)
	}

	return nil

}
