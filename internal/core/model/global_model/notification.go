package globalmodel

type Notification struct {
	To          string
	Title       string
	Body        string
	Icon        string
	Name        string
	DeviceToken string
	Data        map[string]string
}
