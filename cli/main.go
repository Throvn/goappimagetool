package main

import (
	"encoding/hex"
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

func main() {
	location, err := os.Getwd()
	ait.Check(err)
	outFileName := filepath.Join(location, "test_fs.squashfs")

	os.Remove(outFileName)
	ait.CreateSquashFSFromFolder(filepath.Join(location, "test_fs"), outFileName)
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
