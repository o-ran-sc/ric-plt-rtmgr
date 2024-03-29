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
	Mnemonic:	rmrpush_test.go
	Abstract:
	Date:		3 May 2019
*/
package sbi

import (
	//"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"os"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/stub"
	"testing"
	"time"
)

type Consumer struct{}

func (m Consumer) Consume(params *xapp.RMRParams) (err error) {
	xapp.SdlStorage.Store(rtmgr.RTMGR_SDL_NS, "myKey", params.Payload)
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
func resetTestPushDataset(instance RmrPush, testdata []rtmgr.Endpoint) {
	rtmgr.RMRConnStatus = make(map[string]bool)
	rtmgr.Eps = make(map[string]*rtmgr.Endpoint)
	for _, endpoint := range testdata {
		ep := endpoint
		//ep.Socket, _ = instance.NewSocket()
		rtmgr.Eps[ep.Uuid] = &ep
	}
}

/*
rmrpush.Initialize() method is empty, nothing to be tested
*/
func TestRmrPushInitialize(t *testing.T) {
	var rmrpush = RmrPush{}

	_ = rmrpush.Initialize("")
}

/*
rmrpush.Terminate() method is empty, nothing to be tested
*/
func TestRmrPushTerminate(t *testing.T) {
	var rmrpush = RmrPush{}

	_ = rmrpush.Terminate()
}

/*
rmrpush.UpdateEndpoints() is testd against stub.ValidXApps dataset
*/
func TestRmrPushUpdateEndpoints(t *testing.T) {
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints)

	rmrpush.UpdateEndpoints(&stub.ValidRicComponents)
	if rtmgr.Eps == nil {
		t.Errorf("rmrpush.UpdateEndpoints() result was incorrect, got: %v, want: %v.", nil, "rtmgr.Endpoints")
	}
}

/*
rmrpush.AddEndpoint() is tested for happy path case
*/
func TestRmrPushAddEndpoint(t *testing.T) {
	//	var err error
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints)
	_ = rmrpush.AddEndpoint(rtmgr.Eps["localhost"])
	/*	if err != nil {
		t.Errorf("rmrpush.AddEndpoint() return was incorrect, got: %v, want: %v.", err, "nil")
	}*/
}

/*
rmrpush.DistributeAll() is tested for happy path case
*/
func TestRmrPushDistributeAll(t *testing.T) {
	var err error
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints)

	rmrcallid = 200
	err = rmrpush.DistributeAll(stub.ValidPolicies)
	t.Log(err)
	/*if err != nil {
		t.Errorf("rmrpush.DistributeAll(policies) was incorrect, got: %v, want: %v.", err, "nil")
	}*/
}

/*
rmrpush.DistributeToEp() is tested for Sending case
*/
func TestDistributeToEp(t *testing.T) {
	var err error
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints)

	rmrdynamiccallid = 255
	err = rmrpush.DistributeToEp(stub.ValidPolicies, "localhost:4561", 100)
	if err != nil {
		t.Errorf("rmrpush.DistributetoEp(policies) was incorrect, got: %v, want: %v.", err, "nil")
	}
}

func TestDeleteEndpoint(t *testing.T) {
	var err error
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints)

	err = rmrpush.DeleteEndpoint(rtmgr.Eps["localhost"])
	if err != nil {
		t.Errorf("rmrpush.DeleteEndpoint() was incorrect, got: %v, want: %v.", err, "nil")
	}
}

func TestCheckEndpoint(t *testing.T) {
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints1)
	rmrpush.CheckEndpoint("192.168.0.1:0")
	rmrpush.CheckEndpoint("10.2.2.1:0")
	rmrpush.CheckEndpoint("localhost:0")
}

func TestCreateEndpoint(t *testing.T) {
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints1)
	rmrpush.CreateEndpoint("Src=127.0.0.1:4561 hello")
}

/*
Initialize and send policies
*/
func TestRmrPushInitializeandsendPolicies(t *testing.T) {
	var rmrpush = RmrPush{}
	resetTestPushDataset(rmrpush, stub.ValidEndpoints)
	policies := []string{"hello", "welcome"}
	rmrpush.send_data(rtmgr.Eps["localhost"], &policies, 1)
}

func TestString(t *testing.T) {
	var params xapp.RMRParams
	params.Payload = []byte("abcdefgh")
	params.Meid = &xapp.RMRMeid{}
	msg := RMRParams{&params}
	msg.String()

}

func TestSenddata(t *testing.T) {
	var rmrpush = RmrPush{}
	ep := rtmgr.Endpoint{Whid: -1, Ip: "1.1.1.1"}
	policies := []string{"mse|12345|-1|local.com"}
	rmrpush.send_data(&ep, &policies, 300)
}

func TestSendDynamicdata(t *testing.T) {
	var rmrpush = RmrPush{}
	ep := "1.1.1.1"
	policies := []string{"mse|12345|-1|local.com"}
	rmrpush.sendDynamicRoutes(ep, 1, &policies, 300)
}
