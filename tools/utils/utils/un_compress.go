package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"dataProcess/constants"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnZipFile(zipPath, zipFileName string) (err error) {
	zipReader, err := zip.OpenReader(zipPath + string(os.PathSeparator) + zipFileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	cond := false
	for _, file := range zipReader.File {
		if preCheckRunEntry(file.Name) {
			cond = true
			break
		}
	}

	if !cond {
		err = errors.New("this compress file not contain run.sh")
		return
	}

	var inFile io.ReadCloser
	var outFile *os.File

	for _, file := range zipReader.File {

		fPath := zipPath + string(os.PathSeparator) + file.Name

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(fPath, os.ModePerm)
			if err != nil {
				return
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
				return
			}

			inFile, err = file.Open()
			if err != nil {
				return
			}

			outFile, err = os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return
			}

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}

			inFile.Close()
			outFile.Close()
		}
	}

	return
}

func UnTarFile(tarPath, tarFileName string) (err error) {

	fr, err := os.Open(tarPath + string(os.PathSeparator) + tarFileName)
	if err != nil {
		return err
	}
	defer fr.Close()

	tr := tar.NewReader(fr)

	var h *tar.Header
	var outFile *os.File

	cond := false

	for {
		h, err = tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fPath := tarPath + string(os.PathSeparator) + h.Name

		if h.FileInfo().IsDir() {
			err = os.MkdirAll(fPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if preCheckRunEntry(h.Name) {
				cond = true
			}
			outFile, err = os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, tr)
			if err != nil {
				return err
			}

			outFile.Close()
		}
	}

	if !cond {
		err = errors.New("this compress file not contain run.sh")
		return err
	}

	return nil
}

func UnTarGzFile(tarGzPath, tarGzFileName string) (err error) {
	fr, err := os.Open(tarGzPath + string(os.PathSeparator) + tarGzFileName)
	if err != nil {
		return
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	if err != nil {
		return
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	var head *tar.Header
	var outFile *os.File

	cond := false

	for {
		head, err = tr.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				return
			}
		}

		if head == nil {
			break
		}

		fPath := tarGzPath + string(os.PathSeparator) + head.Name

		if head.FileInfo().IsDir() {
			err = os.MkdirAll(fPath, os.ModePerm)
			if err != nil {
				return
			}
		} else {
			if preCheckRunEntry(head.Name) {
				cond = true
			}
			outFile, err = os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return
			}

			_, err = io.Copy(outFile, tr)
			if err != nil {
				return
			}

			outFile.Close()
		}
	}

	if !cond {
		err = errors.New("this compress file not contain run.sh")
		return err
	}

	return
}

func preCheckRunEntry(name string) (cond bool) {
	return strings.HasSuffix(name, constants.RunEntry)
}
