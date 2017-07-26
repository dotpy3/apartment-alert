package main

import (
	"context"

	"github.com/TheThingsNetwork/packet_forwarder/util"
	"github.com/dotpy3/apartment-alert/alerts"
	"github.com/dotpy3/apartment-alert/alerts/twilio"
	"github.com/dotpy3/apartment-alert/feed/kamernet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Short: "start",
	Long:  "start begins the apartment hunt!",
	Run: func(cmd *cobra.Command, args []string) {
		alerters := make([]alerts.ApartmentAlerter, 0)

		log := util.GetLogger()

		if viper.GetBool("twilio.enable") {
			fromNumber := viper.GetString("twilio.from-number")
			toNumber := viper.GetString("twilio.to-number")
			sid := viper.GetString("twilio.sid")
			token := viper.GetString("twilio.token")

			log.WithField("FromNumber", fromNumber).WithField("ToNumber", toNumber).Info("Enabling Twilio")

			alerters = append(alerters, twilio.Alerter(sid, token, fromNumber, toNumber))
		}

		var cliEnabled bool
		if cliEnabled = viper.GetBool("cli.enable"); cliEnabled {
			log.Info("CLI display enabled")
		} else {
			log.Warn("CLI display disabled")
		}

		f := kamernet.NewKamernetFeed(context.Background(), log, viper.GetString("city"))
		for apt := range f.Apartments() {
			if cliEnabled {
				log.WithField("Address", apt.Address).WithField("Price", apt.Price).Info("New apartment detected")
			}
			for _, alerter := range alerters {
				err := alerter.Push(apt)
				if err != nil {
					log.WithError(err).Error("Couldn't alert of new apartment")
				}
			}
		}

	},
}

func init() {
	startCmd.PersistentFlags().Bool("twilio.enable", false, "Enable cellular alerts")
	startCmd.PersistentFlags().String("twilio.from-number", "", "Number from which Twilio alerts should be sent")
	startCmd.PersistentFlags().String("twilio.to-number", "", "Number to which Twilio alerts should be sent")
	startCmd.PersistentFlags().String("twilio.sid", "", "Twilio Auth SID")
	startCmd.PersistentFlags().String("twilio.token", "", "Twilio Auth token")

	startCmd.PersistentFlags().String("city", "amsterdam", "City in the Netherlands where you want to look for an apartment")

	startCmd.PersistentFlags().Bool("cli.enable", true, "Enable CLI output")

	viper.BindPFlags(startCmd.PersistentFlags())
}
