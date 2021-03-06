export class AlbumEditController
{
    constructor(data)
    {
        this._album = data;

        this._loadAllArtists();
    }

    _closeAllModals()
    {
        const e = new CustomEvent("reloadAlbums", {});
        document.body.dispatchEvent(e);
        document.querySelectorAll(".modal").forEach(function(item)
        {
            item.parentNode.removeChild(item);
        });
    }

    _loadAllArtists()
    {
        const that = this;

        fetch(new Request("/api/v1/artist"))
            .then(function(response)
            {
                if (!response.ok)
                {
                    throw new Error("Unable to retrieve artist list");
                }

                return response.json();
            })
            .then(function(data)
            {
                console.log(data);
                that._displayModal(data);
            })
            .catch(console.error);
    }

    _deleteAlbum(album)
    {
        const that = this;
        const json = JSON.stringify(album);

        fetch(new Request("/api/v1/album/" + album.id), {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json"
            },
            body: json
        })
            .then(function(response)
            {
                if (response.ok)
                {
                    return new Promise((kept, broken) => {
                        kept(album);
                    });
                }
                return response.json();
            })
            .then(function(data)
            {
                that._closeAllModals();
                if (data.error === undefined)
                {
                    return;
                }

                console.error("Encountered a problem: " + data);
            })
            .catch(console.error);
    }

    _saveAlbum(album)
    {
        if (!album.title || album.title == "" || !album.artist || !album.artist.name || album.artist.name == "")
        {
            console.log("Malformed album");
            return;
        }

        const that = this;
        const json = JSON.stringify(album);
        let uri = "/api/v1/album";
        let method = "POST";
        if (album.id && album.id > 0)
        {
            uri = "/api/v1/album/" + album.id;
            method = "PUT";
        }

        fetch(new Request(uri), {
            method: method,
            headers: {
                "Content-Type": "application/json"
            },
            body: json
        })
            .then(function(response)
            {
                if (method == "PUT" && response.ok)
                {
                    return new Promise((kept, broken) => {
                        kept(album);
                    });
                }
                return response.json();
            })
            .then(function(data)
            {
                if (data.error === undefined)
                {
                    that._album = data;
                    that._closeAllModals();
                    that._loadAllArtists();
                    return;
                }

                console.error("Encountered a problem: " + data);
            })
            .catch(console.error);
    }

    _displayModal(artistList)
    {
        const that = this;
        let container = document.createElement("div");
        container.classList.add("modal");
        container.addEventListener("click", function(e)
        {
            if (e.target !== this)
            {
                return;
            }

            that._closeAllModals();
        });

        let body = document.createElement("div");
        container.appendChild(body);

        let title = document.createElement("input");
        title.classList.add("title");
        title.setAttribute("type", "text");
        title.setAttribute("value", this._album.title || "");
        title.setAttribute("placeholder", "Album Title");
        title.setAttribute("pattern", ".{5,}");
        body.appendChild(title);

        this._displayArtistList(body, artistList);
        this._displayTracks(body);
        this._displayDeleteButton(body);
        this._displaySaveButton(body);

        document.body.appendChild(container);
    }

    _displayDeleteButton(body)
    {
        if (!this._album.id || this._album.id == 0)
        {
            return;
        }

        const that = this;
        let del = document.createElement("button");
        del.classList.add("delete");
        del.textContent = "Delete";
        del.addEventListener("click", function(e)
        {
            that._deleteAlbum(that._album);
        });

        body.appendChild(del);
    }

    _displaySaveButton(body)
    {
        const that = this;
        let save = document.createElement("button");
        save.classList.add("save");
        save.textContent = "Save";
        save.addEventListener("click", function(e)
        {
            let album = that._album;

            album.title = body.querySelector(".title").value;
            const artistInput = body.querySelector(".artist");
            if (artistInput.value != "" && artistInput.validity.valid)
            {
                album.artist.id = 0;
                album.artist.name = artistInput.value;
            }
            else
            {
                const sel = body.querySelector("select");
                const opt = sel.options[sel.selectedIndex];
                album.artist.id = parseInt(opt.value);
                album.artist.name = opt.textContent;
            }

            that._saveAlbum(album);
        });

        body.appendChild(save);
    }

    _displayArtistList(body, artistList)
    {
        let select = document.createElement("select");
        for (const artist of artistList)
        {
            console.log(artist);
            let opt = document.createElement("option");
            opt.setAttribute("value", artist.id);
            opt.textContent = artist.name;
            if (artist.id == this._album.artist.id) {
                opt.setAttribute("selected", "true");
            }
            select.appendChild(opt);
        }

        body.appendChild(document.createElement("br"));
        body.appendChild(document.createTextNode("Artist: "));
        body.appendChild(select);
        body.appendChild(document.createTextNode(" -or- "));

        let artistInput = document.createElement("input");
        artistInput.classList.add("artist");
        artistInput.setAttribute("type", "text");
        artistInput.setAttribute("placeholder", "Artist Name");
        artistInput.setAttribute("pattern", ".{2,}");

        body.appendChild(artistInput);
        body.appendChild(document.createElement("br"));
    }

    _displayTracks(body)
    {
        body.appendChild(document.createTextNode("Tracks:"));
        body.appendChild(document.createElement("br"));

        for (const track of this._album.tracks)
        {
            let t = document.createElement("div");
            t.classList.add("track");
            t.textContent = track.title;
            t.setAttribute("data-track", JSON.stringify(track));

            t.addEventListener("dblclick", function(e)
            {
                const target = e.target;

                console.log("Double Clicked!");
            });

            body.appendChild(t);
        }

        this._addNewTrackButton(body);
    }
    
    _addNewTrackButton(body)
    {
        const that = this;
        if (this._album.id && this._album.id > 0)
        {
            let t = document.createElement("div");
            t.classList.add("addTrack");
            t.textContent = "+";

            t.addEventListener("click", function(e)
            {
                const id = that._album.id;
                console.log("Starting new track addition for album ", id);

                let newTrack = document.createElement("div");
                newTrack.classList.add("track");

                let input = document.createElement("input");
                input.setAttribute("type", "text");
                input.setAttribute("placeholder", "Track Title");

                input.addEventListener("blur", function(e)
                {
                    const value = this.value;
                    if (value == "")
                    {
                        newTrack.parentNode.removeChild(newTrack);
                        return;
                    }

                    that._saveTrack(id, value, newTrack);
                });

                newTrack.appendChild(input);
                t.parentNode.insertBefore(newTrack, t);
                input.focus();
            });

            body.appendChild(t);
        }
    }

    _saveTrack(albumId, trackTitle, trackEle)
    {
        const json = JSON.stringify({
            title: trackTitle,
            album: albumId
        });
        console.log("Saving track", json);

        fetch(new Request("/api/v1/track"), {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: json
        })
            .then(function(response)
            {
                return response.json();
            })
            .then(function(data)
            {
                if (data.error === undefined)
                {
                    trackEle.textContent = trackTitle;
                    return;
                }

                console.error("Encountered a problem: " + data);
            })
            .catch(console.error);
    }
}
