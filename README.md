# CoEpi API

## POST /exposureandsymptoms
Body: JSON Object - ExposureAndSymptoms
Sample:
```
TODO
```

Response:
```
200 OK
400 error
```

Behavior:
* Server records ContactIDs (BLE Proximity Event UUIDs) in `contacts` table with { Date, SymptomHash }
* Server records Symptoms in `symptoms` table keyed by SymptomHash

## POST /exposurecheck
Body: JSON Object - ExposureCheck
Sample:
```
TODO
```
Behavior:
* Server gets UUIDs (BLE Proximity Event UUIDs) and checks `contacts` table for potential { Date, SymptomHash } combinations
* With hits, Server fetches Symptoms from `symptoms` table keyed by SymptomHash and returns array of byte blobs


### BigTable Setup

Use `cbt` https://cloud.google.com/bigtable/docs/quickstart-cbt

1. After setting up your BigTable `co-epi` instance and updating the project/instance strings in server.go, create 2 tables and their families with `cbt`
```
cbt createtable contacts
cbt createtable symptoms
cbt createfamily symptoms case
cbt createfamily contacts symptoms
```
Check with:
```
# cbt ls
contacts
symptoms
# cbt ls symptoms
Family Name	GC Policy
-----------	---------
case		<never>
# cbt ls contacts
Family Name	GC Policy
-----------	---------
symptoms	<never>
```

2. Check that you can write some contacts
`go test TestBackendSimple`
