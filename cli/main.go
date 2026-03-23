package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
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
		return nil
	}

	return fmt.Errorf("[flags] System architecture %s is unknown", *arch)
}

func main() {
	// Args for default (mkappdir) command
	arch := flag.String("-arch", "x86_64", "System Architecture on which the AppImage should run. Valid values are: x86_64, aarch64, i686, armhf")
	runtimePath := flag.String("-runtime-file", "", "(Optional) Path of AppImage runtime which is copied into in the AppImage")
	privKeyPath := flag.String("-sign-key", "", "(Optional) Path of PGP key file (.asc) to sign the AppImage")
	passphrase := flag.String("-passphrase", "", "(Optional) Passphrase of encrypted PGP key file. Only use if encrypted.")

	flag.Parse()

	if flag.NArg() > 0 && flag.Arg(0) == "mkkey" {
		if flag.NArg() != 2 {
			ait.Check(fmt.Errorf("command malformed: use appimagetool.go mkdir email@example.com"))
		}
		// TODO: Generate new key and write it to cwd.
		currUser, err := user.Current()
		ait.Check(err)

		secretKey, publicKey, err := ait.GenerateSigningKey(currUser.Username+" - AppImageTool.go", flag.Arg(1), *passphrase)
		ait.Check(err)

		err = os.WriteFile("private.asc", []byte(secretKey), 0o400)
		ait.Check(err)
		os.WriteFile("public.asc", []byte(publicKey), 0o400)
		ait.Check(err)

		return
	}

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
	appImageEngine, err := filepath.Abs(appImageEngine)
	ait.Check(err)

	for i := range cliArgs {
		ait.CreateAppImage(cliArgs[i], appImageEngine)
	}

	// privateKey, err := ait.GeneratePGPPrivateKey("This")
	// fmt.Println(privateKey)
	// signedHash := ait.SignSha256(hash, privateKey, "This")
	// ait.UpdateSha256(fileName, signedHash)
}
