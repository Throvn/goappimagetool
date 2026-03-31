package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ait "github.com/Throvn/goappimagetool/pkg/tool"
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
	arch := flag.String("arch", "x86_64", "System Architecture on which the AppImage should run. Valid values are: x86_64, aarch64, i686, armhf")
	runtimePath := flag.String("runtime-file", "", "(Optional) Path of AppImage runtime which is copied into in the AppImage")
	privKeyPath := flag.String("sign-key", "", "(Optional) Path of PGP key file (.asc) to sign the AppImage")
	passphrase := flag.String("passphrase", "", "(Optional) Passphrase of encrypted PGP key file. Only use if encrypted.")

	flag.Parse()

	if flag.NArg() > 0 && flag.Arg(0) == "mkkey" {
		if flag.NArg() != 2 {
			ait.Check(fmt.Errorf("command malformed: use goappimagetool mkdir email@example.com"))
		}
		email := flag.Arg(1)
		emailLocalPart := strings.SplitN(email, "@", 1)[0]

		secretKey, publicKey, err := ait.GenerateSigningKey(emailLocalPart+" - goAppImageTool", email, *passphrase)
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
		appImageEngine = ait.DownloadAppImageEngine(*arch)
	}
	appImageEngine, err := filepath.Abs(appImageEngine)
	ait.Check(err)

	// Only if -sign-key is supplied, load the file into memory.
	var privKeyArmored []byte
	if *privKeyPath != "" {
		privKeyArmored, err = os.ReadFile(*privKeyPath)
		ait.Check(err)
	}

	for i := range cliArgs {
		key := ait.PGPMaterial{
			Passphrase:        *passphrase,
			PrivateKeyArmored: string(privKeyArmored),
		}
		ait.CreateAppImage(cliArgs[i], appImageEngine, key)
	}

}
