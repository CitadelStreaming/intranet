package dao

import (
    "citadel_intranet/src/db/models"
)

type ArtistDao interface {
    BaseDao

    LoadAll() []models.Artist
    Load(uint64) *models.Artist
    Save(models.Artist) (int64, error)
    Delete(models.Artist) (int64, error)
}
