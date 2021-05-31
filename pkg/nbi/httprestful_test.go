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
	xfmodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/go-openapi/swag"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"routing-manager/pkg/models"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sbi"
	"routing-manager/pkg/sdl"
	"routing-manager/pkg/stub"
	"sync"
	"testing"
	"time"
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

var E2TListResp = []byte(`[{"e2tAddress":"127.0.0.1:0","ranNames":["RanM0","RanN0"]},{"e2tAddress":"127.0.0.1:1","ranNames":["RanM1","RanN1"]},{"e2tAddress":"127.0.0.1:2","ranNames":["RanM2","RanN2"]},{"e2tAddress":"127.0.0.1:3","ranNames":["RanM3","RanN3"]}]`)

var SubscriptionList = []byte(`[{"SubscriptionId":11,"Meid":"Test-Gnb","Endpoint":["127.0.0.1:4056"]}]`)

var InvalidSubResp = []byte(`{"Version":0, "EventType":all}`)

type Consumer struct{}

func (m Consumer) Consume(params *xapp.RMRParams) (err error) {
	xapp.Sdl.Store("myKey", params.Payload)
	return nil
}

// Test cases
func TestMain(m *testing.M) {
	go xapp.RunWithParams(Consumer{}, viper.GetBool("db.waitForSdl"))
	time.Sleep(time.Duration(5) * time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestValidateXappCallbackData_1(t *testing.T) {
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

func TestValidateXappCallbackDataInvalid(t *testing.T) {
	data := models.XappCallbackData{}
	err := validateXappCallbackData(&data)
	t.Log(err)
}

func TestValidateXappSubscriptionsData(t *testing.T) {

	rtmgr.RMRConnStatus = make(map[string]bool)
	ep := make(map[string]*rtmgr.Endpoint)
	ep["dummy"] = &rtmgr.Endpoint{Uuid: "10.0.0.1:0", Name: "E2TERM", XAppType: "app1", Ip: "", Port: 0, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: true, Keepalive: true}
	p := uint16(1234)
	data := models.XappSubscriptionData{
		Address:        swag.String("10.1.1.1"),
		Port:           &p,
		SubscriptionID: swag.Int32(123456)}

	var err error
	err = validateXappSubscriptionData(&data)
	t.Log(err)

	rtmgr.Eps = ep
	p = uint16(0)
	data1 := models.XappSubscriptionData{
		Address:        swag.String(""),
		Port:           &p,
		SubscriptionID: swag.Int32(123456)}
	err = validateXappSubscriptionData(&data1)

	//Validate E2tData
	data2 := models.E2tData{
		E2TAddress: swag.String(""),
	}
	/*err = validateE2tData(&data2)*/

	//e2tchannel := make(chan *models.E2tData, 10)
	_ = CreateNewE2tHandleHandlerImpl(&data2)
	//defer close(e2tchannel)

	//test case for provideXappSubscriptionHandleImp
	//datachannel := make(chan *models.XappSubscriptionData, 10)
	sdlEngine, _ = sdl.GetSdl("file")
	sbiEngine, _ = sbi.GetSbi("rmrpush")
	rpeEngine, _ = rpe.GetRpe("rmrpush")
	_ = ProvideXappSubscriptionHandleImpl(&data1)
	//defer close(datachannel)

	//test case for deleteXappSubscriptionHandleImpl
	_ = DeleteXappSubscriptionHandleImpl(&data1)

	data3 := models.XappSubscriptionData{
		Address:        swag.String("10.55.55.5"),
		Port:           &p,
		SubscriptionID: swag.Int32(123456)}
	//test case for deleteXappSubscriptionHandleImpl
	_ = DeleteXappSubscriptionHandleImpl(&data3)
	data4 := models.XappSubscriptionData{
		Address:        swag.String("1.5.5.5"),
		Port:           &p,
		SubscriptionID: swag.Int32(1236)}
	_ = DeleteXappSubscriptionHandleImpl(&data4)

}

func TestValidateE2tDataEmpty(t *testing.T) {
	data := models.E2tData{
		E2TAddress: swag.String(""),
	}
	err, _ := validateE2tData(&data)
	t.Log(err)
}

func TestValidateE2tDataDNSLookUPfails(t *testing.T) {
	rtmgr.RMRConnStatus = make(map[string]bool)
	data := models.E2tData{
		E2TAddress: swag.String("e2t.1com:1234"),
	}
	err, _ := validateE2tData(&data)
	t.Log(err)
}

func TestValidateE2tDataInvalid(t *testing.T) {
	rtmgr.RMRConnStatus = make(map[string]bool)
	data := models.E2tData{
		E2TAddress: swag.String("10.101.01.1"),
	}
	err, _ := validateE2tData(&data)
	t.Log(err)
}

func TestValidateE2tDatavalid(t *testing.T) {
	rtmgr.RMRConnStatus = make(map[string]bool)
	data := models.E2tData{
		E2TAddress: swag.String("10.101.01.1:8098"),
	}

	err, _ := validateE2tData(&data)
	t.Log(err)

	_ = CreateNewE2tHandleHandlerImpl(&data)

}

func TestValidateE2tDatavalidEndpointPresent(t *testing.T) {
	rtmgr.RMRConnStatus = make(map[string]bool)
	data := models.E2tData{
		E2TAddress: swag.String("10.101.01.1:8098"),
	}

	// insert endpoint for testing purpose
	uuid := "10.101.01.1:8098"
	ep := &rtmgr.Endpoint{
		Uuid: uuid,
	}
	rtmgr.Eps[uuid] = ep

	err, _ := validateE2tData(&data)
	t.Log(err)

	// delete endpoint for at end of test case
	delete(rtmgr.Eps, uuid)

}

func TestValidateDeleteE2tData(t *testing.T) {

	rtmgr.RMRConnStatus = make(map[string]bool)
	// test-1
	data := models.E2tDeleteData{
		E2TAddress: swag.String(""),
	}

	err := validateDeleteE2tData(&data)
	if err.Error() != "E2TAddress is empty!!!" {
		t.Log(err)
	}

	// test-2
	data = models.E2tDeleteData{
		E2TAddress: swag.String("10.101.01.1:8098"),
	}

	err = validateDeleteE2tData(&data)
	if err != nil {
		t.Log(err)
	}

	// test-3
	//################ Create End Point dummy entry
	uuid := "10.101.01.1:8098"
	ep := &rtmgr.Endpoint{
		Uuid: uuid,
	}
	rtmgr.Eps[uuid] = ep
	//#####################

	data = models.E2tDeleteData{
		E2TAddress: swag.String("10.101.01.1:8098"),
		RanAssocList: models.RanE2tMap{
			{E2TAddress: swag.String("10.101.01.1:8098")},
		},
	}

	err = validateDeleteE2tData(&data)
	if err != nil {
		t.Log(err)
	}

	// delete endpoint for at end of test case
	//################ Delete End Point dummy entry
	delete(rtmgr.Eps, uuid)
	//#####################

	// test-4

	//################ Create End Point dummy entry
	uuid = "10.101.01.1:9991"
	ep = &rtmgr.Endpoint{
		Uuid: uuid,
	}
	rtmgr.Eps[uuid] = ep

	uuid = "10.101.01.1:9992"
	ep = &rtmgr.Endpoint{
		Uuid: uuid,
	}
	rtmgr.Eps[uuid] = ep
	//#####################

	data = models.E2tDeleteData{
		E2TAddress: swag.String("10.101.01:8098"),
		RanAssocList: models.RanE2tMap{
			{E2TAddress: swag.String("10.101.01.1:9991")},
			{E2TAddress: swag.String("10.101.01.1:9992")},
		},
	}

	err = validateDeleteE2tData(&data)
	if err != nil {
		t.Log(err)
	}
	//################ Delete End Point dummy entry
	delete(rtmgr.Eps, "10.101.01.1:9991")
	delete(rtmgr.Eps, "10.101.01.1:9992")
	//#####################

	// test-5

	data = models.E2tDeleteData{
		E2TAddress: swag.String("10.101.01:8098"),
		RanAssocList: models.RanE2tMap{
			{E2TAddress: swag.String("10.101.01.19991")},
		},
	}

	err = validateDeleteE2tData(&data)
	if err.Error() != "E2T Delete - RanAssocList E2TAddress is not a proper format like ip:port, 10.101.01.19991" {
		t.Log(err)
	}
}

func TestValidateE2TAddressRANListData(t *testing.T) {

	data := models.RanE2tMap{
		{
			E2TAddress: swag.String(""),
		},
	}
	err := validateE2TAddressRANListData(data)
	if err != nil {
		t.Log(err)
	}

	data = models.RanE2tMap{
		{
			E2TAddress: swag.String("10.101.01.1:8098"),
		},
	}
	err = validateE2TAddressRANListData(data)
	if err != nil {
		t.Log(err)
	}

}

func TestAssociateRanToE2THandlerImpl(t *testing.T) {

	data := models.RanE2tMap{
		{
			E2TAddress: swag.String("10.101.01.1:8098"),
		},
	}
	err := AssociateRanToE2THandlerImpl(data)
	if err != nil {
		t.Log(err)
	}

	//################ Create End Point dummy entry
	uuid := "10.101.01.1:8098"
	ep := &rtmgr.Endpoint{
		Uuid: uuid,
	}
	rtmgr.Eps[uuid] = ep
	//#####################

	data = models.RanE2tMap{
		{
			E2TAddress: swag.String("10.101.01.1:8098"),
		},
	}
	err = AssociateRanToE2THandlerImpl(data)
	if err != nil {
		t.Log(err)
	}

	//################ Delete End Point dummy entry
	delete(rtmgr.Eps, uuid)
	//#####################
}

func TestDisassociateRanToE2THandlerImpl(t *testing.T) {

	data := models.RanE2tMap{
		{
			E2TAddress: swag.String("10.101.01.1:8098"),
		},
	}
	err := DisassociateRanToE2THandlerImpl(data)
	if err != nil {
		t.Log(err)
	}
	//################ Create End Point dummy entry
	uuid := "10.101.01.1:8098"
	ep := &rtmgr.Endpoint{
		Uuid: uuid,
	}
	rtmgr.Eps[uuid] = ep
	//#####################

	data = models.RanE2tMap{
		{
			E2TAddress: swag.String("10.101.01.1:8098"),
		},
	}
	err = DisassociateRanToE2THandlerImpl(data)
	if err != nil {
		t.Log(err)
	}

	//################ Delete End Point dummy entry
	delete(rtmgr.Eps, uuid)
	//#####################
}

func TestDeleteE2tHandleHandlerImpl(t *testing.T) {

	data := models.E2tDeleteData{
		E2TAddress: swag.String(""),
	}
	err := DeleteE2tHandleHandlerImpl(&data)
	if err != nil {
		t.Log(err)
	}

	//################ Create End Point dummy entry
	uuid := "10.101.01.1:8098"
	ep := &rtmgr.Endpoint{
		Uuid: uuid,
	}
	rtmgr.Eps[uuid] = ep
	//#####################

	data = models.E2tDeleteData{
		E2TAddress: swag.String("10.101.01.1:8098"),
	}
	err = DeleteE2tHandleHandlerImpl(&data)
	if err != nil {
		t.Log(err)
	}
	//################ Delete End Point dummy entry
	delete(rtmgr.Eps, uuid)
	//#####################
}

func TestSubscriptionExists(t *testing.T) {
	p := uint16(0)
	data := models.XappSubscriptionData{
		Address:        swag.String("10.0.0.0"),
		Port:           &p,
		SubscriptionID: swag.Int32(1234)}

	rtmgr.Subs = *stub.ValidSubscriptions

	yes_no := subscriptionExists(&data)
	yes_no = addSubscription(&rtmgr.Subs, &data)
	yes_no = addSubscription(&rtmgr.Subs, &data)
	yes_no = delSubscription(&rtmgr.Subs, &data)
	yes_no = delSubscription(&rtmgr.Subs, &data)
	t.Log(yes_no)
}

func TestAddSubscriptions(t *testing.T) {
	p := uint16(1)
	subdata := models.XappSubscriptionData{
		Address:        swag.String("10.0.0.0"),
		Port:           &p,
		SubscriptionID: swag.Int32(1234)}

	rtmgr.Subs = *stub.ValidSubscriptions
	yes_no := addSubscription(&rtmgr.Subs, &subdata)
	t.Log(yes_no)
}

func TestHttpInstance(t *testing.T) {
	sdlEngine, _ := sdl.GetSdl("file")
	rpeEngine, _ := rpe.GetRpe("rmrpush")
	sbiEngine, _ := sbi.GetSbi("rmrpush")
	httpinstance := NewHttpRestful()
	err := httpinstance.Terminate()
	t.Log(err)
	fmt.Printf("sbiEngine = %v", sbiEngine)

	createMockPlatformComponents()
	//ts := createMockAppmgrWithData("127.0.0.1:3000", BasicXAppLists, nil)
	//ts.Start()
	//defer ts.Close()
	var m sync.Mutex
	err = httpinstance.Initialize(XMURL, "httpgetter", "rt.json", "config.json", E2MURL, sdlEngine, rpeEngine, &m)
}

func TestXappCallbackWithData(t *testing.T) {
	data := models.XappCallbackData{
		XApps:   *swag.String("[]"),
		Version: *swag.Int64(1),
		Event:   *swag.String("someevent"),
		ID:      *swag.String("123456")}
	_, _ = recvXappCallbackData(&data)
}

func TestXappCallbackNodata(t *testing.T) {
	//data := *models.XappCallbackData
	_, _ = recvXappCallbackData(nil)
}

func TestE2TwithData(t *testing.T) {
	data2 := models.E2tData{
		E2TAddress:  swag.String("1.2.3.4"),
		RanNamelist: []string{"ran1", "ran2"},
	}
	_, _, _ = recvNewE2Tdata(&data2)
}

func TestE2TwithNoData(t *testing.T) {
	_, _, _ = recvNewE2Tdata(nil)
}

func TestProvideXappSubscriptionHandleImpl(t *testing.T) {
	p := uint16(0)
	data := models.XappSubscriptionData{
		Address:        swag.String("10.0.0.0"),
		Port:           &p,
		SubscriptionID: swag.Int32(1234)}
	_ = ProvideXappSubscriptionHandleImpl(&data)
}

func createMockAppmgrWithData(url string, g []byte, p []byte, t []byte) *httptest.Server {
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
		if r.Method == "GET" && r.URL.String() == "/ric/v1/e2t/list" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(t)
		}

	}))
	ts.Listener.Close()
	ts.Listener = l
	return ts
}

func createMockSubmgrWithData(url string, t []byte) *httptest.Server {
	l, err := net.Listen("tcp", url)
	if err != nil {
		fmt.Println("Failed to create listener: " + err.Error())
	}
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" && r.URL.String() == "//ric/v1/subscriptions" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(t)
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

func TestProvideXappHandleHandlerImpl(t *testing.T) {
	data := models.XappCallbackData{
		XApps:   *swag.String("[]"),
		Version: *swag.Int64(1),
		Event:   *swag.String("someevent"),
		ID:      *swag.String("123456")}
	err := ProvideXappHandleHandlerImpl(&data)

	//Empty XappCallbackdata
	data1 := models.XappCallbackData{}
	err = ProvideXappHandleHandlerImpl(&data1)
	t.Log(err)
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
	ts := createMockAppmgrWithData("127.0.0.1:3000", BasicXAppLists, nil, nil)

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
	err := retrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", E2MURL, sdlEngine)
	if err == nil {
		t.Error("Cannot retrieve startup data: " + err.Error())
	}
	os.Remove("rt.json")
	os.Remove("config.json")
}

func TestRetrieveStartupData(t *testing.T) {
	ts := createMockAppmgrWithData("127.0.0.1:3000", BasicXAppLists, SubscriptionResp, nil)
	ts.Start()
	defer ts.Close()

	ts1 := createMockAppmgrWithData("127.0.0.1:8085", nil, nil, E2TListResp)
	ts1.Start()
	defer ts1.Close()

	ts2 := createMockSubmgrWithData("127.0.0.1:8089", SubscriptionList)
	ts2.Start()
	defer ts2.Close()

	sdlEngine, _ := sdl.GetSdl("file")
	var httpRestful, _ = GetNbi("httpRESTful")
	createMockPlatformComponents()
	httpRestful.(*HttpRestful).RetrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", E2MURL, sdlEngine)
	//err := retrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", sdlEngine)
	/*if err != nil {
		t.Error("Cannot retrieve startup data: " + err.Error())
	}*/
	os.Remove("rt.json")
	os.Remove("config.json")
}

func TestRetrieveStartupDataWithInvalidSubResp(t *testing.T) {
	ts := createMockAppmgrWithData("127.0.0.1:3000", BasicXAppLists, InvalidSubResp, nil)
	ts.Start()
	defer ts.Close()
	sdlEngine, _ := sdl.GetSdl("file")
	var httpRestful, _ = GetNbi("httpRESTful")
	createMockPlatformComponents()
	err := httpRestful.(*HttpRestful).RetrieveStartupData(XMURL, "httpgetter", "rt.json", "config.json", E2MURL, sdlEngine)
	if err == nil {
		t.Error("Cannot retrieve startup data: " + err.Error())
	}
	os.Remove("rt.json")
	os.Remove("config.json")
}

func TestInvalidarguments(t *testing.T) {
	_ = PostSubReq("\n", "nbifinterface")
	_ = PostSubReq("xmurl", "\n")
}

func TestInitEngine(t *testing.T) {
	initRtmgr()
}

func TestUpdateXappSubscription(t *testing.T) {
	ep := make(map[string]*rtmgr.Endpoint)
	ep["dummy"] = &rtmgr.Endpoint{Uuid: "10.0.0.1:0", Name: "E2TERM", XAppType: "app1", Ip: "10.1.1.1", Port: 1234, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: true, Keepalive: true}

	rtmgr.Eps = ep

	p := uint16(1234)
	xapp := models.XappElement{
		Address: swag.String("10.1.1.1"),
		Port:    &p,
	}

	var b models.XappList
	b = append(b, &xapp)
	_ = UpdateXappSubscriptionHandleImpl(&b, 10)

	//Test case when subscriptions already exist
	data := models.XappSubscriptionData{
		Address:        swag.String("10.0.0.0"),
		Port:           &p,
		SubscriptionID: swag.Int32(12345)}

	rtmgr.Subs = *stub.ValidSubscriptions

	subscriptionExists(&data)
	addSubscription(&rtmgr.Subs, &data)
	_ = UpdateXappSubscriptionHandleImpl(&b, 10)

}

func TestDumpDebugdata(t *testing.T) {
	_, _ = DumpDebugData()
}

func TestManagerRequest(t *testing.T) {
	var params xapp.RMRParams
	var rmrmeid xapp.RMRMeid
	sdlEngine, _ = sdl.GetSdl("file")
	sbiEngine, _ = sbi.GetSbi("rmrpush")
	rpeEngine, _ = rpe.GetRpe("rmrpush")
	rmrmeid.RanName = "gnb1"
	c := Control{make(chan *xapp.RMRParams)}
	params.Payload = []byte{1, 2, 3, 4}
	params.Mtype = 1234
	params.SubId = -1
	params.Meid = &rmrmeid
	params.Src = "sender"
	params.PayloadLen = 4
	c.handleUpdateToRoutingManagerRequest(&params)
}

func TestRecievermr(t *testing.T) {
	var params xapp.RMRParams
	var rmrmeid xapp.RMRMeid
	sdlEngine, _ = sdl.GetSdl("file")
	sbiEngine, _ = sbi.GetSbi("rmrpush")
	rpeEngine, _ = rpe.GetRpe("rmrpush")
	rmrmeid.RanName = "gnb1"

	params.Payload = []byte{1, 2, 3, 4}
	params.SubId = -1
	params.Meid = &rmrmeid
	params.Src = "sender"
	params.PayloadLen = 4

	c := Control{make(chan *xapp.RMRParams)}
	params.Mtype = xapp.RICMessageTypes["RMRRM_REQ_TABLE"]
	c.recievermr(&params)
	params.Mtype = xapp.RICMessageTypes["RMRRM_TABLE_STATE"]
	c.recievermr(&params)
	params.Mtype = 1234
	c.recievermr(&params)

	rtmgr.Rtmgr_ready = true
	params.Mtype = xapp.RICMessageTypes["RMRRM_REQ_TABLE"]
	c.recievermr(&params)
	time.Sleep(time.Duration(5) * time.Second)
}

func TestAddDelRmr(t *testing.T) {
	sdlEngine, _ = sdl.GetSdl("file")
	sbiEngine, _ = sbi.GetSbi("rmrpush")
	rpeEngine, _ = rpe.GetRpe("rmrpush")
	var routelist models.Routelist
	mtype := uint32(1234)
	tendpoint := "goofle.com"
	listofroutes := models.AddRmrRoute{SubscriptionID: 0, SenderEndPoint: "nokia.com", MessageType: &mtype, TargetEndPoint: &tendpoint}
	listofroutes2 := models.AddRmrRoute{SubscriptionID: 1, SenderEndPoint: "", MessageType: &mtype, TargetEndPoint: &tendpoint}
	listofroutes3 := models.AddRmrRoute{MessageType: &mtype, TargetEndPoint: &tendpoint}
	Adddelrmrroute(routelist, false)
	routelist = append(routelist, &listofroutes)
	routelist = append(routelist, &listofroutes2)
	routelist = append(routelist, &listofroutes3)
	routelist = append(routelist, &listofroutes3)
	Adddelrmrroute(routelist, true)

	Adddelrmrroute(routelist, false)
}

func TestPopulateSubscription(t *testing.T) {
	var sublist xfmodel.SubscriptionList

	subdata := xfmodel.SubscriptionData{ClientEndpoint: []string{"xapp1.com:3800"}, SubscriptionID: -1, Meid: ""}
	subdata2 := xfmodel.SubscriptionData{ClientEndpoint: []string{"xapp2.com:3800"}, SubscriptionID: 11, Meid: ""}
	subdata3 := xfmodel.SubscriptionData{ClientEndpoint: []string{"xapp3.com:3800"}, SubscriptionID: 221, Meid: ""}
	sublist = append(sublist, &subdata)
	sublist = append(sublist, &subdata2)
	sublist = append(sublist, &subdata3)

	PopulateSubscription(sublist)

}
