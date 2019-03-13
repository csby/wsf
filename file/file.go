package file

import (
	"archive/zip"
	"os"
)

// 从内存解压tar或zip格式文件
// src：待解压二进制数据
// dest：解压后文件所在目录路径
func Decompress(src []byte, dest string) error {
	z := &Zip{}
	err := z.DecompressMemory(src, dest)
	if err == nil {
		return err
	}

	if err.Error() == zip.ErrFormat.Error() {
		tar := &Tar{}
		err = tar.DecompressMemory(src, dest)
	}

	return err
}

// 拷贝文件(夹)
// src： 源文件(夹)路径
// dest： 目标文件夹路径
func Copy(src, dest string) error {
	folder := &Folder{}

	return folder.Copy(src, dest)
}

// 判断文件(夹)是否存在
func Exist(name string) bool {
	_, err := os.Stat(name)

	return err == nil || os.IsExist(err)
}

// 删除文件或空文件夹
func Delete(name string) error {
	_, err := os.Stat(name)
	if err == nil || os.IsExist(err) {
		err = os.Remove(name)
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除文件夹
func DeleteAll(name string) error {
	return os.RemoveAll(name)
}
