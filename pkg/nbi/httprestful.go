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
  Mnemonic:	httprestful.go
  Abstract:	HTTP Restful API NBI implementation
                Based on Swagger generated code
  Date:		25 March 2019
*/

package nbi

//noinspection GoUnresolvedReference,GoUnresolvedReference,GoUnresolvedReference,GoUnresolvedReference,GoUnresolvedReference,GoUnresolvedReference
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"net/url"
	"os"
	"routing-manager/pkg/models"
	"routing-manager/pkg/restapi"
	"routing-manager/pkg/restapi/operations"
	"routing-manager/pkg/restapi/operations/handle"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sdl"
	"strconv"
	"time"
)

//var myClient = &http.Client{Timeout: 1 * time.Second}

type HttpRestful struct {
	Engine
	LaunchRest                   LaunchRestHandler
	RecvXappCallbackData         RecvXappCallbackDataHandler
	ProvideXappHandleHandlerImpl ProvideXappHandleHandlerImpl
	RetrieveStartupData          RetrieveStartupDataHandler
}

func NewHttpRestful() *HttpRestful {
	instance := new(HttpRestful)
	instance.LaunchRest = launchRest
	instance.RecvXappCallbackData = recvXappCallbackData
	instance.ProvideXappHandleHandlerImpl = provideXappHandleHandlerImpl
	instance.RetrieveStartupData = retrieveStartupData
	return instance
}

// ToDo: Use Range over channel. Read and return only the latest one.
func recvXappCallbackData(dataChannel <-chan *models.XappCallbackData) (*[]rtmgr.XApp, error) {
	var xappData *models.XappCallbackData
	// Drain the channel as we are only looking for the latest value until
	// xapp manager sends all xapp data with every request.
	length := len(dataChannel)
	//rtmgr.Logger.Info(length)
	for i := 0; i <= length; i++ {
		rtmgr.Logger.Info("data received")
		// If no data received from the REST, it blocks.
		xappData = <-dataChannel
	}
	if nil != xappData {
		var xapps []rtmgr.XApp
		err := json.Unmarshal([]byte(xappData.XApps), &xapps)
		return &xapps, err
	} else {
		rtmgr.Logger.Info("No data")
	}

	rtmgr.Logger.Debug("Nothing received on the Http interface")
	return nil, nil
}

func validateXappCallbackData(callbackData *models.XappCallbackData) error {
	if len(callbackData.XApps) == 0 {
		return fmt.Errorf("invalid Data field: \"%s\"", callbackData.XApps)
	}
	var xapps []rtmgr.XApp
	err := json.Unmarshal([]byte(callbackData.XApps), &xapps)
	if err != nil {
		return fmt.Errorf("unmarshal failed: \"%s\"", err.Error())
	}
	return nil
}

func provideXappHandleHandlerImpl(datach chan<- *models.XappCallbackData, data *models.XappCallbackData) error {
	if data != nil {
		rtmgr.Logger.Debug("Received callback data")
	}
	err := validateXappCallbackData(data)
	if err != nil {
		rtmgr.Logger.Warn("XApp callback data validation failed: " + err.Error())
		return err
	} else {
		datach <- data
		return nil
	}
}

func validateXappSubscriptionData(data *models.XappSubscriptionData) error {
	var err = fmt.Errorf("XApp instance not found: %v:%v", *data.Address, *data.Port)
	for _, ep := range rtmgr.Eps {
		if ep.Ip == *data.Address && ep.Port == *data.Port {
			err = nil
			break
		}
	}
	return err
}

func provideXappSubscriptionHandleImpl(subchan chan<- *models.XappSubscriptionData,
	data *models.XappSubscriptionData) error {
	rtmgr.Logger.Debug("Invoked provideXappSubscriptionHandleImpl")
	err := validateXappSubscriptionData(data)
	if err != nil {
		rtmgr.Logger.Error(err.Error())
		return err
	}
	subchan <- data
	//var val = string(*data.Address + ":" + strconv.Itoa(int(*data.Port)))
	rtmgr.Logger.Debug("Endpoints: %v", rtmgr.Eps)
	return nil
}

func subscriptionExists(data *models.XappSubscriptionData) bool {
	present := false
	sub := rtmgr.Subscription{SubID: *data.SubscriptionID, Fqdn: *data.Address, Port: *data.Port}
	for _, elem := range rtmgr.Subs {
		if elem == sub {
			present = true
			break
		}
	}
	return present
}

func deleteXappSubscriptionHandleImpl(subdelchan chan<- *models.XappSubscriptionData,
	data *models.XappSubscriptionData) error {
	rtmgr.Logger.Debug("Invoked deleteXappSubscriptionHandleImpl")
	err := validateXappSubscriptionData(data)
	if err != nil {
		rtmgr.Logger.Error(err.Error())
		return err
	}

	if !subscriptionExists(data) {
		rtmgr.Logger.Warn("subscription not found: %d", *data.SubscriptionID)
		err := fmt.Errorf("subscription not found: %d", *data.SubscriptionID)
		return err
	}

	subdelchan <- data
	return nil
}

func launchRest(nbiif *string, datach chan<- *models.XappCallbackData, subchan chan<- *models.XappSubscriptionData,
	subdelchan chan<- *models.XappSubscriptionData) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		//log.Fatalln(err)
		rtmgr.Logger.Error(err.Error())
		os.Exit(1)
	}
	nbiUrl, err := url.Parse(*nbiif)
	if err != nil {
		rtmgr.Logger.Error(err.Error())
		os.Exit(1)
	}
	api := operations.NewRoutingManagerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port, err = strconv.Atoi(nbiUrl.Port())
	if err != nil {
		rtmgr.Logger.Error("Invalid NBI RestAPI port")
		os.Exit(1)
	}
	server.Host = "0.0.0.0"
	// set handlers
	api.HandleProvideXappHandleHandler = handle.ProvideXappHandleHandlerFunc(
		func(params handle.ProvideXappHandleParams) middleware.Responder {
			rtmgr.Logger.Info("Data received on Http interface")
			err := provideXappHandleHandlerImpl(datach, params.XappCallbackData)
			if err != nil {
				rtmgr.Logger.Error("Invalid XApp callback data: " + err.Error())
				return handle.NewProvideXappHandleBadRequest()
			} else {
				return handle.NewGetHandlesOK()
			}
		})
	api.HandleProvideXappSubscriptionHandleHandler = handle.ProvideXappSubscriptionHandleHandlerFunc(
		func(params handle.ProvideXappSubscriptionHandleParams) middleware.Responder {
			err := provideXappSubscriptionHandleImpl(subchan, params.XappSubscriptionData)
			if err != nil {
				return handle.NewProvideXappSubscriptionHandleBadRequest()
			} else {
				//Delay the reponse as add subscription channel needs to update sdl and then sbi sends updated routes to all endpoints
				time.Sleep(1 * time.Second)
				return handle.NewGetHandlesOK()
			}
		})
	api.HandleDeleteXappSubscriptionHandleHandler = handle.DeleteXappSubscriptionHandleHandlerFunc(
		func(params handle.DeleteXappSubscriptionHandleParams) middleware.Responder {
			err := deleteXappSubscriptionHandleImpl(subdelchan, params.XappSubscriptionData)
			if err != nil {
				return handle.NewDeleteXappSubscriptionHandleNoContent()
			} else {
				//Delay the reponse as delete subscription channel needs to update sdl and then sbi sends updated routes to all endpoints
				time.Sleep(1 * time.Second)
				return handle.NewGetHandlesOK()
			}
		})
	// start to serve API
	rtmgr.Logger.Info("Starting the HTTP Rest service")
	if err := server.Serve(); err != nil {
		rtmgr.Logger.Error(err.Error())
	}
}

func httpGetXApps(xmurl string) (*[]rtmgr.XApp, error) {
	rtmgr.Logger.Info("Invoked httprestful.httpGetXApps: " + xmurl)
	r, err := myClient.Get(xmurl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 200 {
		rtmgr.Logger.Debug("http client raw response: %v", r)
		var xapps []rtmgr.XApp
		err = json.NewDecoder(r.Body).Decode(&xapps)
		if err != nil {
			rtmgr.Logger.Warn("Json decode failed: " + err.Error())
		}
		rtmgr.Logger.Info("HTTP GET: OK")
		rtmgr.Logger.Debug("httprestful.httpGetXApps returns: %v", xapps)
		return &xapps, err
	}
	rtmgr.Logger.Warn("httprestful got an unexpected http status code: %v", r.StatusCode)
	return nil, nil
}

func retrieveStartupData(xmurl string, nbiif string, fileName string, configfile string, sdlEngine sdl.Engine) error {
	var readErr error
	var maxRetries = 10
	for i := 1; i <= maxRetries; i++ {
		time.Sleep(2 * time.Second)
		xappData, err := httpGetXApps(xmurl)
		if xappData != nil && err == nil {
			pcData, confErr := rtmgr.GetPlatformComponents(configfile)
			if confErr != nil {
				rtmgr.Logger.Error(confErr.Error())
				return confErr
			}
			rtmgr.Logger.Info("Recieved intial xapp data and platform data, writing into SDL.")
			// Combine the xapps data and platform data before writing to the SDL
			ricData := &rtmgr.RicComponents{XApps: *xappData, Pcs: *pcData}
			writeErr := sdlEngine.WriteAll(fileName, ricData)
			if writeErr != nil {
				rtmgr.Logger.Error(writeErr.Error())
			}
			// post subscription req to appmgr
			readErr = PostSubReq(xmurl, nbiif)
			if readErr == nil {
				return nil
			}
		} else if err == nil {
			readErr = errors.New("unexpected HTTP status code")
		} else {
			rtmgr.Logger.Warn("cannot get xapp data due to: " + err.Error())
			readErr = err
		}
	}
	return readErr
}

func (r *HttpRestful) Initialize(xmurl string, nbiif string, fileName string, configfile string,
	sdlEngine sdl.Engine, rpeEngine rpe.Engine, triggerSBI chan<- bool) error {
	err := r.RetrieveStartupData(xmurl, nbiif, fileName, configfile, sdlEngine)
	if err != nil {
		rtmgr.Logger.Error("Exiting as nbi failed to get the initial startup data from the xapp manager: " + err.Error())
		return err
	}

	datach := make(chan *models.XappCallbackData, 10)
	subschan := make(chan *models.XappSubscriptionData, 10)
	subdelchan := make(chan *models.XappSubscriptionData, 10)
	rtmgr.Logger.Info("Launching Rest Http service")
	go func() {
		r.LaunchRest(&nbiif, datach, subschan, subdelchan)
	}()

	go func() {
		for {
			data, err := r.RecvXappCallbackData(datach)
			if err != nil {
				rtmgr.Logger.Error("cannot get data from rest api dute to: " + err.Error())
			} else if data != nil {
				rtmgr.Logger.Debug("Fetching all xApps deployed in xApp Manager through GET operation.")
				alldata, err1 := httpGetXApps(xmurl)
				if alldata != nil && err1 == nil {
					sdlEngine.WriteXApps(fileName, alldata)
					triggerSBI <- true
				}
			}
		}
	}()

	go func() {
		for {
			data := <-subschan
			rtmgr.Logger.Debug("received XApp subscription data")
			addSubscription(&rtmgr.Subs, data)
			triggerSBI <- true
		}
	}()

	go func() {
		for {
			data := <-subdelchan
			rtmgr.Logger.Debug("received XApp subscription delete data")
			delSubscription(&rtmgr.Subs, data)
			triggerSBI <- true
		}
	}()

	return nil
}

func (r *HttpRestful) Terminate() error {
	return nil
}

func addSubscription(subs *rtmgr.SubscriptionList, xappSubData *models.XappSubscriptionData) bool {
	var b = false
	sub := rtmgr.Subscription{SubID: *xappSubData.SubscriptionID, Fqdn: *xappSubData.Address, Port: *xappSubData.Port}
	for _, elem := range *subs {
		if elem == sub {
			rtmgr.Logger.Warn("rtmgr.addSubscription: Subscription already present: %v", elem)
			b = true
		}
	}
	if b == false {
		*subs = append(*subs, sub)
	}
	return b
}

func delSubscription(subs *rtmgr.SubscriptionList, xappSubData *models.XappSubscriptionData) bool {
	rtmgr.Logger.Debug("Deleteing the subscription from the subscriptions list")
	var present = false
	sub := rtmgr.Subscription{SubID: *xappSubData.SubscriptionID, Fqdn: *xappSubData.Address, Port: *xappSubData.Port}
	for i, elem := range *subs {
		if elem == sub {
			present = true
			// Since the order of the list is not important, we are swapping the last element
			// with the matching element and replacing the list with list(n-1) elements.
			(*subs)[len(*subs)-1], (*subs)[i] = (*subs)[i], (*subs)[len(*subs)-1]
			*subs = (*subs)[:len(*subs)-1]
			break
		}
	}
	if present == false {
		rtmgr.Logger.Warn("rtmgr.delSubscription: Subscription = %v, not present in the existing subscriptions", xappSubData)
	}
	return present
}
