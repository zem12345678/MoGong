package request

// LoginRequest 登录请求结构体
type LoginRequest struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
