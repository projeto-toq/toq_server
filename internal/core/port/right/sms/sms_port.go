package smsport

import globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

type SMSPortInterface interface {
	SendSms(notification globalmodel.Notification) error
}
