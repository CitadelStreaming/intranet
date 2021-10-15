package dao

import (
    "citadel_intranet/src/db/model"
)

type ArtistDao interface {
    BaseDao

    LoadAll() []model.Artist
    Load(uint64) *model.Artist
    Save(model.Artist) (int64, error)
    Delete(model.Artist) (int64, error)
}
