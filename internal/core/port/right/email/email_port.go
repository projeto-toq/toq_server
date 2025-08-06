package emailport

import globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

type EmailPortInterface interface {
	SendEmail(notification globalmodel.Notification) (err error)
}
