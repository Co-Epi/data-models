##Exposure and Symptoms Record (POST)
Request:
- Symptoms record (Answers from symptoms interview)
- Datestamp
 - Geohash(es) covering the geo-location of some/all exposure records
 - Exposure Record (List of all BLE Proximity Event UUIDs within that geohash)
Response: just a confirmation code or error message

Server assembles all exposure records for each datestamp and (optionally, as we scale) shards them by geohash
- each BLE Proximity Event UUIDs is truncated to include only the first 64 bits of the version-4 UUIDs

##Exposure Check (GET)
- All Exposure Records (with shortened UUIDs) for a given date (and optionally geohash)
 - Date (and optionally geohash) are components of the API GET URI
- As we scale, the response should have cache-control headers with an appropriate TTL (maybe 1h) and be cached so it doesn’t need to be pulled from the DB every time.

Client checks its own local record of all BLE Proximity Event UUIDs against those in the Exposure Check response
If a match is found, the client alerts the user, displays date/time and location of the match, and requests permission to share exposure(s) and request case record(s)

##Confirmed Exposure (POST)
Request:
- Full BLE Proximity Event UUID(s) (all 128 bits)
Response:
- Symptoms Record, including reported symptoms + severity, dates, etc.

