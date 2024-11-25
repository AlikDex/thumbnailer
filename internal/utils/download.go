package utils

import (
	"app/pkg/filesystem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	activeDownloads sync.Map
)

func DownloadFile(url string, outFilepath string) error {
	if _, loaded := activeDownloads.LoadOrStore(url, true); loaded {
		return fmt.Errorf("already downloading")
	}

	defer activeDownloads.Delete(url)

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status should be 200, got: %d", resp.StatusCode)
	}

	if resp.ContentLength == 0 {
		return fmt.Errorf("content is empty")
	}

	if !filesystem.IsDir(filepath.Dir(outFilepath)) {
		filesystem.MkDir(filepath.Dir(outFilepath))
	}

	tmpFilePath := outFilepath + ".tmp"
	tmpFile, err := os.Create(tmpFilePath)

	if err != nil {
		return err
	}

	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)

	if err != nil {
		os.Remove(tmpFile.Name())

		return err
	}

	err = os.Rename(tmpFile.Name(), outFilepath)

	if err != nil {
		os.Remove(tmpFile.Name())

		return err
	}

	lastModified := resp.Header.Get("Last-Modified")
	modifiedTime := time.Now()

	if lastModified != "" {
		modifiedTime, err = time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", lastModified)

		if err != nil {
			modifiedTime = time.Now()
		}
	}

	return os.Chtimes(outFilepath, time.Now(), modifiedTime)
}
