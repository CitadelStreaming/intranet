package mysql_test

import (
	"errors"
	"testing"

	"citadel_intranet/src/db/dao/mysql"
	"citadel_intranet/src/db/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestArtistDao(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	artist := model.Artist{
		Name: "James",
	}

	mock.ExpectExec(`
        INSERT INTO artist\(
            id,
            name
        \)
        VALUES\(
            \?,
            \?
        \)
        ON DUPLICATE KEY UPDATE
            name = VALUES\(name\)
    `).
		WithArgs(artist.Id, artist.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(int64(1), "James")
	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
        WHERE id = \?
    `).
		WithArgs(1).
		WillReturnRows(mockRows)

	mock.ExpectExec(`
        DELETE
        FROM artist
        WHERE id = \?
    `).
		WithArgs(artist.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	dao := mysql.NewArtistDao(db)
	defer dao.Close()

	rows, err := dao.Save(artist)
	assert.Nil(err)
	assert.Equal(int64(1), rows)

	result := dao.Load(int64(1))
	assert.NotNil(result)
	assert.Equal(artist.Name, result.Name)

	rows, err = dao.Delete(artist)
	assert.Nil(err)
	assert.Equal(int64(1), rows)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestArtistDaoLoadError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
        WHERE id = \?
    `).
		WithArgs(1).
		WillReturnError(errors.New("Something bad happened"))

	dao := mysql.NewArtistDao(db)
	defer dao.Close()

	result := dao.Load(int64(1))
	assert.Nil(result)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestArtistDaoDeleteError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	artist := model.Artist{
		Id:   42,
		Name: "Bobby",
	}

	mock.ExpectExec(`
        DELETE
        FROM artist
        WHERE id = \?
    `).
		WithArgs(artist.Id).
		WillReturnError(errors.New("That's not a real user"))

	dao := mysql.NewArtistDao(db)
	defer dao.Close()

	rowsAffected, err := dao.Delete(artist)
	assert.Equal(int64(0), rowsAffected)
	assert.NotNil(err)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestArtistDaoSaveError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	artist := model.Artist{
		Id:   42,
		Name: "Bobby",
	}

	mock.ExpectExec(`
        INSERT INTO artist\(
            id,
            name
        \)
        VALUES\(
            \?,
            \?
        \)
        ON DUPLICATE KEY UPDATE
            name = VALUES\(name\)
    `).
		WithArgs(artist.Id, artist.Name).
		WillReturnError(errors.New("That's not a real user"))

	dao := mysql.NewArtistDao(db)
	defer dao.Close()

	lastId, err := dao.Save(artist)
	assert.Equal(int64(0), lastId)
	assert.NotNil(err)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestArtistDaoLoadAllError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
    `).
		WillReturnError(errors.New("Something bad happened"))

	dao := mysql.NewArtistDao(db)
	defer dao.Close()

	result := dao.LoadAll()
	assert.Nil(result)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestArtistDaoLoadAllScanError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("cat", "James")
	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
    `).
		WillReturnRows(mockRows)

	dao := mysql.NewArtistDao(db)
	defer dao.Close()

	result := dao.LoadAll()
	assert.NotNil(result)
	assert.Len(result, 0)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestArtistDaoLoadAll(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(int64(1), "James").
		AddRow(int64(2), "Bobby").
		AddRow(int64(3), "Frank").
		RowError(1, errors.New("Keep him away from the tables!"))
	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
    `).
		WillReturnRows(mockRows)

	dao := mysql.NewArtistDao(db)
	defer dao.Close()

	result := dao.LoadAll()
	assert.NotNil(result)
	assert.Len(result, 1)

	assert.Nil(mock.ExpectationsWereMet())
}
