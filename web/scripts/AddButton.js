export class AddButton
{
    constructor()
    {
        this._button = document.createElement("button");
        this._button.classList.add("addButton");
        this._button.addEventListener("click", this._addAlbumModal);
        this._button.innerHTML = "+";

        document.body.appendChild(this._button);
    }

    _addAlbumModal(e)
    {
        console.log(e);
    }
}
