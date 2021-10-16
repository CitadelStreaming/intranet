package application

type Application interface {
	/*
	   Run the application. This function should only return when the application
	   has been completed, or Closed.
	*/
	Run()

	/*
	   Close the application, tearing down all internals and preparing for
	   graceful exit.
	*/
	Close()
}
