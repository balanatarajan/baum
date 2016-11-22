package main

import (
	"config"
	"fmt"
)

type dbConfig struct {
	location string
}

type fileConfig struct {
	location string
}

func main() {

	// Read our config
	c, err := ReadConfig("storage.cdb")

	var a infra.Acceptor
	err := a.Open()

	defer a.Close()

	// Run the event loop
	err = a.Run()

	fmt.Println("About to Wait")

	// Wait for all go routines to clean up before an exit
	err = a.Wait()
}
