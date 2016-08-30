package vanities

// Vanity Contains the vanity web url.
type Vanity struct {
	WebURL string `json:"webUrl"`
}

// VanityService Requests the URL
type VanityService interface {
	GetVanity() Vanity
}

// GetVanity Given (data tbc) it will lookup the vanity url and return it as a Vanity struct.
// TODO Move this back to handlers?
func (v Vanity) GetVanity() Vanity {
	return Vanity{""}
}
