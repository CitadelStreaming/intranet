package model

type Album struct {
	Id        uint64
	Title     string
	Artist    Artist
	Tracks    []Track
	Published bool
	Rating    uint
}
