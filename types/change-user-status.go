package types

type ChangeUsersStatus struct {
	UserId   string `json:"id" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
	Status   int    `json:"status"`
	Msg      string `json:"msg"`
}

type ChangeUsersStatusDTO struct {
	ChangeUsersStatus []ChangeUsersStatus `json:"users" binding:"required"`
}
