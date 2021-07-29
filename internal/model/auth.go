package model

type Auth struct {
	*Model
	AppKey    string `json:"app_key"`
	AppSecret string `app_secret`
}
