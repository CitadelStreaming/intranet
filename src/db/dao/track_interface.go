package dao

import (
    "citadel_intranet/src/db/model"
)

type TrackDao interface {
    BaseDao

    LoadAll() []model.Track
    Load(uint64) *model.Track
    Save(model.Track) (int64, error)
    Delete(model.Track) (int64, error)
}
