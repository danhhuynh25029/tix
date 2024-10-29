package pkg

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func ZipFiles(archive string, files ...string) (err error) {
	zipFile, err := os.Create(archive)
	if err != nil {
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, filePath := range files {
		inFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		fileBytes, err := io.ReadAll(inFile)
		if err != nil {
			return err
		}
		outFile, err := os.Create("README.tx")
		if err != nil {
			return err
		}
		outBytes := base64.StdEncoding.EncodeToString(fileBytes)
		if _, err := outFile.Write([]byte(outBytes)); err != nil {
			return err
		}
		_ = outFile.Close()
		err = ProcessFile(zipWriter, filePath)
		if err != nil {
			return err
		}
		//os.Remove("README.tx")
		_ = inFile.Close()
	}
	return
}

func ProcessFile(zipWriter *zip.Writer, filePath string) error {
	var err error
	var fileInfo fs.FileInfo
	var header *zip.FileHeader
	var headerWriter io.Writer
	var file *os.File

	fileInfo, err = os.Stat(filePath)
	if err != nil {
		return err
	}

	header, err = zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Method = zip.Deflate
	header.Name, err = filepath.Rel(filepath.Dir("."), filePath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		header.Name += "/"
	}

	headerWriter, err = zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return nil
	}

	file, err = os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(headerWriter, file)
	return err
}
