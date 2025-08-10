package usermodel

type deviceToken struct {
	id           int64
	userID       int64
	device_token string
}

func (u *deviceToken) GetID() int64 {
	return u.id
}

func (u *deviceToken) SetID(id int64) {
	u.id = id
}

func (u *deviceToken) GetUserID() int64 {
	return u.userID
}

func (u *deviceToken) SetUserID(userID int64) {
	u.userID = userID
}
func (u *deviceToken) GetDeviceToken() string {
	return u.device_token
}

func (u *deviceToken) SetDeviceToken(token string) {
	u.device_token = token
}
