package upload

import (
	"github.com/fishblog/global"
	"github.com/fishblog/pkg/util"
	"os"
	"path"
	"strings"
)

type FileType int

const TypeImage FileType = iota + 1

func GetFileName(name string) string {
	suffix := GetFileSuffix(name)
	fileName := strings.TrimSuffix(name, suffix)
	fileName = util.EncodeMD5(fileName)
	return fileName + suffix
}

func GetFileSuffix(name string) string {
	return path.Ext(name)
}

func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}


