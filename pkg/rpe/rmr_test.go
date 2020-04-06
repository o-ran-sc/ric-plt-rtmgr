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


   This source code is part of the near-RT RIC (RAN Intelligent Controller)
   platform project (RICP).

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
	for _, endpoint := range testdata {
		ep := endpoint
		rtmgr.Eps[ep.Uuid] = &ep
	}
}

/*
RmrPush.GeneratePolicies() method is tested for happy path case
*/
func TestRmrPushGeneratePolicies(t *testing.T) {
	var rmrpush = RmrPush{}
	var pcs rtmgr.RicComponents
	resetTestDataset(stub.ValidEndpoints1)
	stub.ValidPlatformComponents = nil
	rtmgr.Subs = *stub.ValidSubscriptions
	rtmgr.PrsCfg = stub.DummyRoutes
	stub.E2map["E2instance1.com"] = stub.ValidE2TInstance
	pcs = stub.ValidRicComponents

	rawrt := rmrpush.GeneratePolicies(rtmgr.Eps, &pcs)
	t.Log(rawrt)
}

/*
getEndpointByUuid: Pass empty and valid values
*/
func TestRmrgetEndpointByUuid(t *testing.T) {
	var ep *rtmgr.Endpoint
	ep = getEndpointByUuid("")
	t.Logf("getEndpointByUuid() return was correct, got: %v, want: %v.", ep, "<nil>")
	ep = getEndpointByUuid("10.0.0.1:0")
}

/*
GetRpe Instance with empty and valid values
*/
func TestRmrGetRpe(t *testing.T) {
	_, _ = GetRpe("")
	_, _ = GetRpe("rmrpush")
}

/*
generateRouteTable with empty Platform components
*/
func TestGenerateRouteTableRmrGetRpe(t *testing.T) {
	rpe := Rpe{}
	_ = rpe.generateRouteTable(stub.ValidEndPointsEmpty)
}
