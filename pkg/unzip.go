package pkg

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	destDir = filepath.Clean(destDir)

	for _, file := range reader.File {
		fmt.Println(file)
		if err := ExtractFile(file, destDir); err != nil {
			return err
		}
	}

	return nil
}

func ExtractFile(file *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, file.Name)
	destPath = filepath.Clean(destPath)
	if !strings.HasPrefix(destPath, destDir) {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(destPath, file.Mode()); err != nil {
			return err
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(srcFile)
		respBytes := buf.String()

		decodeBytes, err := base64.StdEncoding.DecodeString(respBytes)
		if err != nil {
			return err
		}
		reader := bytes.NewReader(decodeBytes)
		defer srcFile.Close()

		if _, err := io.Copy(destFile, reader); err != nil {
			return err
		}
	}

	return nil
}
