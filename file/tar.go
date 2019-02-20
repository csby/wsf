package file

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

type Tar struct {
}

// 从内存解压
// source：待解压二进制数据
// destination：解压后文件所在目录路径
func (s *Tar) DecompressMemory(source []byte, destination string) error {
	br := bytes.NewReader(source)
	gr, err := gzip.NewReader(br)
	if err != nil {
		return err
	}
	reader := tar.NewReader(gr)

	for header, err := reader.Next(); err != io.EOF; header, err = reader.Next() {
		if err != nil {
			return err
		}

		fileInfo := header.FileInfo()
		path := filepath.Join(destination, header.Name)
		if fileInfo.IsDir() {
			os.MkdirAll(path, fileInfo.Mode())
		} else {
			folderPath, err := filepath.Abs(filepath.Dir(path))
			if err == nil {
				os.MkdirAll(folderPath, fileInfo.Mode())
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(f, reader)
			f.Close()

			if err != nil {
				return err
			}
		}
	}

	return nil
}
