import {AddButton} from "./AddButton.js";
import {AlbumListController} from "./AlbumListController.js";

export class Application
{
    constructor(bodySelector)
    {
        this._albums = new AlbumListController(bodySelector);
        this._addButton = new AddButton();
    }
}
