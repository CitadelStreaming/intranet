import {AlbumEditController} from "./AlbumEditController.js";

export class AlbumListController
{
    constructor(selector)
    {
        this._selector = selector;
        this._albums = [];
        this._request = new Request("/api/v1/album");

        const that = this;
        document.body.addEventListener("reloadAlbums", function()
        {
            that._loadAllAlbums();
        });
        this._loadAllAlbums();
    }

    _loadAllAlbums()
    {
        const that = this;

        fetch(this._request)
            .then(function(response)
            {
                if (!response.ok)
                {
                    throw new Error("Unable to retrieve album list");
                }

                return response.json();
            })
            .then(function(data)
            {
                that._albums = data;
                that._drawAlbums();
            })
            .catch(console.error);
    }

    _loadAllArtists(target)
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
                that._artistSelect(target, data);
            })
            .catch(console.error);
    }

    _saveAlbum(album)
    {
        const json = JSON.stringify(album);
        fetch(new Request("/api/v1/album/" + album.id), {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: json
        })
            .then(function(response)
            {
                if (response.ok)
                {
                    return;
                }

                return response.json();
            })
            .then(function(data)
            {
                if (data === undefined)
                {
                    return;
                }

                console.error("Encountered a problem: " + data);
            })
            .catch(console.error);
    }

    _artistSelect(target, data)
    {
        const that = this;

        let select = document.createElement("select");
        for (const artist of data)
        {
            let opt = document.createElement("option");
            opt.setAttribute("value", artist.id);
            opt.textContent = artist.name;

            if (opt.textContent === target.textContent) {
                opt.setAttribute("selected", "true");
            }

            select.appendChild(opt);
        }

        select.addEventListener("change", function(e)
        {
            let opt = this.options[this.selectedIndex];
            let album = JSON.parse(this.parentNode.parentNode.getAttribute("data-album"));
            album.artist.id = parseInt(opt.value);
            album.artist.name = opt.textContent;

            target.textContent = opt.textContent;
            this.parentNode.appendChild(target);
            this.parentNode.removeChild(this);

            that._saveAlbum(album);
        });
        select.addEventListener("blur", function(e)
        {
            this.parentNode.appendChild(target);
            this.parentNode.removeChild(this);
        });

        target.parentNode.appendChild(select);
        target.parentNode.removeChild(target);
        select.focus();
    }

    _drawAlbums()
    {
        const that = this;
        let container = document.querySelector(this._selector);
        container.innerHTML = "";

        for (const album of this._albums)
        {
            let item = document.createElement("div");
            item.classList.add("album");
            item.setAttribute("data-album", JSON.stringify(album));
            
            let title = document.createElement("h2");
            item.appendChild(title).textContent = album.title;
            title.addEventListener("click", function(e)
            {
                const editor = new AlbumEditController(album);
            });

            let byline = document.createElement("span");
            item.appendChild(byline).appendChild(document.createTextNode(" by "));

            let artist = document.createElement("em");
            byline.appendChild(artist).textContent = album.artist.name;
            artist.addEventListener("click", function(e)
            {
                that._loadAllArtists(e.target);
            });

            container.appendChild(item);
        }
    }
}
