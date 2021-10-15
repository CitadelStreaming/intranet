package dao

import (
    "citadel_intranet/src/db/model"
)

/*
CRUD operations for an artist
*/
type ArtistDao interface {
    BaseDao

    /*
    Load all artists
    */
    LoadAll() []model.Artist

    /*
    Load an artist from their id

    Returns nil if no artist is found
    */
    Load(uint64) *model.Artist

    /*
    Save an artist via upsert.

    Returns the last inserted id and an error
    */
    Save(model.Artist) (int64, error)

    /*
    Delete an artist based on its id

    Returns rows affected and an error
    */
    Delete(model.Artist) (int64, error)
}
