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
	Mnemonic:	nngpush_test.go
	Abstract:
	Date:		3 May 2019
*/
package sbi

import (
	//"errors"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/stub"
	"time"
	"os"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"testing"
)

type Consumer struct{}

func (m Consumer) Consume(params *xapp.RMRParams) (err error) {
        xapp.Sdl.Store("myKey", params.Payload)
        return nil
}

// Test cases
func TestMain(m *testing.M) {
        go xapp.RunWithParams(Consumer{}, false)
        time.Sleep(time.Duration(5) * time.Second)
        code := m.Run()
        os.Exit(code)
}

/*
Resets the EndpointList according to argumnets
*/
func resetTestPushDataset(instance NngPush, testdata []rtmgr.Endpoint) {
	rtmgr.Eps = make(map[string]*rtmgr.Endpoint)
	for _, endpoint := range testdata {
		ep := endpoint
		//ep.Socket, _ = instance.NewSocket()
		rtmgr.Eps[ep.Uuid] = &ep
	}
}

/*
nngpush.Initialize() method is empty, nothing to be tested
*/
func TestNngPushInitialize(t *testing.T) {
	var nngpush = NngPush{}

	_ = nngpush.Initialize("")
}

/*
nngpush.Terminate() method is empty, nothing to be tested
*/
func TestNngPushTerminate(t *testing.T) {
	var nngpush = NngPush{}

	_ = nngpush.Terminate()
}

/*
nngpush.UpdateEndpoints() is testd against stub.ValidXApps dataset
*/
func TestNngPushUpdateEndpoints(t *testing.T) {
	var nngpush = NngPush{}
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	nngpush.UpdateEndpoints(&stub.ValidRicComponents)
	if rtmgr.Eps == nil {
		t.Errorf("nngpush.UpdateEndpoints() result was incorrect, got: %v, want: %v.", nil, "rtmgr.Endpoints")
	}
}

/*
nngpush.AddEndpoint() is tested for happy path case
*/
func TestNngPushAddEndpoint(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	resetTestPushDataset(nngpush, stub.ValidEndpoints)
	err = nngpush.AddEndpoint(rtmgr.Eps["localhost"])
	if err != nil {
		t.Errorf("nngpush.AddEndpoint() return was incorrect, got: %v, want: %v.", err, "nil")
	}
}


/*
nngpush.DistributeAll() is tested for happy path case
*/
func TestNngPushDistributeAll(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.DistributeAll(stub.ValidPolicies)
	if err != nil {
		t.Errorf("nngpush.DistributeAll(policies) was incorrect, got: %v, want: %v.", err, "nil")
	}
}

/*
nngpush.DistributeToEp() is tested for Sending case
*/
func TestDistributeToEp(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.DistributeToEp(stub.ValidPolicies,rtmgr.Eps["localhost"])
	if err != nil {
		t.Errorf("nngpush.DistributetoEp(policies) was incorrect, got: %v, want: %v.", err, "nil")
	}
}

func TestDeleteEndpoint(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.DeleteEndpoint(rtmgr.Eps["localhost"])
	if err != nil {
		t.Errorf("nngpush.DeleteEndpoint() was incorrect, got: %v, want: %v.", err, "nil")
	}
}

func TestCreateEndpoint(t *testing.T) {
	var nngpush = NngPush{}
	resetTestPushDataset(nngpush, stub.ValidEndpoints1)
	 nngpush.CreateEndpoint("192.168.0.1:0")
	 nngpush.CreateEndpoint("localhost:4560")
}
/*
Initialize and send policies
*/
func TestNngPushInitializeandsendPolicies(t *testing.T) {
        var nngpush = NngPush{}
	resetTestPushDataset(nngpush, stub.ValidEndpoints)
        policies := []string{"hello","welcome"}
	nngpush.send(rtmgr.Eps["localhost"],&policies)
}
