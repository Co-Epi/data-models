package main

import (
	"testing"
)

func TestBackend(t *testing.T) {
  // TODO: setup N random users and have them generate Contact records between themselves in N goroutines, each doing a ExposureCheck every few seconds
  //  Then have N/10 users post a ExposureAndSymptoms, which should have some of the go routines generating symptom responses   
}
