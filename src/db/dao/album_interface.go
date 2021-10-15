package dao

import (
    "citadel_intranet/src/db/model"
)

type AlbumDao interface {
    BaseDao

    LoadAll() []model.Album
    Load(uint64) *model.Album
    Save(model.Album) (int64, error)
    Delete(model.Album) (int64, error)
}
