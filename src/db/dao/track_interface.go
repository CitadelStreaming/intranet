package dao

import (
    "citadel_intranet/src/db/models"
)

type TrackDao interface {
    BaseDao

    LoadAll() []models.Track
    Load(uint64) *models.Track
    Save(models.Track) (int64, error)
    Delete(models.Track) (int64, error)
}
