package main

import (
	"log"
	"testing"
)

var a App

// M struct below contains the method to run other tests in the package

func TestMain(m *testing.M) {
	err := app.Initialise()

	if err != nil {
		log.Fatal("Error occured while initialising database")
	}

	m.Run() // RUn all tests within the package
}
