package handlers

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamSignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamJoinGroup struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type ParamNewGroup struct {
	GroupName   string `json:"group_name" binding:"required"`
	Description string `json:"description"`
}

type ParamFindGroup struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type ParamLeaveGroup struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type ParamPullMessage struct {
	LastSeq  int64 `json:"last_seq"`
	TargetID int64 `json:"target_id"`
	Size     int   `json:"size" binding:"required"`
}
