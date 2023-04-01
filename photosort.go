package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
)

var (
	acceptedExtensions = []string{".JPG", ".JPEG", ".RAW", ".ARW", ".jpg"}
	numWorkers         = runtime.NumCPU()
)

func main() {
	var source, dest string

	// Prompt the user to enter the source directory, and use the current value as the default
	fmt.Print("Enter the source directory (default: D:/photosort-main/tmp1/): ")
	fmt.Scanln(&source)
	if source == "" {
		source = "D:/photosort-main/tmp1/"
	}

	// Prompt the user to enter the destination directory, and use the current value as the default
	fmt.Print("Enter the destination directory (default: D:/photosort-main/tmp2/): ")
	fmt.Scanln(&dest)
	if dest == "" {
		dest = "D:/photosort-main/tmp2/"
	}

	fileInfos, err := ioutil.ReadDir(source)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var wg sync.WaitGroup
	jobs := make(chan os.FileInfo, numWorkers*2)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fileInfo := range jobs {
				err := copyPhoto(dest, fileInfo, source)
				if err != nil {
					fmt.Println("Error copying photo:", err)
				}
			}
		}()
	}

	for _, fileInfo := range fileInfos {
		jobs <- fileInfo
	}

	close(jobs)
	wg.Wait()
}

func copyPhoto(destDir string, fileInfo os.FileInfo, sourceDir string) error {
	filename := fileInfo.Name()
	extension := strings.ToUpper(filepath.Ext(filename))

	if !stringInSlice(extension, acceptedExtensions) {
		return nil
	}

	fullpath := filepath.Join(sourceDir, filename)

	f, err := os.Open(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return err
	}

	t, err := x.DateTime()
	if err != nil {
		return err
	}

	// date := t.Format("2006/01/02")
	year, month, day := t.Date()

	var destpath string
	if extension == ".JPG" {
		destpath = filepath.Join(destDir, "foto", strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day))
	} else {
		destpath = filepath.Join(destDir, "foto "+extension[1:], strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day))
	}

	if _, err := os.Stat(destpath); os.IsNotExist(err) {
		os.MkdirAll(destpath, 0755)
	}

	destFilepath := filepath.Join(destpath, filename)
	_, err = os.Stat(destFilepath)
	if err == nil {
		fmt.Println(destFilepath, "already exists")
		return nil
	}

	df, err := os.Create(destFilepath)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, f)
	if err != nil {
		return err
	}

	fmt.Println("Copied file", filename, "to", destFilepath)

	return nil
}

func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
