# CoEpi Go Server 

## POST `/exposureandsymptoms`
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

## POST `/exposurecheck`
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

In the `symptoms` table, there is map between symptomHashes and a blob of bytes.
```
# cbt read symptoms
----------------------------------------
b93a90a843ed293522aa803781298dac436040fa231a189e52c6994a5d591f09
  case:symptoms                            @ 2020/03/21-02:03:55.604000
    "JSONBLOB:severe fever,coughing"
```

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


## How it works

