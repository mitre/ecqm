# ecqm

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

## Requirements

* Git
* Go >= 1.5
* Ruby >= 2.2
* MongoDB >= 3.2

## Install

_Note_: These instructions are currently broken on MongoDB 3.2.0 and 3.2.1 due
to [an issue with inserting large documents](https://jira.mongodb.org/browse/SERVER-22167).
We will need to wait for MongoDB 3.2.2 to be released for these to fully work.


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
