package appimagetoolgo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const DOWNLOAD_URL = "https://github.com/AppImage/type2-runtime/releases/download/continuous"

func downloadFromUrl(url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	output, err := os.Create(fileName)
	if os.IsExist(err) {
		// File already exists, skip downloading.
		return fileName
	} else if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return ""
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}

	fmt.Println(n, "bytes downloaded.")
	return fileName
}

// Downloads the AppImage Engine from the official source
// and returns the location of the downloaded file.
func DownloadAppImageEngine(arch string) string {
	dlLocation := downloadFromUrl(DOWNLOAD_URL + "/runtime-" + arch)
	return dlLocation
}
