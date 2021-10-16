ALTER TABLE track
DROP CONSTRAINT track_ibfk_1;

ALTER TABLE track
ADD CONSTRAINT track_to_album_mapping
    FOREIGN KEY (album)
    REFERENCES album(id)
    ON DELETE CASCADE;


ALTER TABLE album
DROP CONSTRAINT album_ibfk_1;

ALTER TABLE album
ADD CONSTRAINT album_to_artist_mapping
    FOREIGN KEY (artist)
    REFERENCES artist(id)
    ON DELETE CASCADE;
