package file

import (
	"archive/zip"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"
)

type Zip struct {
}

// 解压文件
// source：待解压文件路径
// destination：解压后文件所在目录路径
func (s *Zip) DecompressFile(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(destination, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 从内存解压
// source：待解压二进制数据
// destination：解压后文件所在目录路径
func (s *Zip) DecompressMemory(source []byte, destination string) error {
	br := bytes.NewReader(source)
	reader, err := zip.NewReader(br, int64(len(source)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		fileName := file.Name
		if file.NonUTF8 {
			if !utf8.ValidString(fileName) {
				transformName, _, err := transform.String(simplifiedchinese.GBK.NewDecoder(), fileName)
				if err == nil {
					fileName = transformName
				}
			}
		}
		path := filepath.Join(destination, fileName)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			folderPath, err := filepath.Abs(filepath.Dir(path))
			if err == nil {
				//os.MkdirAll(folderPath, file.Mode())
				os.MkdirAll(folderPath, 0777)
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				rc.Close()
				return err
			}

			_, err = io.Copy(f, rc)
			f.Close()

			if err != nil {
				rc.Close()
				return err
			}
		}
		rc.Close()
	}

	return nil
}
