import {AddButton} from "./AddButton.js";
import {AlbumController} from "./AlbumController.js";

export class Application
{
    constructor(bodySelector)
    {
        this._albums = new AlbumController(bodySelector);
        this._addButton = new AddButton();
    }
}
