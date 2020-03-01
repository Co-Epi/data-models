package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestBackendSimple(t *testing.T) {
	backend, err := NewBackend(DefaultProject, DefaultInstance)
	if err != nil {
		t.Fatalf("%s", err)
	}

	eas := new(ExposureAndSymptoms)
	eas.Contacts = []Contact{Contact{UUID: "ax", Date: "2020-03-04"}, Contact{UUID: "by", Date: "2020-03-15"}, Contact{UUID: "cz", Date: "2020-03-20"}}
	eas.Symptoms = []byte("JSONBLOB:severe fever,coughing")
	easJSON, err := json.Marshal(eas)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	fmt.Printf("ExposureAndSymptoms Sample: %s\n", easJSON)

	err = backend.ProcessExposureAndSymptoms(eas)
	if err != nil {
		t.Fatalf("ProcessExposureAndSymptoms: %s", err)
	}

	check1 := new(ExposureCheck)
	check1.Contacts = []Contact{Contact{UUID: "by", Date: "2020-03-04"}}
	check1JSON, err := json.Marshal(check1)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	fmt.Printf("ExposureCheck Sample: %s\n", check1JSON)

	symptoms, err := backend.ProcessExposureCheck(check1)
	if err != nil {
		t.Fatalf("ProcessExposureCheck(check1): %s", err)
	}
	if len(symptoms) != 1 {
		t.Fatalf("ProcessExposureCheck(check1) Expected 1 response, got %d", len(symptoms))
	}

	if !bytes.Equal(eas.Symptoms, symptoms[0]) {
		t.Fatalf("ProcessExposureCheck(check1) Expected 1 response, got %d", len(symptoms))
	}
	fmt.Printf("ProcessExposureCheck(check1) SUCCESS: [%s]\n", symptoms[0])

	check0 := new(ExposureCheck)
	check0.Contacts = []Contact{Contact{UUID: "00", Date: "2020-03-21"}}
	symptoms, err = backend.ProcessExposureCheck(check0)
	if err != nil {
		t.Fatalf("ProcessExposureCheck(check0): %s", err)
	}
	if len(symptoms) != 0 {
		t.Fatalf("processExposureCheck(check0) Expected 0 responses, but got %d", len(symptoms))
	}
	fmt.Printf("processExposureCheck(check0) SUCCESS: []\n")
}

// TODO: setup N random users and have them generate Contact records between themselves in N goroutines, each doing a ExposureCheck every few seconds
//  Then have N/10 users post a ExposureAndSymptoms, which should have some of the go routines generating symptom responses
