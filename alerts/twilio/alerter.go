package twilio

import (
	"fmt"

	"github.com/dotpy3/apartment-alert/alerts"
	"github.com/dotpy3/apartment-alert/feed"
	"github.com/pkg/errors"
	"github.com/sfreiberg/gotwilio"
)

type twilioAlerter struct {
	cli        *gotwilio.Twilio
	fromNumber string
	toNumber   string
}

func Alerter(sid, token, fromNumber, toNumber string) alerts.ApartmentAlerter {
	return &twilioAlerter{
		cli:        gotwilio.NewTwilioClient(sid, token),
		fromNumber: fromNumber,
		toNumber:   toNumber,
	}
}

func (a *twilioAlerter) Push(apt feed.Apartment) error {
	_, _, err := a.cli.SendSMS(a.fromNumber, a.toNumber,
		fmt.Sprintf("A new apartment is available at %s (%s), for a price of %s!", apt.Address, apt.Postcode, apt.Price),
		"", "")
	return errors.Wrap(err, "couldn't send SMS")
}
