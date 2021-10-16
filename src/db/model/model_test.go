package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"citadel_intranet/src/db/model"
)

func TestModelInstantiation(t *testing.T) {
	assert := assert.New(t)
	album := &model.Album{}
	assert.NotNil(album)

	track := &model.Track{}
	assert.NotNil(track)

	artist := &model.Artist{}
	assert.NotNil(artist)
}
