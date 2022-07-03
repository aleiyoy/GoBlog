package upload

import (
	"GoBlog/global"
	"GoBlog/pkg/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

//iota 相当于是一个 const 的常量计数器，你也可以理解为枚举值，第一个声明的 iota 的值为 0，在新的一行被使用时，它的值都会自动递增。
const (
	TypeImage FileType = iota + 1
	TypeExcel
	TypeTxt
)

func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

func GetFileExt(name string) string {
	return path.Ext(name)
}

func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}
// ################################################################################################
// 对相关文件信息进行检查
func CheckSavePath(dst string) bool {
	// 检查保存的路径
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

func CheckContainExt(t FileType, name string) bool {
	// 检查文件后缀是否在允许的文件后缀返回内：
	//		- .jpg
	//      - .jpeg
	//      - .png
	ext := GetFileExt(name)
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}

	}

	return false
}

func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}

	return false
}


func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

// ################################################################################################
// 文件保存相关
func CreateSavePath(dst string, perm os.FileMode) error {
	//创建在上传文件时所使用的保存目录，在方法内部调用的 os.MkdirAll 方法，
	//该方法将会以传入的 os.FileMode 权限位去递归创建所需的所有目录结构，
	//若涉及的目录均已存在，则不会进行任何操作，直接返回 nil。
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}

	return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	// 调用 os.Create 方法创建目标地址的文件，
	//再通过 file.Open 方法打开源地址的文件，
	//结合 io.Copy 方法实现两者之间的文件内容拷贝。
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}