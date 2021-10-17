package model

type Track struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	AlbumId int64  `json:"album"`
	Rating  uint   `json:"rating"`
}
