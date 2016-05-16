# ecqm [![Build Status](https://api.travis-ci.org/mitre/ecqm.svg)](https://travis-ci.org/mitre/ecqm)

This project contains the components to serve information about electronic
Clinical Quality Measures. It is assumed that this project will be used with the
[Go based FHIR server](http://https://github.com/intervention-engine/fhir) for
the storage of patient information. Clinical quality measure calculation will be
handled by the
[node.js based quality measure engine](https://github.com/mitre/node-qme).

This library provides a RESTful JSON API for accessing information on quality
measures as well as running calculations and accessing their results. This can
also be used with the [eCQM Frontend](https://github.com/mitre/ecqm-frontend)
for a web based interface to this information.

This project also includes the services for the [Patient Matching Test Harness](https://github.com/mitre/ptmatch).
That allows the project to calculate quality measures, function as a patient matching test harness or work
as a plain FHIR server.

## Requirements

* Git
* Go >= 1.6
* Ruby >= 2.2
* MongoDB >= 3.2

## Install

Get the code:

```
mkdir -p $GOPATH/src/github.com/mitre
cd $GOPATH/src/github.com/mitre
git clone https://github.com/mitre/ecqm.git
```

This project currently depends on resources from the
[health-data-standards](https://github.com/projectcypress/health-data-standards)
library to pull in measure bundles. The following should pull over enough of
that library to get started:

```
cd ecqm
gem install health-data-standards
gem install highline
rake bundle:download_and_install
```

This requires an NLM username and password to obtain a
[Cypress bundle](http://projectcypress.org/test_data.html) which
contains the measures.

## Running and Testing

This project uses [Godep](https://github.com/tools/godep) to manage dependencies. All of the needed related
libraries are included in the vendor directory.

To run all of the tests for this project, run:

    go test $(go list ./... | grep -v /vendor/)

in this directory.

To start the server, run:
    
    go run server.go -assets PATH_TO_ASSETS

In this case, PATH_TO_ASSETS should be a location where a version of either the
[eCQM Frontend](https://github.com/mitre/ecqm-frontend) or [Patient Match Frontend](https://github.com/mitre/ptmatch-frontend)
has been built.

## License

Copyright 2016 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.