package mysql

import (
	"database/sql"

	"citadel_intranet/src/db/dao"
	"citadel_intranet/src/db/model"

	"github.com/sirupsen/logrus"
)

type trackDao struct {
    db *sql.DB
}

func NewTrackDao(db *sql.DB) dao.TrackDao {
    return trackDao{
        db: db,
    }
}

func (this trackDao) Close() {
    logrus.Debug("Closing Track DAO")
}

func (this trackDao) scanAll(rows *sql.Rows) []model.Track {
    var ret []model.Track = make([]model.Track, 0)

    for rows.Next() {
        var track model.Track
        err := rows.Scan(&track.Id, &track.Title, &track.AlbumId, &track.Rating)

        if err != nil {
            logrus.Warn(err.Error())
        } else {
            ret = append(ret, track)
        }
    }

    return ret
}

func (this trackDao) LoadForAlbum(id uint64) []model.Track {
    rows, err := this.db.Query(`
        SELECT
            *
        FROM track
        WHERE album = ?
    `, id)

    if err != nil {
        logrus.Warn("Unable to load tracks ", err.Error())
        return nil
    }

    return this.scanAll(rows)
}

func (this trackDao) Load(id uint64) *model.Track {
    var track *model.Track = &model.Track{}

    row := this.db.QueryRow(`
        SELECT
            *
        FROM track
        WHERE id = ?
    `, id)

    err := row.Scan(&track.Id, &track.Title, &track.AlbumId, &track.Rating)

    if err != nil {
        logrus.Warn("Loading failed for ", id, " ", err.Error())
        return nil
    }

    return track
}

func (this trackDao) LoadAll() []model.Track {
    rows, err := this.db.Query(`
        SELECT
            *
        FROM track
    `)

    if err != nil {
        logrus.Warn("Unable to load tracks ", err.Error())
        return nil
    }

    return this.scanAll(rows)
}

func (this trackDao) Delete(track model.Track) (int64, error) {
    result, err := this.db.Exec(`
        DELETE
        FROM track
        WHERE id = ?
    `, track.Id)

    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}

func (this trackDao) Save(track model.Track) (int64, error) {
    result, err := this.db.Exec(`
        INSERT INTO track(
            id,
            title,
            album,
            rating
        )
        VALUES(
            ?,
            ?,
            ?,
            ?
        )
        ON DUPLICATE KEY UPDATE
            title = VALUES(title),
            album = VALUES(album),
            rating = VALUES(rating)
    `,
        track.Id,
        track.Title,
        track.AlbumId,
        track.Rating,
    )

    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}
