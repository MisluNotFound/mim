package handlers

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Avatar     string `json:"avatar"`
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
	LastSeq   int64 `json:"last_seq"`
	SessionID int64 `json:"target_id" binding:"required"`
	Size      int   `json:"size" binding:"required"`
	IsGroup   bool  `json:"is_group"`
}

type ParamPullOfflineMessage struct {
	IsGroup   bool  `json:"is_group"`
	SessionID int64 `json:"session_id" binding:"required"`
}

type ParamNearby struct {
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

type ParamAddFriend struct {
	FriendID int64 `json:"friend_id" binding:"required"`
}

type ParamRemoveFriend struct {
	FriendID int64 `json:"friend_id" binding:"required"`
}

type ParamUpdatePhoto struct {
	Avatar string `json:"avatar" binding:"required"`
}

type ParamUpdatePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	RePassword  string `json:"re_password" binding:"required,eqfield=NewPassword"`
}

type ParamUpdateName struct {
	Name string `json:"name"`
}

type ParamUpdateFriendRemark struct {
	FriendID int64  `json:"friend_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type ParamFindFriend struct {
	UserID int64 `json:"user_id" binding:"required"`
}
