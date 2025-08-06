package smsadapter

import (
	"encoding/json"
	"fmt"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func (s *SmsAdapter) SendSms(notification globalmodel.Notification) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(notification.To)
	params.SetFrom(s.myNumber)
	params.SetBody(notification.Body)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
	return err
}
