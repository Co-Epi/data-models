package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigtable"
	"github.com/wolkdb/cloudstore/wolk/core"
)

// Backend holds a client to connect to the BigTable backend
type Backend struct {
	client *bigtable.Client
}

// Contact represents a BLE pairing between 2 devices
type Contact struct {
	ContactID string // this is the hash of a pair of BLE ids
	Date      string // this is when the pair came into contact, used to set TTLs
}

// ExposureAndSymptoms payload is sent by client to /exposureandsymptoms when user reports symptoms
type ExposureAndSymptoms struct {
	Symptoms []byte    // this is expected to be a JSON blob but the server doesn't need to parse it
	Contacts []Contact // these are the contacts th
}

// ExposureCheck payload is sent by client to /exposurecheck to try to
type ExposureCheck struct {
	Contacts []Contact
}

// TableContacts stores the mapping between contactIDs and symptomHash.
const TableContacts = "contacts"

// TableSymptoms stores the mapping between symptomHash and symptoms.   The content of the symptoms string is a JSON document that clients need to power the UI but the server does not need to process it
const TableSymptoms = "symptoms"

// NewBackend sets up a client connection to BigTable to manage incoming payloads
func NewBackend(project, instance string) (backend *Backend, err error) {
	backend = new(Backend)
	client, err := bigtable.NewClient(context.Background(), project, instance)
	if err != nil {
		return backend, err
	}
	backend.client = client
	return backend, nil
}

// POST /exposureandsymptoms
//  Input: ExposureAndSymptoms
//  Output: error
func (backend *Backend) processExposureAndSymptoms(payload *ExposureAndSymptoms) (err error) {
	// store symptoms in the symptoms table
	symptoms := payload.Symptoms
	symptomsTable := backend.client.Open(TableSymptoms)
	symptomsHash := core.Computehash(symptoms)
	mut := bigtable.NewMutation()
	mut.Set("case", "symptoms", bigtable.Now(), []byte(symptoms))
	err = symptomsTable.Apply(context.Background(), fmt.Sprintf("%x", symptomsHash), mut)
	if err != nil {
		return err
	}

	contactsTable := backend.client.Open(TableContacts)
	// store the first 64 one cell per observation
	for _, contact := range payload.Contacts {
		mut := bigtable.NewMutation()
		mut.Set("symptoms", contact.Date, bigtable.Now(), symptomsHash)
		err = contactsTable.Apply(context.Background(), contact.ContactID, mut)
		if err != nil {
			return err
		}
	}
	return nil
}

// POST /exposurecheck
//  Input: ExposureCheck
//  Output: array of byte blobs
func (backend *Backend) processExposureCheck(payload *ExposureCheck) (symptomsList [][]byte, err error) {
	tableContacts := backend.client.Open(TableContacts)
	// store one cell per observation
	symptomsHashes := make([]string, 0)
	for _, contact := range payload.Contacts {
		rr := bigtable.PrefixRange(contact.ContactID)
		err := tableContacts.ReadRows(context.Background(), rr, func(r bigtable.Row) bool {
			for k, xv := range r {
				switch k {
				case "symptoms":
					for _, yv := range xv {
						// TODO: get the date from yv.Column ("symptoms:2020-03-22") and filter it
						// MATCH DETECTED
						symptomsHashes = append(symptomsHashes, string(yv.Value))
					}
				}
			}
			return true // Keep going.
		}, bigtable.RowFilter(bigtable.FamilyFilter("links")))
		if err != nil {
			// TODO: handle err.
		}
	}

	// For all symptomHashes, get Symptoms byte blobs (containing reported symptoms + severity, dates, etc.)
	symptomsList = make([][]byte, 0)
	for _, symptomsHash := range symptomsHashes {
		tableSymptoms := backend.client.Open(TableSymptoms)
		rr := bigtable.PrefixRange(symptomsHash)
		err := tableSymptoms.ReadRows(context.Background(), rr, func(r bigtable.Row) bool {
			for k, xv := range r {
				switch k {
				case "case":
					for _, yv := range xv {
						if yv.Column == "case:symptoms" {
							symptomsList = append(symptomsList, yv.Value)
						}
					}
				}
			}
			return true // Keep going.
		}, bigtable.RowFilter(bigtable.FamilyFilter("case")))
		if err != nil {
			// TODO: handle err.
		}
	}
	return symptomsList, nil
}
