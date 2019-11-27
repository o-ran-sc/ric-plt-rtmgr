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
	"errors"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/stub"
	"testing"
)

/*
Returns an error free Socket instance
*/
func createNewStubPushSocket() (NngSocket, error) {
	socket := stub.MangosSocket{}
	return socket, nil
}

/*
Returns a SocketError
*/
func createNewStubPushSocketError() (NngSocket, error) {
	return nil, errors.New("stub generated Create Socket error")
}

/*
Returns a Socket which always generates error on Close()
*/
func createNewStubPushSocketCloseError() (NngSocket, error) {
	socket := stub.MangosSocket{}
	socket.GenerateSocketCloseError = true
	return socket, nil
}

/*
Returns a Socket which always generates error on Send()
*/
func createNewStubPushSocketSendError() (NngSocket, error) {
	socket := stub.MangosSocket{}
	socket.GenerateSocketSendError = true
	return socket, nil
}

/*
Returns a Socket which always generates error on Dial()
*/
func createNewStubPushSocketDialError() (NngSocket, error) {
	socket := stub.MangosSocket{}
	socket.GenerateSocketDialError = true
	return socket, nil
}

/*
Resets the EndpointList according to argumnets
*/
func resetTestPushDataset(instance NngPush, testdata []rtmgr.Endpoint) {
	rtmgr.Eps = make(map[string]*rtmgr.Endpoint)
	for _, endpoint := range testdata {
		ep := endpoint
		ep.Socket, _ = instance.NewSocket()
		rtmgr.Eps[ep.Uuid] = &ep
	}
}

/*
nngpush.Initialize() method is empty, nothing to be tested
*/
func TestNngPushInitialize(t *testing.T) {
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocket

	_ = nngpush.Initialize("")
}

/*
nngpush.Terminate() method is empty, nothing to be tested
*/
func TestNngPushTerminate(t *testing.T) {
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocket

	_ = nngpush.Terminate()
}

/*
nngpush.UpdateEndpoints() is testd against stub.ValidXApps dataset
*/
func TestNngPushUpdateEndpoints(t *testing.T) {
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocket
	rtmgr.Eps = make(rtmgr.Endpoints)

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
	nngpush.NewSocket = createNewStubPushSocket
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.AddEndpoint(rtmgr.Eps["10.0.0.1:0"])
	if err != nil {
		t.Errorf("nngpush.AddEndpoint() return was incorrect, got: %v, want: %v.", err, "nil")
	}
	if rtmgr.Eps["10.0.0.1:0"].Socket == nil {
		t.Errorf("nngpush.AddEndpoint() was incorrect, got: %v, want: %v.", nil, "Socket")
	}
}

/*
nngpush.AddEndpoint() is tested for Socket creating error case
*/
func TestNngPushAddEndpointWithSocketError(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocketError
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.AddEndpoint(rtmgr.Eps["10.0.0.1:0"])
	if err == nil {
		t.Errorf("nngpush.AddEndpoint() was incorrect, got: %v, want: %v.", nil, "error")
	}
	if rtmgr.Eps["10.0.0.1:0"].Socket != nil {
		t.Errorf("nngpush.AddEndpoint() was incorrect, got: %v, want: %v.", rtmgr.Eps["10.0.0.1:0"].Socket, nil)
	}
}

/*
nngpush.AddEndpoint() is tested for Dialing error case
*/
func TestNngPushAddEndpointWithSocketDialError(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocketDialError
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.AddEndpoint(rtmgr.Eps["10.0.0.1:0"])
	if err == nil {
		t.Errorf("nngpush.AddEndpoint() was incorrect, got: %v, want: %v.", nil, "error")
	}
}

/*
nngpush.DistributeAll() is tested for happy path case
*/
func TestNngPushDistributeAll(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocket
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.DistributeAll(stub.ValidPolicies)
	if err != nil {
		t.Errorf("nngpush.DistributeAll(policies) was incorrect, got: %v, want: %v.", err, "nil")
	}
}

/*
nngpush.DistributeAll() is tested for Sending error case
*/
func TestNngPushDistributeAllSocketSendError(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocketSendError
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.DistributeAll(stub.ValidPolicies)
	if err != nil {
		t.Errorf("nngpush.DistributeAll(policies) was incorrect, got: %v, want: %v.", err, "nil")
	}
}

func TestNngPushDeleteEndpoint(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocket
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.DeleteEndpoint(rtmgr.Eps["10.0.0.1:0"])
	if err != nil {
		t.Errorf("nngpush.DeleteEndpoint() was incorrect, got: %v, want: %v.", err, "nil")
	}
}

func TestNngPushDeleteEndpointWithSocketCloseError(t *testing.T) {
	var err error
	var nngpush = NngPush{}
	nngpush.NewSocket = createNewStubPushSocketCloseError
	resetTestPushDataset(nngpush, stub.ValidEndpoints)

	err = nngpush.DeleteEndpoint(rtmgr.Eps["10.1.1.1:0"])
	if err == nil {
		t.Errorf("nngpush.DeleteEndpoint() was incorrect, got: %v, want: %v.", nil, "error")
	}
}
