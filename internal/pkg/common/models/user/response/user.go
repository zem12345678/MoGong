package response

// LoginResponse 定义登录返回结构体
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	ExpireAt    int64  `json:"expireAt"`
	TimeStamp   int64  `json:"timeStamp"`
}
