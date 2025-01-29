package infrastructure

type RegisterInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
type RegisterEvent struct {
	Topic string
	Body  RegisterInfo
}
