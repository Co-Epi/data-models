package main

import (
	"testing"
)

const (
	defaultProject  = "us-west1-wlk"
	defaultInstance = "co-epi"
)

func TestBackendSimple(t *testing.T) {
	backend, err := NewBackend(defaultProject, defaultInstance)
	if err != nil {
		t.Fatalf("%s", err)
	}

	eas := new(ExposureAndSymptoms)
	eas.Contacts = []Contact{Contact{UUID: "ax", Date: "2020-03-04"}, Contact{UUID: "by", Date: "2020-03-15"}, Contact{UUID: "cz", Date: "2020-03-20"}}
	eas.Symptoms = []byte("JSONBLOB:severe fever,coughing")

	err = backend.processExposureAndSymptoms(eas)
	if err != nil {
		t.Fatalf("%s", err)
	}

	// TODO: setup N random users and have them generate Contact records between themselves in N goroutines, each doing a ExposureCheck every few seconds
	//  Then have N/10 users post a ExposureAndSymptoms, which should have some of the go routines generating symptom responses

}
