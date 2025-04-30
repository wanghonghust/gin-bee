package request

import "gin-bee/apps/system/model"

type CreateUserParam struct {
	Username    string        `json:"username" binding:"required"`
	Password    string        `json:"password" binding:"required"`
	PasswordC   string        `json:"password_c" binding:"required"`
	Nickname    string        `json:"nickname"`
	Email       string        `json:"email" binding:"required"`
	Role        []uint        `json:"role"`
	IsSuperUser bool          `json:"isSuperUser"`
	Limiter     model.Limiter `json:"limiter"`
}

type UpdateUserParam struct {
	ID          uint          `json:"id"`
	Nickname    string        `json:"nickname"`
	Email       string        `json:"email"`
	Role        []uint        `json:"role"`
	IsSuperUser bool          `json:"isSuperUser"`
	Limiter     model.Limiter `json:"limiter"`
}
