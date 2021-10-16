package mysql

import (
	"database/sql"

	"citadel_intranet/src/db/dao"
	"citadel_intranet/src/db/model"

	"github.com/sirupsen/logrus"
)

type albumDao struct {
	db        *sql.DB
	artistDao dao.ArtistDao
	trackDao  dao.TrackDao
}

func NewAlbumDao(db *sql.DB, artistDao dao.ArtistDao, trackDao dao.TrackDao) dao.AlbumDao {
	return albumDao{
		db:        db,
		artistDao: artistDao,
		trackDao:  trackDao,
	}
}

func (this albumDao) Close() {
	logrus.Debug("Closing Album DAO")
}

func (this albumDao) loadArtistAndTracksForAlbum(album *model.Album, artistId uint64) {
	artist := this.artistDao.Load(artistId)
	if artist == nil {
		logrus.Error("Unable to find artist with ID=", artistId)
	} else {
		album.Artist = *artist
	}

	album.Tracks = this.trackDao.LoadForAlbum(album.Id)
}

func (this albumDao) Load(id uint64) *model.Album {
	var album *model.Album = &model.Album{}

	row := this.db.QueryRow(`
        SELECT
            *
        FROM album
        WHERE id = ?
    `, id)

	var artistId uint64
	err := row.Scan(&album.Id, &album.Title, &artistId, &album.Published, &album.Rating)

	if err != nil {
		logrus.Warn("Loading failed for ", id, " ", err.Error())
		return nil
	}

	this.loadArtistAndTracksForAlbum(album, artistId)

	return album
}

func (this albumDao) LoadAll() []model.Album {
	var ret []model.Album = make([]model.Album, 0)

	rows, err := this.db.Query(`
        SELECT
            *
        FROM album
    `)

	if err != nil {
		logrus.Warn("Unable to load albums ", err.Error())
		return nil
	}

	var artistId uint64
	for rows.Next() {
		var album model.Album
		err := rows.Scan(&album.Id, &album.Title, &artistId, &album.Published, &album.Rating)

		if err != nil {
			logrus.Warn(err.Error())
		} else {
			this.loadArtistAndTracksForAlbum(&album, artistId)
			ret = append(ret, album)
		}
	}

	return ret
}

func (this albumDao) Delete(album model.Album) (int64, error) {
	result, err := this.db.Exec(`
        DELETE
        FROM album
        WHERE id = ?
    `, album.Id)

	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this albumDao) Save(album model.Album) (int64, error) {
	result, err := this.db.Exec(`
        INSERT INTO album(
            id,
            title,
            artist,
            published,
            rating
        )
        VALUES(
            ?,
            ?,
            ?,
            ?,
            ?
        )
        ON DUPLICATE KEY UPDATE
            title = VALUES(title),
            artist = VALUES(artist),
            published = VALUES(published),
            rating = VALUES(rating)
    `,
		album.Id,
		album.Title,
		album.Artist.Id,
		album.Published,
		album.Rating,
	)

	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
