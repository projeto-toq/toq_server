package smsport

import globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"

type SMSPortInterface interface {
	SendSms(notification globalmodel.Notification) error
}
