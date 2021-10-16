package dao

import (
	"citadel_intranet/src/db/model"
)

/*
Used for CRUD operations on an Album
*/
type AlbumDao interface {
	BaseDao

	/*
	   Load all albums from the database.
	*/
	LoadAll() []model.Album

	/*
	   Load a single album based on id, will return nil if the album cannot be
	   found.
	*/
	Load(int64) *model.Album

	/*
	   Save an album. This should perform an upsert style insert or update to an
	   album in the case where it already exists.

	   Returns the last inserted id and an error
	*/
	Save(model.Album) (int64, error)

	/*
	   Delete an album, based on the id of the album.

	   Returns the affected rows and an error
	*/
	Delete(model.Album) (int64, error)
}
