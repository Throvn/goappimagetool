package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ait "github.com/Throvn/appimagetool.go"
)

const (
	ARCH_x86_64  = "x86_64"
	ARCH_aarch64 = "aarch64"
	ARCH_i686    = "i686"
	ARCH_armhf   = "armhf"
)

func checkArch(arch *string) error {
	switch *arch {
	case ARCH_x86_64, ARCH_aarch64, ARCH_i686, ARCH_armhf:
		return nil
	}

	return fmt.Errorf("[flags] System architecture %s is unknown", *arch)
}

// Taken from https://stackoverflow.com/a/28371044/10408987
func copyFile(src string, dst string) {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	ait.Check(err)
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	ait.Check(err)
}

func safeFileBase(path string) string {
	dir, file := filepath.Split(path)

	// Remove file extension, if exists.
	fileParts := strings.Split(file, ".")
	if len(fileParts) > 1 {
		file = strings.Join(fileParts[:len(fileParts)-1], ".")
	}

	return filepath.Join(dir, file)
}

func createAppImage(path string, appImageEngine string) {

	fileName := safeFileBase(path) + ".AppImage"

	copyFile(appImageEngine, fileName)

	location, err := os.Getwd()
	ait.Check(err)

	outFileName := safeFileBase(path) + ".squashfs"

	os.Remove(outFileName)
	ait.CreateSquashFSFromFolder(filepath.Join(location, path), outFileName)

	ait.AppendToFile(outFileName, fileName)
	hash := ait.CalculateMD5(fileName)
	ait.UpdateMD5(fileName, hash)
	fmt.Println(hex.EncodeToString(hash))
	ait.MakeExecutable(fileName)
	hash = ait.CalculateSha256(fileName)

	err = os.Remove(outFileName)
	ait.Check(err)

	fmt.Printf("Created %s\n", fileName)
}

func main() {
	arch := flag.String("arch", "x86_64", "System Architecture on which the AppImage should run. Valid values are: x86_64, aarch64, i686, armhf")
	runtimePath := flag.String("runtime-file", "", "(Optional) Path of AppImage runtime which BECOMES the AppImage")
	privKeyPath := flag.String("sign-key", "", "(Optional) Path of PGP private key file to sign the AppImage")
	passphrase := flag.String("passphrase", "", "(Optional) Passphrase of encrypted PGP key file. Only use if encrypted.")
	flag.Parse()

	if err := checkArch(arch); err != nil {
		if *runtimePath == "" {
			ait.Check(fmt.Errorf("Unknown system architecture supplied. Supply -runtime-file or choose from (x86_64, aarch64, i686, armhf)"))
		}
	}

	if *privKeyPath != "" || *passphrase != "" {
		fmt.Println("[flags] Warning: Code signing is not yet implemented")
	}

	cliArgs := flag.Args()
	if len(cliArgs) <= 0 {
		ait.Check(fmt.Errorf("[flags] No AppDir supplied"))
	}

	fmt.Printf("Args: %s\n", cliArgs[0])

	// Use -runtime-file as engine blueprint if given.
	// Otherwise download the AppImage engine from the official source.
	appImageEngine := *runtimePath
	if appImageEngine == "" {
		appImageEngine = ait.DownloadAppImageEngine(ARCH_x86_64)
	}

	for i := range cliArgs {
		createAppImage(cliArgs[i], appImageEngine)
	}

	// privateKey, err := ait.GeneratePGPPrivateKey("This")
	// fmt.Println(privateKey)
	// signedHash := ait.SignSha256(hash, privateKey, "This")
	// ait.UpdateSha256(fileName, signedHash)
}
