package alerts

import (
	"github.com/dotpy3/apartment-alert/feed"
)

type ApartmentAlerter interface {
	Push(apt feed.Apartment) error
}
