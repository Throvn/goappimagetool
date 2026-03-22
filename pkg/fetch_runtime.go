package appimagetoolgo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const DOWNLOAD_URL = "https://github.com/AppImage/type2-runtime/releases/download/continuous"

func downloadFromUrl(url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	output, err := os.Create(fileName)
	if os.IsExist(err) {
		// File already exists, skip downloading.
		return fileName, nil
	} else if err != nil {
		return "", fmt.Errorf("File creation failed for %s - %v", fileName, err)
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Download failed for %s - %v", url, err)
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to write response to file for %s - %v", url, err)
	}

	fmt.Println(n, "bytes downloaded.")
	return fileName, nil
}

// Downloads the AppImage Engine from the official source
// and returns the location of the downloaded file.
func DownloadAppImageEngine(arch string) string {
	dlLocation, err := downloadFromUrl(DOWNLOAD_URL + "/runtime-" + arch)
	Check(err)
	return dlLocation
}
