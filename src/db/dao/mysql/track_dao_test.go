package mysql_test

import (
	"errors"
	"testing"

	"citadel_intranet/src/db/dao/mysql"
	"citadel_intranet/src/db/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTrackDao(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	track := model.Track{
		Id:      456,
		Title:   "Something Awesome",
		AlbumId: 123,
		Rating:  5,
	}

	mock.ExpectExec(`
        INSERT INTO track\(
            id,
            title,
            album,
            rating
        \)
        VALUES\(
            \?,
            \?,
            \?,
            \?
        \)
        ON DUPLICATE KEY UPDATE
            title = VALUES\(title\),
            album = VALUES\(album\),
            rating = VALUES\(rating\)
    `).
		WithArgs(track.Id, track.Title, track.AlbumId, track.Rating).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockRows := sqlmock.NewRows([]string{"id", "title", "album", "rating"}).
		AddRow(uint64(456), "Something Awesome", uint64(123), uint(5))
	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
        WHERE id = \?
    `).
		WithArgs(456).
		WillReturnRows(mockRows)

	mock.ExpectExec(`
        DELETE
        FROM track
        WHERE id = \?
    `).
		WithArgs(track.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	rows, err := dao.Save(track)
	assert.Nil(err)
	assert.Equal(int64(1), rows)

	result := dao.Load(uint64(456))
	assert.NotNil(result)
	assert.Equal(track.Title, result.Title)
	assert.Equal(track.AlbumId, result.AlbumId)
	assert.Equal(track.Rating, result.Rating)

	rows, err = dao.Delete(track)
	assert.Nil(err)
	assert.Equal(int64(1), rows)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoLoadForAlbum(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "album", "rating"}).
		AddRow(uint64(457), "Track 1", uint64(123), uint(5)).
		AddRow(uint64(456), "Track 2", uint64(123), uint(5)).
		AddRow(uint64(458), "Track 3", uint64(123), uint(5))
	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
        WHERE album = \?
    `).
		WithArgs(1).
		WillReturnRows(mockRows)

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	result := dao.LoadForAlbum(uint64(1))
	assert.Len(result, 3)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoLoadForAlbumError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
        WHERE album = \?
    `).
		WithArgs(1).
		WillReturnError(errors.New("Something bad happened"))

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	result := dao.LoadForAlbum(uint64(1))
	assert.Nil(result)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoLoadError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
        WHERE id = \?
    `).
		WithArgs(1).
		WillReturnError(errors.New("Something bad happened"))

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	result := dao.Load(uint64(1))
	assert.Nil(result)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoDeleteError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	track := model.Track{
		Id:    42,
		Title: "Bobby",
	}

	mock.ExpectExec(`
        DELETE
        FROM track
        WHERE id = \?
    `).
		WithArgs(track.Id).
		WillReturnError(errors.New("No song found"))

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	rowsAffected, err := dao.Delete(track)
	assert.Equal(int64(0), rowsAffected)
	assert.NotNil(err)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoSaveError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	track := model.Track{
		Id:    42,
		Title: "Bobby",
	}

	mock.ExpectExec(`
        INSERT INTO track\(
            id,
            title,
            album,
            rating
        \)
        VALUES\(
            \?,
            \?,
            \?,
            \?
        \)
        ON DUPLICATE KEY UPDATE
            title = VALUES\(title\),
            album = VALUES\(album\),
            rating = VALUES\(rating\)
    `).
		WithArgs(track.Id, track.Title, track.AlbumId, track.Rating).
		WillReturnError(errors.New("That's not a real user"))

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	lastId, err := dao.Save(track)
	assert.Equal(int64(0), lastId)
	assert.NotNil(err)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoLoadAllError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
    `).WillReturnError(errors.New("Something bad happened"))

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	result := dao.LoadAll()
	assert.Nil(result)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoLoadAllScanError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "album", "rating"}).
		AddRow("cat", "Track 1", uint64(123), uint(5))
	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
    `).WillReturnRows(mockRows)

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	result := dao.LoadAll()
	assert.NotNil(result)
	assert.Len(result, 0)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestTrackDaoLoadAll(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "album", "rating"}).
		AddRow(uint64(457), "Track 1", uint64(123), uint(5)).
		AddRow(uint64(456), "Track 2", uint64(123), uint(5)).
		AddRow(uint64(458), "Track 3", uint64(123), uint(5))
	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
    `).WillReturnRows(mockRows)

	dao := mysql.NewTrackDao(db)
	defer dao.Close()

	result := dao.LoadAll()
	assert.NotNil(result)
	assert.Len(result, 3)

	assert.Nil(mock.ExpectationsWereMet())
}
