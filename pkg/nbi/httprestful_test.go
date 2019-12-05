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
  Mnemonic:     httprestful_test.go
  Abstract:     HTTPRestful unit tests
  Date:         15 May 2019
*/

package nbi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"routing-manager/pkg/models"
	"routing-manager/pkg/sdl"
	"routing-manager/pkg/stub"
	"testing"
	"time"

	"github.com/go-openapi/swag"
)

var BasicXAppLists = []byte(`[
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

var SubscriptionResp = []byte(`{"ID":"deadbeef1234567890", "Version":0, "EventType":"all"}`)

var InvalidSubResp = []byte(`{"Version":0, "EventType":all}`)

func createMockAppmgrWithData(url string, g []byte, p []byte) *httptest.Server {
	l, err := net.Listen("tcp", url)
	if err != nil {
		fmt.Println("Failed to create listener: " + err.Error())
	}
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.String() == "/ric/v1/xapps" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(g)
		}
		if r.Method == "POST" && r.URL.String() == "/ric/v1/subscriptions" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(p)
		}

	}))
	ts.Listener.Close()
	ts.Listener = l
	return ts
}

func createMockPlatformComponents() {
	var filename = "config.json"
	file, _ := json.MarshalIndent(stub.ValidPlatformComponents, "", "")
	filestr := string(file)
	filestr = "{\"PlatformComponents\":" + filestr + "}"
	file = []byte(filestr)
	_ = ioutil.WriteFile(filename, file, 644)
}

func TestRecvXappCallbackData(t *testing.T) {
	data := models.XappCallbackData{
		XApps:   *swag.String("[]"),
		Version: *swag.Int64(1),
		Event:   *swag.String("any"),
		ID:      *swag.String("123456"),
	}

	ch := make(chan *models.XappCallbackData)
	defer close(ch)
	httpRestful := NewHttpRestful()
	go func() { ch <- &data }()
	time.Sleep(1 * time.Second)
	t.Log(string(len(ch)))
	xappList, err := httpRestful.RecvXappCallbackData(ch)
	if err != nil {
		t.Error("Receive failed: " + err.Error())
	} else {
		if xappList == nil {
			t.Error("Expected an XApp notification list")
		} else {
			t.Log("whatever")
		}
	}
}

func TestProvideXappHandleHandlerImpl(t *testing.T) {
	datach := make(chan *models.XappCallbackData, 10)
	defer close(datach)
	data := models.XappCallbackData{
		XApps:   *swag.String("[]"),
		Version: *swag.Int64(1),
		Event:   *swag.String("someevent"),
		ID:      *swag.String("123456")}
	var httpRestful, _ = GetNbi("httpRESTful")
	err := httpRestful.(*HttpRestful).ProvideXappHandleHandlerImpl(datach, &data)
	if err != nil {
		t.Error("Error occured: " + err.Error())
	} else {
		recv := <-datach
		if recv == nil {
			t.Error("Something gone wrong: " + err.Error())
		} else {
			if recv != &data {
				t.Error("Malformed data on channel")
			}
		}
	}
}

func TestValidateXappCallbackData(t *testing.T) {
	data := models.XappCallbackData{
		XApps:   *swag.String("[]"),
		Version: *swag.Int64(1),
		Event:   *swag.String("someevent"),
		ID:      *swag.String("123456")}

	err := validateXappCallbackData(&data)
	if err != nil {
		t.Error("Invalid XApp callback data: " + err.Error())
	}
}

func TestValidateXappCallbackDataWithInvalidData(t *testing.T) {
	data := models.XappCallbackData{
		XApps:   *swag.String("{}"),
		Version: *swag.Int64(1),
		Event:   *swag.String("someevent"),
		ID:      *swag.String("123456")}

	err := validateXappCallbackData(&data)
	if err == nil {
		t.Error("Invalid XApp callback data: " + err.Error())
	}
}

func TestHttpGetXAppsInvalidData(t *testing.T) {
	_, err := httpGetXApps(XMURL)
	if err == nil {
		t.Error("No XApp data received: " + err.Error())
	}
}

func TestHttpGetXAppsWithValidData(t *testing.T) {
	var expected = 1
	ts := createMockAppmgrWithData("127.0.0.1:3000", BasicXAppLists, nil)

	ts.Start()
	defer ts.Close()
	xapplist, err := httpGetXApps(XMURL)
	if err != nil {
		t.Error("Error occured: " + err.Error())
	} else {
		if len(*xapplist) != expected {
			t.Error("Invalid XApp data: got " + string(len(*xapplist)) + ", expected " + string(expected))
		}
	}
}

func TestRetrieveStartupDataTimeout(t *testing.T) {
	sdlEngine, _ := sdl.GetSdl("file")
	createMockPlatformComponents()
	err := retrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", sdlEngine)
	if err == nil {
		t.Error("Cannot retrieve startup data: " + err.Error())
	}
	os.Remove("rt.json")
	os.Remove("config.json")
}

func TestRetrieveStartupData(t *testing.T) {
	ts := createMockAppmgrWithData("127.0.0.1:3000", BasicXAppLists, SubscriptionResp)
	ts.Start()
	defer ts.Close()
	sdlEngine, _ := sdl.GetSdl("file")
	var httpRestful, _ = GetNbi("httpRESTful")
	createMockPlatformComponents()
	err := httpRestful.(*HttpRestful).RetrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", sdlEngine)
	//err := retrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", sdlEngine)
	if err != nil {
		t.Error("Cannot retrieve startup data: " + err.Error())
	}
	os.Remove("rt.json")
	os.Remove("config.json")
}

func TestRetrieveStartupDataWithInvalidSubResp(t *testing.T) {
	ts := createMockAppmgrWithData("127.0.0.1:3000", BasicXAppLists, InvalidSubResp)
	ts.Start()
	defer ts.Close()
	sdlEngine, _ := sdl.GetSdl("file")
	var httpRestful, _ = GetNbi("httpRESTful")
	createMockPlatformComponents()
	err := httpRestful.(*HttpRestful).RetrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", sdlEngine)
	if err == nil {
		t.Error("Cannot retrieve startup data: " + err.Error())
	}
	os.Remove("rt.json")
	os.Remove("config.json")
}