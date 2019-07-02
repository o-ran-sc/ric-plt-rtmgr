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
func createNewStubPubSocket() (NngSocket, error) {
	socket := stub.MangosSocket{}
	return socket, nil
}

/*
Returns a SocketError
*/
func createNewStubPubSocketError() (NngSocket, error) {
	return nil, errors.New("stub generated Create Socket error")
}

/*
Returns a Socket which always generates error on Close()
*/
func createNewStubPubSocketCloseError() (NngSocket, error) {
	socket := stub.MangosSocket{}
	socket.GenerateSocketCloseError = true
	return socket, nil
}

/*
Returns a Socket which always generates error on Listen()
*/
func createNewStubPubSocketListenError() (NngSocket, error) {
	socket := stub.MangosSocket{}
	socket.GenerateSocketListenError = true
	return socket, nil
}

/*
Returns a Socket which always generates error on Send()
*/
func createNewStubPubSocketSendError() (NngSocket, error) {
	socket := stub.MangosSocket{}
	socket.GenerateSocketSendError = true
	return socket, nil
}

/*
Resets the EndpointList according to argumnets
*/
func resetTestPubDataset(instance NngPub, testdata []rtmgr.Endpoint) {
	rtmgr.Eps = make(map[string]*rtmgr.Endpoint)
	for _, endpoint := range testdata {
		ep := endpoint
		ep.Socket, _ = instance.NewSocket()
		rtmgr.Eps[ep.Uuid] = &ep
	}
}

/*
nngPub.Initialize() method is tested for happy path case
*/
func TestNngPubInitialize(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocket

	err := nngpub.Initialize("")
	if err != nil {
		t.Errorf("nngPub.Initialize() was incorrect, got: %v, want: %v.", err, nil)
	}
}

/*
nngPub.Initialize() is tested for Socket creating error case
*/
func TestNngPubInitializeWithSocketError(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocketError

	err := nngpub.Initialize("")
	if err == nil {
		t.Errorf("nngPub.Initialize() was incorrect, got: %v, want: %v.", nil, "error")
	}
}

/*
nngPub.Initialize() is tested for Socket listening error case
*/
func TestNngPubInitializeWithSocketListenError(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocketListenError

	err := nngpub.Initialize("")
	if err == nil {
		t.Errorf("nngPub.Initialize() was incorrect, got: %v, want: %v.", nil, "error")
	}
}

/*
nngPub.Terminate() method is empty, nothing to be tested
*/
func TestNngPubTerminate(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocket
	nngpub.Initialize("")

	err := nngpub.Terminate()
	if err != nil {
		t.Errorf("nngPub.Terminate() was incorrect, got: %v, want: %v.", err, nil)
	}
}

/*
nngPub.Terminate() is tested for Socket closing error case
*/
func TestNngPubTerminateWithSocketCloseError(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocketCloseError
	nngpub.Initialize("")

	err := nngpub.Terminate()
	if err == nil {
		t.Errorf("nngPub.Terminate() was incorrect, got: %v, want: %v.", nil, "error")
	}
}

/*
nngPub.UpdateEndpoints() is testd against stub.ValidXapps dataset
*/
func TestNngPubUpdateEndpoints(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocket
	nngpub.Initialize("")
	rtmgr.Eps = make(rtmgr.Endpoints)
	nngpub.UpdateEndpoints(&stub.ValidRicComponents)
	if rtmgr.Eps == nil {
		t.Errorf("nngPub.UpdateEndpoints() result was incorrect, got: %v, want: %v.", nil, "rtmgr.Endpoints")
	}
}

/*
nngPub.AddEndpoint() method is empty, nothing to be tested
*/
func TestNngPubAddEndpoint(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocket

	_ = nngpub.AddEndpoint(new(rtmgr.Endpoint))
}

/*
nngPub.DeleteEndpoint() method is empty, nothing to be tested
*/
func TestNngPubDeleteEndpoint(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocket

	_ = nngpub.DeleteEndpoint(new(rtmgr.Endpoint))
}

/*
nngPub.DistributeAll() is tested for happy path case
*/
func TestNngPubDistributeAll(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocket
	nngpub.Initialize("")
	resetTestPubDataset(nngpub, stub.ValidEndpoints)

	err := nngpub.DistributeAll(stub.ValidPolicies)
	if err != nil {
		t.Errorf("nngPub.DistributeAll(policies) was incorrect, got: %v, want: %v.", err, nil)
	}
}

/*
nngPub.DistributeAll() is tested for Sending error case
*/
func TestNngPubDistributeAllSocketSendError(t *testing.T) {
	var nngpub = NngPub{}
	nngpub.NewSocket = createNewStubPubSocketSendError
	nngpub.Initialize("")
	resetTestPubDataset(nngpub, stub.ValidEndpoints)

	err := nngpub.DistributeAll(stub.ValidPolicies)
	if err == nil {
		t.Errorf("nngPub.DistributeAll(policies) was incorrect, got: %v, want: %v.", nil, "error")
	}
}
