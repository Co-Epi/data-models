# CoEpi Go Server

The client-server flow is as follows:`
1. Users record symptoms in their CoEpi app, resulting in a POST to `/exposureandsymptoms` endpoint with their `Symptoms` and UUIDs/Dates of `Contacts`
2. All CoEpi apps poll every N mins with POST to `/exposurecheck` to see if the server has received any symptoms of the above matching those that the device has seen.

## CoEpi API Endpoints

Test Endpoint: https://coepi.wolk.com:8081

### POST `/exposureandsymptoms`
Body: JSON Object - `ExposureAndSymptoms`
```
{"Symptoms":"SlNPTkJMT0I6c2V2ZXJlIGZldmVyLGNvdWdoaW5n","Contacts":[{"UUID":"ax","Date":"2020-03-04"},{"UUID":"by","Date":"2020-03-15"},{"UUID":"cz","Date":"2020-03-20"}]}
```

Response:
```
200 OK
400 Error
```

Behavior:
* Server records ContactIDs (BLE Proximity Event UUIDs) in `contacts` table with { Date, SymptomHash }
* Server records Symptoms in `symptoms` table keyed by SymptomHash

### POST `/exposurecheck`
Body: JSON Object - `ExposureCheck`
Sample:
```
{"Contacts":[{"UUID":"by","Date":"2020-03-04"}]}
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

2. Check that you can write some contacts with `go test -run TestBackendSimple`

```
[root@d5 server]# go test -run TestBackendSimple

processExposureCheck(check1) SUCCESS: [JSONBLOB:severe fever,coughing]
processExposureCheck(check0) SUCCESS: []
PASS
ok	github.com/wolkdb/coepi-backend-go/server	0.412s

```

## Build + Run

```
$ make coepi
go build -o bin/coepi
Done building coepi.  Run "/bin/coepi" to launch coepi.
```

## Test

After getting your SSL Certs in the right spot with a DNS entry that matches and running `bin/coepi`, you can run this test:
```
# go test -run TestCoepi
ExposureAndSymptoms Sample: {"Symptoms":"SlNPTkJMT0I6c2V2ZXJlIGZldmVyLGNvdWdoaW5n","Contacts":[{"UUID":"ax","Date":"2020-03-04"},{"UUID":"by","Date":"2020-03-15"},{"UUID":"cz","Date":"2020-03-20"}]}
exposureandsymptoms[OK]exposurecheck(check1) SUCCESS: [JSONBLOB:severe fever,coughing]
exposurecheck(check0) SUCCESS: []
PASS
ok	github.com/wolkdb/coepi-backend-go	0.589s
```

which does the same things as the above backend test except going through the HTTP Server.

### How it works (at a glance)

In the `contacts` table, there is a map between each UUID and a `symptomHash`

```
# cbt read contacts
2020/03/21 02:36:20 -creds flag unset, will use gcloud credential
----------------------------------------
ax
  symptoms:2020-03-04                      @ 2020/03/21-02:34:55.307000
    "b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09"
----------------------------------------
by
  symptoms:2020-03-15                      @ 2020/03/21-02:34:55.389000
    "b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09"
----------------------------------------
cz
  symptoms:2020-03-20                      @ 2020/03/21-02:34:55.443000
    "b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09"
```

In the `symptoms` table, there is map between symptomHashes and a blob of bytes.
```
# cbt read symptoms
----------------------------------------
b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09
  case:symptoms                            @ 2020/03/21-02:03:55.604000
    "JSONBLOB:severe fever,coughing"
```
