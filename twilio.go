package main

import (
	"fmt"

	"github.com/sfreiberg/gotwilio"
)

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	From       string
	To         []string
}

type Twilio struct {
	Client *gotwilio.Twilio
	From   string
	To     []string
}

func NewTwilio(cfg *TwilioConfig) *Twilio {
	return &Twilio{
		Client: gotwilio.NewTwilioClient(cfg.AccountSID, cfg.AuthToken),
		From:   cfg.From,
		To:     cfg.To,
	}
}

func (t *Twilio) PublishReportChanges(r1, r2 *Report) error {
	msg := fmt.Sprintf("COVID Report - %s:\n%s ICU: %d IN: %d Total: %d\n%s ICU: %d IN: %d Total: %d", r1.Location,
		r1.Time.Format("2006/01/02 03:04PM MST"), r1.ICUUnits, r1.InpatientUnits, r1.Total(),
		r2.Time.Format("2006/01/02 03:04PM MST"), r2.ICUUnits, r2.InpatientUnits, r2.Total())

	for _, to := range t.To {
		_, ex, err := t.Client.SendSMS(t.From, to, msg, "", "")
		if err != nil {
			return fmt.Errorf("send sms error: %w", err)
		}
		if ex != nil {
			return fmt.Errorf("send sms exception: %w", ex)
		}
	}

	return nil
}
