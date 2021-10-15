package mysql

import (
	"database/sql"

	"citadel_intranet/src/db/dao"
	"citadel_intranet/src/db/model"

	"github.com/sirupsen/logrus"
)

type artistDao struct {
    db *sql.DB
}

func NewArtistDao(db *sql.DB) dao.ArtistDao {
    return artistDao{
        db: db,
    }
}

func (this artistDao) Close() {
}

func (this artistDao) Load(id uint64) *model.Artist {
    var artist *model.Artist = &model.Artist{}

    row := this.db.QueryRow(`
        SELECT
            *
        FROM artist
        WHERE id = ?
    `, id)

    err := row.Scan(&artist.Id, &artist.Name)

    if err != nil {
        logrus.Warn("Loading failed for ", id, " ", err.Error())
        return nil
    }

    return artist
}

func (this artistDao) LoadAll() []model.Artist {
    var ret []model.Artist

    rows, err := this.db.Query(`
        SELECT
            *
        FROM artist
    `)

    if err != nil {
        logrus.Warn("Unable to load artists ", err.Error())
        return nil
    }

    for rows.Next() {
        var artist model.Artist
        err := rows.Scan(&artist.Id, &artist.Name)

        if err != nil {
            logrus.Warn(err.Error())
        } else {
            ret = append(ret, artist)
        }
    }
    
    return ret
}

func (this artistDao) Delete(artist model.Artist) (int64, error) {
    result, err := this.db.Exec(`
        DELETE
        FROM artist
        WHERE id = ?
    `, artist.Id)

    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}

func (this artistDao) Save(artist model.Artist) (int64, error) {
    result, err := this.db.Exec(`
        INSERT INTO artist(
            id,
            name
        )
        VALUES(
            ?,
            ?
        )
        ON DUPLICATE KEY UPDATE
            name = VALUES(name)
    `,
        artist.Id,
        artist.Name,
    )

    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}
