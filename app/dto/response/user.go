package response

type User struct {
	ID   uint   `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Address string `json:"address"`
	Password string `json:"password"`
}
