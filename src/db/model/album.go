package model

type Album struct {
	Id        int64   `json:"id"`
	Title     string  `json:"title"`
	Artist    Artist  `json:"artist"`
	Tracks    []Track `json:"tracks"`
	Published bool    `json:"published"`
	Rating    uint    `json:"rating"`
}
