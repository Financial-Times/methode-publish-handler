package vanities

// Vanity Contains the vanity web url.
type Vanity struct {
	WebURL string `json:"webUrl"`
}

// GetVanity Given (data tbc) it will lookup the vanity url and return it as a Vanity struct.
// TODO Move this back to handlers?
func GetVanity() Vanity {
	return Vanity{""}
}
