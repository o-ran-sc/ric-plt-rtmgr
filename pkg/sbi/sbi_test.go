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
	Mnemonic:	sbi_test.go
	Abstract:
	Date:		25 April 2019
*/
package sbi

import (
	"errors"
	"reflect"
	"routing-manager/pkg/rtmgr"
	"testing"
)

func TestGetSbi(t *testing.T) {
	var errtype = errors.New("")
	var sbitype = new(RmrPush)
	var invalids = []string{"rmrpus", ""}

	sbii, err := GetSbi("rmrpush")
	if err != nil {
		t.Errorf("GetSbi(rmrpub) was incorrect, got: %v, want: %v.", reflect.TypeOf(err), nil)
	}
	if reflect.TypeOf(sbii) != reflect.TypeOf(sbitype) {
		t.Errorf("GetSbi(rmrpub) was incorrect, got: %v, want: %v.", reflect.TypeOf(sbii), reflect.TypeOf(sbitype))
	}

	for _, arg := range invalids {
		_, err := GetSbi(arg)
		if err == nil {
			t.Errorf("GetSbi("+arg+") was incorrect, got: %v, want: %v.", reflect.TypeOf(err), reflect.TypeOf(errtype))
		}
	}
}

func TestUpdateE2TendPoint(t *testing.T) {
	var err error
	var sbi = Sbi{}
	sbii, err := GetSbi("rmrpush")

	var EP = make(map[string]*rtmgr.Endpoint)
	EP["127.0.0.2"] = &rtmgr.Endpoint{Uuid: "127.0.0.2", Name: "E2TERM", XAppType: "app1", Ip: "127.0.0.2", Port: 4562, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: true, Keepalive: false}
	rtmgr.Eps = EP

	var rmrpush = RmrPush{}
	rmrpush.AddEndpoint(rtmgr.Eps["127.0.0.2"])

	var E2map = make(map[string]rtmgr.E2TInstance)

	E2map["127.0.0.2:4562"] = rtmgr.E2TInstance{
		Name:    "E2Tinstance1",
		Fqdn:    "127.0.0.2:4562",
		Ranlist: []string{"1", "2"},
	}

	sbi.updateE2TEndpoints(&E2map, sbii)
	t.Log(err)
}

func TestPruneEndpointList(t *testing.T) {
	var sbi = Sbi{}
	var err error
	sbii, err := GetSbi("rmrpush")
	var EP = make(map[string]*rtmgr.Endpoint)
	EP["127.0.0.2"] = &rtmgr.Endpoint{Uuid: "127.0.0.2", Name: "E2TERM", XAppType: "app1", Ip: "127.0.0.1", Port: 4562, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: true, Keepalive: false}
	rtmgr.Eps = EP

	var rmrpush = RmrPush{}
	rmrpush.AddEndpoint(rtmgr.Eps["127.0.0.2"])

	sbi.pruneEndpointList(sbii)
	t.Log(err)
}
