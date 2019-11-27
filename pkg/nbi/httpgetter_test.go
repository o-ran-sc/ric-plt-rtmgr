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
  Mnemonic:     httpgetter.go
  Abstract:     HTTPgetter unit tests
  Date:         14 May 2019
*/

package nbi

import (
	"net"
	"net/http"
	"net/http/httptest"
	"routing-manager/pkg/rtmgr"
	"testing"
)

var (
	XMURL = "http://127.0.0.1:3000/ric/v1/xapps"
)

func TestFetchXappListInvalidData(t *testing.T) {
	var httpGetter = NewHttpGetter()
	_, err := httpGetter.FetchAllXApps(XMURL)
	if err == nil {
		t.Error("No XApp data received: " + err.Error())
	}
}

func TestFetchXappListWithInvalidData(t *testing.T) {
	var expected = 0
	rtmgr.SetLogLevel("debug")
	b := []byte(`{"ID":"deadbeef1234567890", "Version":0, "EventType":"all"}`)
	l, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		t.Error("Failed to create listener: " + err.Error())
	}
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//t.Log(r.Method)
		//t.Log(r.URL)
		if r.Method == "GET" && r.URL.String() == "/ric/v1/xapps" {
			//t.Log("Sending reply")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	}))
	ts.Listener.Close()
	ts.Listener = l

	ts.Start()
	defer ts.Close()
	var httpGetter = NewHttpGetter()
	xapplist, err := httpGetter.FetchAllXApps(XMURL)
	if err == nil {
		t.Error("Error occured: " + err.Error())
	} else {
		//t.Log(len(*xapplist))
		if len(*xapplist) != expected {
			t.Error("Invalid XApp data: got " + string(len(*xapplist)) + ", expected " + string(expected))
		}
	}
}

func TestFetchAllXAppsWithValidData(t *testing.T) {
	var expected = 1
	b := []byte(`[
 {
 "name":"xapp-01","status":"unknown","version":"1.2.3",
    "instances":[
        {"name":"xapp-01-instance-01","status":"pending","ip":"172.16.1.103","port":4555,
            "txMessages":["ControlIndication"],
            "rxMessages":["LoadIndication","Reset"]
        },
        {"name":"xapp-01-instance-02","status":"pending","ip":"10.244.1.12","port":4561,
            "txMessages":["ControlIndication","SNStatusTransfer"],
            "rxMessages":["LoadIndication","HandoverPreparation"]
        }
    ]
}
]`)
	l, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		t.Error("Failed to create listener: " + err.Error())
	}
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//t.Log(r.Method)
		//t.Log(r.URL)
		if r.Method == "GET" && r.URL.String() == "/ric/v1/xapps" {
			//t.Log("Sending reply")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	}))
	ts.Listener.Close()
	ts.Listener = l

	ts.Start()
	defer ts.Close()
	var httpGetter = NewHttpGetter()
	xapplist, err := httpGetter.FetchAllXApps(XMURL)
	if err != nil {
		t.Error("Error occured: " + err.Error())
	} else {
		if len(*xapplist) != expected {
			t.Error("Invalid XApp data: got " + string(len(*xapplist)) + ", expected " + string(expected))
		}
	}
}
