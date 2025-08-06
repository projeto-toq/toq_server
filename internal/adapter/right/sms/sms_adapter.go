package smsadapter

import "github.com/twilio/twilio-go"

type SmsAdapter struct {
	client   *twilio.RestClient
	myNumber string
}

func NewSmsAdapter(accountSid string, authToken string, myNumber string) *SmsAdapter {

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	return &SmsAdapter{
		client:   client,
		myNumber: myNumber,
	}
}
