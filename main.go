package main

import (
	"context"

	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/packet_forwarder/util"
	"github.com/dotpy3/apartment-alert/feed/kamernet"
)

func main() {
	log := util.GetLogger()
	f := kamernet.NewKamernetFeed(context.Background(), log, "amsterdam")
	for apt := range f.Apartments() {
		log.WithFields(ttnlog.Fields{
			"Address": apt.Address,
			"Price":   apt.Price,
		}).Info("ZOMG NEW APARTMENT")
	}
}
