package mysql_test

import (
	"errors"
	"testing"

	"citadel_intranet/src/db/dao/mock"
	"citadel_intranet/src/db/dao/mysql"
	"citadel_intranet/src/db/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAlbumDaoLoadError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
        WHERE id = \?
    `).WithArgs(1).WillReturnError(errors.New("Something bad happened"))

	dao := mysql.NewAlbumDao(db, nil, nil)
	defer dao.Close()

	result := dao.Load(uint64(1))
	assert.Nil(result)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDaoLoad(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	assert := assert.New(t)

	mockArtistDao := mock.NewMockArtistDao(ctrl)
	mockTrackDao := mock.NewMockTrackDao(ctrl)
	db, mock, err := sqlmock.New()
	assert.Nil(err)

	dao := mysql.NewAlbumDao(db, mockArtistDao, mockTrackDao)
	defer dao.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "artist", "published", "rating"}).
		AddRow(1, "Waffle Irons", 42, false, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
        WHERE id = \?
    `).WithArgs(1).WillReturnRows(mockRows)

	mockArtistDao.EXPECT().
		Load(gomock.Eq(uint64(42))).
		DoAndReturn(func(id uint64) *model.Artist {
			return &model.Artist{
				Id:   id,
				Name: "Bobby",
			}
		}).Times(1)

	mockTrackDao.EXPECT().
		LoadForAlbum(gomock.Eq(uint64(1))).
		DoAndReturn(func(_ uint64) []model.Track {
			return []model.Track{}
		}).Times(1)

	album := dao.Load(uint64(1))
	assert.NotNil(album)

	assert.Equal(uint64(1), album.Id)
	assert.Equal(uint64(42), album.Artist.Id)
	assert.Equal("Bobby", album.Artist.Name)
	assert.Equal("Waffle Irons", album.Title)
	assert.Equal(false, album.Published)
	assert.Equal(uint(5), album.Rating)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDaoLoadNoArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	assert := assert.New(t)

	mockArtistDao := mock.NewMockArtistDao(ctrl)
	mockTrackDao := mock.NewMockTrackDao(ctrl)
	db, mock, err := sqlmock.New()
	assert.Nil(err)

	dao := mysql.NewAlbumDao(db, mockArtistDao, mockTrackDao)
	defer dao.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "artist", "published", "rating"}).
		AddRow(1, "Waffle Irons", 42, false, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
        WHERE id = \?
    `).WithArgs(1).WillReturnRows(mockRows)

	mockArtistDao.EXPECT().
		Load(gomock.Eq(uint64(42))).
		DoAndReturn(func(id uint64) *model.Artist {
			return nil
		}).Times(1)

	mockTrackDao.EXPECT().
		LoadForAlbum(gomock.Eq(uint64(1))).
		DoAndReturn(func(_ uint64) []model.Track {
			return []model.Track{}
		}).Times(1)

	album := dao.Load(uint64(1))
	assert.NotNil(album)

	assert.Equal(uint64(1), album.Id)
	assert.Equal(uint64(0), album.Artist.Id)
	assert.Equal("", album.Artist.Name)
	assert.Equal("Waffle Irons", album.Title)
	assert.Equal(false, album.Published)
	assert.Equal(uint(5), album.Rating)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDaoDeleteError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	album := model.Album{
		Id: 42,
	}

	mock.ExpectExec(`
        DELETE
        FROM album
        WHERE id = \?
    `).
		WithArgs(album.Id).
		WillReturnError(errors.New("That's not a real album"))

	dao := mysql.NewAlbumDao(db, nil, nil)
	defer dao.Close()

	rowsAffected, err := dao.Delete(album)
	assert.Equal(int64(0), rowsAffected)
	assert.NotNil(err)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDaoSaveError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	album := model.Album{
		Id: 42,
	}

	mock.ExpectExec(`
        INSERT INTO album\(
            id,
            title,
            artist,
            published,
            rating
        \)
        VALUES\(
            \?,
            \?,
            \?,
            \?,
            \?
        \)
        ON DUPLICATE KEY UPDATE
            title = VALUES\(title\),
            artist = VALUES\(artist\),
            published = VALUES\(published\),
            rating = VALUES\(rating\)
    `).
		WithArgs(album.Id, album.Title, album.Artist.Id, album.Published, album.Rating).
		WillReturnError(errors.New("Album save died"))

	dao := mysql.NewAlbumDao(db, nil, nil)
	defer dao.Close()

	lastId, err := dao.Save(album)
	assert.Equal(int64(0), lastId)
	assert.NotNil(err)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDao(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	album := model.Album{
		Id: 42,
	}

	mock.ExpectExec(`
        INSERT INTO album\(
            id,
            title,
            artist,
            published,
            rating
        \)
        VALUES\(
            \?,
            \?,
            \?,
            \?,
            \?
        \)
        ON DUPLICATE KEY UPDATE
            title = VALUES\(title\),
            artist = VALUES\(artist\),
            published = VALUES\(published\),
            rating = VALUES\(rating\)
    `).
		WithArgs(album.Id, album.Title, album.Artist.Id, album.Published, album.Rating).
		WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectExec(`
        DELETE
        FROM album
        WHERE id = \?
    `).
		WithArgs(album.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	dao := mysql.NewAlbumDao(db, nil, nil)
	defer dao.Close()

	lastId, err := dao.Save(album)
	assert.Equal(int64(42), lastId)
	assert.Nil(err)

	affectedRows, err := dao.Delete(album)
	assert.Equal(int64(1), affectedRows)
	assert.Nil(err)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDaoLoadAllError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
    `).
		WillReturnError(errors.New("Something bad happened"))

	dao := mysql.NewAlbumDao(db, nil, nil)
	defer dao.Close()

	result := dao.LoadAll()
	assert.Nil(result)

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDaoLoadAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	assert := assert.New(t)

	mockArtistDao := mock.NewMockArtistDao(ctrl)
	mockTrackDao := mock.NewMockTrackDao(ctrl)
	db, mock, err := sqlmock.New()
	assert.Nil(err)

	dao := mysql.NewAlbumDao(db, mockArtistDao, mockTrackDao)
	defer dao.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "artist", "published", "rating"}).
		AddRow(1, "Waffle Irons", 42, false, 5).
		AddRow(2, "Something New", 42, false, 3).
		AddRow(3, "Something New (Deluxe)", 42, true, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
    `).
		WillReturnRows(mockRows)

	mockArtistDao.EXPECT().
		Load(gomock.Eq(uint64(42))).
		DoAndReturn(func(id uint64) *model.Artist {
			return &model.Artist{
				Id:   id,
				Name: "Bobby",
			}
		}).
		Times(3)

	for i := 1; i < 4; i++ {
		mockTrackDao.EXPECT().
			LoadForAlbum(gomock.Eq(uint64(i))).
			DoAndReturn(func(_ uint64) []model.Track {
				return []model.Track{}
			}).Times(1)
	}

	albums := dao.LoadAll()
	assert.Len(albums, 3)

	for index, row := range []struct {
		Title     string
		Published bool
		Rating    uint
	}{
		{"Waffle Irons", false, 5},
		{"Something New", false, 3},
		{"Something New (Deluxe)", true, 5},
	} {
		album := albums[index]
		assert.Equal(uint64(index+1), album.Id)
		assert.Equal(uint64(42), album.Artist.Id)
		assert.Equal("Bobby", album.Artist.Name)
		assert.Equal(row.Title, album.Title)
		assert.Equal(row.Published, album.Published)
		assert.Equal(row.Rating, album.Rating)
	}

	assert.Nil(mock.ExpectationsWereMet())
}

func TestAlbumDaoLoadAllScanError(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	dao := mysql.NewAlbumDao(db, nil, nil)
	defer dao.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title", "artist", "published", "rating"}).
		AddRow("cat", "Waffle Irons", 42, false, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
    `).
		WillReturnRows(mockRows)

	albums := dao.LoadAll()
	assert.NotNil(albums)
	assert.Len(albums, 0)

	assert.Nil(mock.ExpectationsWereMet())
}

func disableTestAlbumDaoLoadAll(t *testing.T) {
	assert := assert.New(t)

	db, mock, err := sqlmock.New()
	assert.Nil(err)

	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(uint64(1), "Album 1").
		AddRow(uint64(2), "Little Bobby Tables").
		AddRow(uint64(3), "Album 42").
		RowError(1, errors.New("Keep him away from the tables!"))
	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
    `).
		WillReturnRows(mockRows)

	dao := mysql.NewAlbumDao(db, nil, nil)
	defer dao.Close()

	result := dao.LoadAll()
	assert.NotNil(result)
	assert.Len(result, 1)

	assert.Nil(mock.ExpectationsWereMet())
}
