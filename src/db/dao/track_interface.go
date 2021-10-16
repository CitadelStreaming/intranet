package dao

import (
	"citadel_intranet/src/db/model"
)

type TrackDao interface {
	BaseDao

	/*
	   Load all tracks
	*/
	LoadAll() []model.Track

	/*
	   Load an track from their id

	   Returns nil if no track is found
	*/
	Load(uint64) *model.Track

	/*
	   Load all tracks associated with an album id.
	*/
	LoadForAlbum(uint64) []model.Track

	/*
	   Save an track via upsert.

	   Returns the last inserted id and an error
	*/
	Save(model.Track) (int64, error)

	/*
	   Delete an track based on its id

	   Returns rows affected and an error
	*/
	Delete(model.Track) (int64, error)
}
