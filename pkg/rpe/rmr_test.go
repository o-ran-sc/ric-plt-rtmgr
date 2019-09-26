/*
==================================================================================
   Copyright (c) 2019 AT&T Intellectual Property.
   Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/
/*
	Mnemonic:	nngpub_test.go
	Abstract:
	Date:		25 April 2019
*/
package rpe

import (
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/stub"
	"testing"
)

func resetTestDataset(testdata []rtmgr.Endpoint) {
	rtmgr.Eps = make(map[string]*rtmgr.Endpoint)
	for _, endpoint := range stub.ValidEndpoints {
		ep := endpoint
		rtmgr.Eps[ep.Uuid] = &ep
	}
}

/*
RmrPush.GeneratePolicies() method is tested for happy path case
*/
func TestRmrPushGeneratePolicies(t *testing.T) {
	var rmrpush = RmrPush{}
	resetTestDataset(stub.ValidEndpoints)

	rawrt := rmrpush.GeneratePolicies(rtmgr.Eps)
	t.Log(rawrt)
}
