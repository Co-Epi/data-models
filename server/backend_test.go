package main

import (
	"bytes"
	"fmt"
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
		t.Fatalf("processExposureAndSymptoms: %s", err)
	}

	check1 := new(ExposureCheck)
	check1.Contacts = []Contact{Contact{UUID: "by", Date: "2020-03-04"}}
	symptoms, err := backend.processExposureCheck(check1)
	if err != nil {
		t.Fatalf("processExposureCheck(check1): %s", err)
	}
	if len(symptoms) != 1 {
		t.Fatalf("processExposureCheck(check1) Expected 1 response, got %d", len(symptoms))
	}

	if !bytes.Equal(eas.Symptoms, symptoms[0]) {
		t.Fatalf("processExposureCheck(check1) Expected 1 response, got %d", len(symptoms))
	}
	fmt.Printf("processExposureCheck(check1) SUCCESS: [%s]\n", symptoms[0])

	check0 := new(ExposureCheck)
	check0.Contacts = []Contact{Contact{UUID: "00", Date: "2020-03-21"}}
	symptoms, err = backend.processExposureCheck(check0)
	if err != nil {
		t.Fatalf("processExposureCheck(check0): %s", err)
	}
	if len(symptoms) != 0 {
		t.Fatalf("processExposureCheck(check0) Expected 0 responses, but got %d", len(symptoms))
	}
	fmt.Printf("processExposureCheck(check0) SUCCESS: []\n")
}

// TODO: setup N random users and have them generate Contact records between themselves in N goroutines, each doing a ExposureCheck every few seconds
//  Then have N/10 users post a ExposureAndSymptoms, which should have some of the go routines generating symptom responses