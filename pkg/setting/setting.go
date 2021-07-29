package setting

import (
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting, error) {
	// vp 读取文件时候，设置的文件路径和文件后缀没有添加在文件组成中。需排查
	vp := viper.New()
	vp.SetConfigFile("configs/config.yaml")
	//可以通过不断调用addconfigpath添加配置文件路径，增加容错
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp}, nil
}
