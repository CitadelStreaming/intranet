package dao

type BaseDao interface {
    /*
    Close the DAO, this is meant to be the last thing called on a DAO as the
    application is gracefully shutting down and should take any actions
    necessary to close down prepared statements or other resources the DAO
    stood up on instantiation.
    */
    Close()
}
