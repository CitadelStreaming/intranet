package dao

import (
    "citadel_intranet/src/db/models"
)

type AlbumDao interface {
    BaseDao

    LoadAll() []models.Album
    Load(uint64) *models.Album
    Save(models.Album) (int64, error)
    Delete(models.Album) (int64, error)
}
