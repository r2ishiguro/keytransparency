// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"errors"

	ktpb "github.com/google/keytransparency/core/proto/keytransparency_v1_types"
	"github.com/google/trillian"
)

var (
	// ErrAlreadyStored is raised if the caller tries storing a response which
	// has already been stored.
	ErrAlreadyStored = errors.New("already stored epoch")
	// ErrNotFound is raised if the caller tries to retrieve data for an epoch
	// which has not been processed and stored yet.
	ErrNotFound = errors.New("data for epoch not found")
)

// MonitoringResult stores all data
type MonitoringResult struct {
	// Smr contains the map root signed by the monitor in case all verifications
	// have passed.
	Smr *trillian.SignedMapRoot
	// Seen is the the unix timestamp at which the mutations response has been
	// received.
	Seen int64
	// Errors contains a string representation of the verifications steps that
	// failed.
	Errors []error
	// Response contains the original mutations API response from the server
	// in case at least one verification step failed.
	Response *ktpb.GetMutationsResponse
}

type Storage interface {
	// Set internally stores the given data as a MonitoringResult which can be
	// retrieved by Get.
	Set(epoch int64, seenNanos int64, smr *trillian.SignedMapRoot, response *ktpb.GetMutationsResponse, errorList []error) error
	// Get returns the MonitoringResult for the given epoch. It returns an error
	// if the result does not exist.
	Get(epoch int64) (*MonitoringResult, error)
	// LatestEpoch is a convenience method to retrieve the latest stored epoch.
	LatestEpoch() int64
}
