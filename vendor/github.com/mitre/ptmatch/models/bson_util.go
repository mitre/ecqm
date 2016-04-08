/*
Copyright 2016 The MITRE Corporation. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import (
	"errors"

	logger "github.com/mitre/ptmatch/logger"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func ToBsonObjectID(idStr string) (bson.ObjectId, error) {
	var id bson.ObjectId

	logger.Log.WithFields(logrus.Fields{"id": idStr}).Debug("ToBsonObjectID")
	if bson.IsObjectIdHex(idStr) {
		id = bson.ObjectIdHex(idStr)
	} else {
		return bson.ObjectId(0), errors.New("Invalid value: " + idStr)
	}
	return id, nil
}