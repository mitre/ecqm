# Quality Reports

## POST /QualityReport

Create a new quality report

* Content-Type: "application/x-www-form-urlencoded"
* Accept: "application/json"

```
measureId=abcd&effectiveDate=1234
```
===
### Response

* Status: 200
* Content-Type: application/json; charset=utf-8
* Data.measureId = "abcd"
* Data.id: /[\da-f]{24}/

## GET /QualityReport/56bd06841cd462774f2af485

Get a new quality report

* Accept: "application/json"

===
### Response

* Status: 200
* Content-Type: application/json; charset=utf-8
* Data.measureId = "efg"
* Data.id: "56bd06841cd462774f2af485"

## GET /QualityReport/56bd06841cd462774f2af485/PatientResults?population=initialPatientPopulation

Get all people in the initial patient population for this measure.

* Accept: "application/json"

===
### Response

* Status: 200
* Content-Type: application/json; charset=utf-8
* Data.total = 2