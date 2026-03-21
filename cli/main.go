package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
		fmt.Printf("[flags] System architecture %s is known\n", *arch)
		return nil
	}

	return fmt.Errorf("[flags] System architecture %s is unknown", *arch)
}

func main() {
	location, err := os.Getwd()
	ait.Check(err)

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

	outFileName := filepath.Join(location, "test_fs.squashfs")

	os.Remove(outFileName)
	ait.CreateSquashFSFromFolder(filepath.Join(location, "test.AppDir"), outFileName)
	fileName := ait.DownloadAppImageEngine(ARCH_x86_64)
	ait.AppendToFile(outFileName, fileName)
	hash := ait.CalculateMD5(fileName)
	ait.UpdateMD5(fileName, hash)
	fmt.Println(hex.EncodeToString(hash))
	ait.MakeExecutable(fileName)
	hash = ait.CalculateSha256(fileName)

	// privateKey, err := ait.GeneratePGPPrivateKey("This")
	// fmt.Println(privateKey)
	// signedHash := ait.SignSha256(hash, privateKey, "This")
	// ait.UpdateSha256(fileName, signedHash)
}
