package daemon

import (
	"os"
)

//this package will be used to maintain any running processes to maintain the database

//the main function should be used to spawn maintenance goroutines which run indefinitely

func Main() {
	err := make(chan error, 1)




	go func (c chan error) {
		c <- momentRetriever()
	}(err)



	msg := <- err

	panic(msg)
	os.Exit(1)
}