package goappimagetool

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/squashfs"
)

func CreateSquashFSFromFolder(srcFolder string, outputFileName string) string {
	// TODO: Explain why we need to set the logical block size and which values should be used
	var LogicalBlocksize diskfs.SectorSize = diskfs.SectorSize4k

	var diskSize int64 = 10 * int64(LogicalBlocksize)

	// Create the disk image
	mydisk, err := diskfs.Create(outputFileName, diskSize, LogicalBlocksize)
	Check(err)

	// Create the ISO filesystem on the disk image
	fspec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeSquashfs,
		VolumeLabel: "label",
	}
	fs, err := mydisk.CreateFilesystem(fspec)
	Check(err)

	// Walk the source folder to copy all files and folders to the SquashFS filesystem
	err = filepath.Walk(srcFolder, func(path string, info os.FileInfo, err error) error {
		Check(err)

		relPath, err := filepath.Rel(srcFolder, path)
		Check(err)

		// If the current path is a folder, create the folder in the ISO filesystem
		if info.IsDir() {
			// Create the directory in the SquashFS file
			err = fs.Mkdir(relPath)
			Check(err)
			return nil
		}

		// If the current path is a file, copy the file to the SquashFS filesystem
		if !info.IsDir() {
			// Open the file in the ISO file for writing
			rw, err := fs.OpenFile(relPath, os.O_CREATE|os.O_RDWR)
			Check(err)
			defer rw.Close()

			// Open the source file for reading
			in, errorOpeningFile := os.Open(path)
			if errorOpeningFile != nil {
				return errorOpeningFile
			}
			defer in.Close()

			// Copy the contents of the source file to the ISO file
			_, err = io.Copy(rw, in)
			Check(err)
		}

		return nil
	})
	Check(err)

	sqfs, ok := fs.(*squashfs.FileSystem)
	if !ok {
		Check(fmt.Errorf("not an squashfs filesystem"))
	}

	err = sqfs.Finalize(squashfs.FinalizeOptions{
		Compression: &squashfs.CompressorZstd{},
	})

	Check(err)

	return outputFileName
}
