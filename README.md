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

2. Check that you can write some contacts with `go test -run TestBackendSimple`

```
# go test -run TestBackendSimple
PASS
ok  	github.com/wolkdb/coepi-backend-go/server	0.522s

# cbt read symptoms
----------------------------------------
b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09
  case:symptoms                            @ 2020/03/21-02:03:55.604000
    "JSONBLOB:severe fever,coughing"

# cbt read contacts
----------------------------------------
ax
  symptoms:2020-03-04                      @ 2020/03/21-02:03:55.750000
    "\xb9:\x90\xa8C\xed)5\"\xaa\x807\x81)\x8d\xacC`@\xfa#\x1a\x18\x9eRƙJ]Y\x1f\t"
----------------------------------------
by
  symptoms:2020-03-15                      @ 2020/03/21-02:03:55.923000
    "\xb9:\x90\xa8C\xed)5\"\xaa\x807\x81)\x8d\xacC`@\xfa#\x1a\x18\x9eRƙJ]Y\x1f\t"
----------------------------------------
cz
  symptoms:2020-03-20                      @ 2020/03/21-02:03:55.983000
    "\xb9:\x90\xa8C\xed)5\"\xaa\x807\x81)\x8d\xacC`@\xfa#\x1a\x18\x9eRƙJ]Y\x1f\t"
```
