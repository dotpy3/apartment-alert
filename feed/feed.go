package feed

type Feed interface {
	Apartments() chan Apartment
}
