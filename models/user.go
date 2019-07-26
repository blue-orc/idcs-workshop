package models

type User struct {
	UserID      int    `json:"userID"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	CreatedOn   string `json:"createdOn"`
	IdcsID      string `json:"idcsID"`
	RoleID      int    `json:"roleID"`
}
