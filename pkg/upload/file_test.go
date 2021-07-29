package upload

import (
	"github.com/fishblog/global"
	"testing"
)

func TestCheckSavePath(t *testing.T) {
	if !CheckSavePath(global.AppSetting.UploadSavePath) {
		t.Error("open file error")
	}
}
