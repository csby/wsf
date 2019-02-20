package file

import "archive/zip"

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
