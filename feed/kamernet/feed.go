package kamernet

import (
	"context"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/dotpy3/apartment-alert/feed"
	"github.com/pkg/errors"
)

func NewKamernetFeed(ctx context.Context, logCtx log.Interface, city string) feed.Feed {
	apartments := make(chan feed.Apartment)
	k := &kamernetFeed{
		channel:          apartments,
		city:             city,
		seenApartmentIds: make([]string, 0),
	}

	go func() {
		var passedFirst bool
		logCtx.Info("Started looking for apartments...")
		for {
			select {
			case <-ctx.Done():
				close(apartments)
			default:
				newApartments, err := k.FetchNewApartments()
				if err != nil {
					logCtx.WithError(err).Error("Fetching new apartments failed")
					continue
				}

				if !passedFirst {
					logCtx.Info("Received first batch of apartments, looking for new apartments...")
					passedFirst = true
				}

				if len(newApartments) == 0 {
					logCtx.Info("No new apartments")
				}

				for _, apartment := range newApartments {
					apartments <- apartment
				}
				time.Sleep(time.Minute)
			}
		}
	}()

	return k
}

type kamernetFeed struct {
	channel chan feed.Apartment

	city string

	seenApartmentIds []string
}

func (k *kamernetFeed) Apartments() chan feed.Apartment {
	return k.channel
}

func (k *kamernetFeed) FetchNewApartments() (newApartments []feed.Apartment, err error) {
	newApartments = make([]feed.Apartment, 0)

	apartments, err := k.GetApartments()
	if err != nil {
		err = errors.Wrap(err, "couldn't get apartments")
		return
	}

newApartmentsLoop:
	for _, apt := range apartments {
		for _, seenAptId := range k.seenApartmentIds {
			if seenAptId == apt.Id {
				// Already seen this apt!
				continue newApartmentsLoop
			}
		}

		newApartments = append(newApartments, apt)
		k.seenApartmentIds = append(k.seenApartmentIds, apt.Id)
	}
	return newApartments, nil
}

func (k kamernetFeed) processApartment(s *goquery.Selection) *feed.Apartment {
	apt := &feed.Apartment{}
	id, exists := s.Find(".rowSearchResultRoom").Attr("data-roomid")
	if !exists {
		return nil
	}
	apt.Id = id

	apt.Address = s.Find(".tile-block-1 .title").Text()
	apt.Price = s.Find(".tile-block-2 .rent").Text()
	s.Find("meta").Each(func(_ int, sel *goquery.Selection) {
		if prop, exists := sel.Attr("itemprop"); exists && prop == "postalCode" {
			if postcode, exists := sel.Attr("content"); exists {
				apt.Postcode = postcode
			}
		}
	})
	return apt
}

func (k kamernetFeed) GetApartments() (apartments []feed.Apartment, err error) {
	apartments = make([]feed.Apartment, 0)

	doc, err := goquery.NewDocument(fmt.Sprintf("https://kamernet.nl/en/for-rent/room-%s", k.city))
	if err != nil {
		err = errors.Wrap(err, "couldn't fetch kamernet list")
		return
	}

	var foundResults bool
	doc.Find(".search-result-page").Each(func(i int, s *goquery.Selection) {
		// Found results
		foundResults = true
		s.Children().Each(func(i int, s *goquery.Selection) {
			apt := k.processApartment(s)
			if apt == nil {
				return
			}

			apartments = append(apartments, *apt)
		})
	})

	if !foundResults {
		err = errors.New("Kamernet didn't return any results")
		return
	}

	return
}
