/**
 * @Notes:
 * @class file
 * @package
 * @author: 原始
 * @Time: 2023/6/11   20:15
 */
package helper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetFileName(filePath string, isfex bool) string {
	fileName := path.Base(filePath)
	if isfex {
		return fileName
	}
	// 使用 strings.LastIndex 函数获取最后一个点号的位置
	dotIndex := strings.LastIndex(fileName, ".")
	// 使用字符串切片操作获取不带文件后缀的文件名
	return fileName[:dotIndex]
}

// GetCurrentPath 获取当前的执行路径
func GetCurrentPath() string {
	path, _ := os.Getwd()
	return path
}

// CurrentFile 获取当前文件的详细路径
func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(interface{}("Can not get current file info"))
	}
	return file
}
func Exists(path string, isAll, isNew bool) bool {

	_, err := os.Stat(path) //os.Stat获取文件信息

	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if isNew {
			if isAll {
				err = os.MkdirAll(path, os.ModePerm)
			} else {
				err = os.Mkdir(path, os.ModePerm)
			}
			if err != nil {
				return false
			}
		}
		return false

	}

	return true

}

// 判断所给路径是否为文件夹

func IsDir(path string) bool {

	s, err := os.Stat(path)

	if err != nil {

		return false

	}

	return s.IsDir()

}

// 判断所给路径是否为文件

func IsFile(path string) bool {

	return !IsDir(path)

}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
