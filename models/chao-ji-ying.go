package models

type ChaoJiYing struct {
	ErrNo  int    `json:"err_no"`
	ErrStr string `json:"err_str"`
	PicID  string `json:"pic_id"`
	PicStr string `json:"pic_str"`
	MD5    string `json:"md5"`
}
