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
  Mnemonic:	httprestful.go
  Abstract:	HTTP Restful API NBI implementation
                Based on Swagger generated code
  Date:		25 March 2019
*/

package nbi

import (
	"fmt"
	"os"
	"time"
	"net/url"
	"strconv"
	"errors"
	"encoding/json"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/sdl"
	"routing-manager/pkg/models"
	"routing-manager/pkg/restapi"
	"routing-manager/pkg/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"routing-manager/pkg/restapi/operations/handle"
	loads "github.com/go-openapi/loads"
)

//var myClient = &http.Client{Timeout: 1 * time.Second}

type HttpRestful struct {
	NbiEngine
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
		return fmt.Errorf("Invalid Data field: \"%s\"", callbackData.XApps)
	}
	var xapps []rtmgr.XApp
        err := json.Unmarshal([]byte(callbackData.XApps), &xapps)
        if err != nil {
		return fmt.Errorf("Unmarshal failed: \"%s\"", err.Error())
	}
	return nil
}

func provideXappHandleHandlerImpl(datach chan<- *models.XappCallbackData, data *models.XappCallbackData) error {
	if data != nil {
		rtmgr.Logger.Debug("Received callback data")
	}
	err := validateXappCallbackData(data)
	if err != nil {
		rtmgr.Logger.Debug("XApp callback data validation failed: "+err.Error())
		return err
	} else {
		datach<-data
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

func launchRest(nbiif *string, datach chan<- *models.XappCallbackData, subchan chan<- *models.XappSubscriptionData) {
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
        server.Host = nbiUrl.Hostname()
        // set handlers
        api.HandleProvideXappHandleHandler = handle.ProvideXappHandleHandlerFunc(
                func(params handle.ProvideXappHandleParams) middleware.Responder {
                rtmgr.Logger.Info("Data received on Http interface")
		err := provideXappHandleHandlerImpl(datach, params.XappCallbackData)
		if err != nil {
			rtmgr.Logger.Error("Invalid XApp callback data: "+err.Error())
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
				return handle.NewGetHandlesOK()
			}
		})
        // start to serve API
        rtmgr.Logger.Info("Starting the HTTP Rest service")
        if err := server.Serve(); err != nil {
                rtmgr.Logger.Error(err.Error())
        }
}

func httpGetXapps(xmurl string) (*[]rtmgr.XApp, error) {
        rtmgr.Logger.Info("Invoked httpgetter.fetchXappList: " + xmurl)
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
                rtmgr.Logger.Debug("httpgetter.fetchXappList returns: %v", xapps)
                return &xapps, err
        }
        rtmgr.Logger.Warn("httpgetter got an unexpected http status code: %v", r.StatusCode)
        return nil, nil
}

func retrieveStartupData(xmurl string, nbiif string, fileName string, configfile string, sdlEngine sdl.SdlEngine) error {
        var readErr error
        var maxRetries = 10

                for i := 1; i <= maxRetries; i++ {
                        time.Sleep(2 * time.Second)

                        xappData, err := httpGetXapps(xmurl)

                        if xappData != nil && err == nil {
				pcData, confErr := rtmgr.GetPlatformComponents(configfile)
				if confErr != nil {
					rtmgr.Logger.Error(confErr.Error())
					return confErr
				}

                                rtmgr.Logger.Info("Recieved intial xapp data and platform data, writing into SDL.")
				// Combine the xapps data and platform data before writing to the SDL
				ricData := &rtmgr.RicComponents{Xapps: *xappData, Pcs: *pcData}

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
                                readErr = errors.New("Unexpected HTTP status code")
                        } else {
                                rtmgr.Logger.Warn("cannot get xapp data due to: " + err.Error())
                                readErr = err
                        }
                }
        return readErr
}

func (r *HttpRestful) Initialize(xmurl string, nbiif string, fileName string, configfile string,
				 sdlEngine sdl.SdlEngine, rpeEngine rpe.RpeEngine, triggerSBI chan<- bool) error {
	err := r.RetrieveStartupData(xmurl, nbiif, fileName, configfile, sdlEngine)
	if err != nil {
		rtmgr.Logger.Error("Exiting as nbi failed to get the intial startup data from the xapp manager: " + err.Error())
		return err
	}

	datach := make(chan *models.XappCallbackData, 10)
	subschan := make(chan *models.XappSubscriptionData, 10)
	rtmgr.Logger.Info("Launching Rest Http service")
	go func() {
		r.LaunchRest(&nbiif, datach, subschan)
	}()

	go func() {
		for {
			data, err := r.RecvXappCallbackData(datach)
			if err != nil {
				rtmgr.Logger.Error("cannot get data from rest api dute to: " + err.Error())
			} else if data != nil {
				sdlEngine.WriteXapps(fileName, data)
				triggerSBI <- true
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

	return nil
}

func (r *HttpRestful) Terminate() error {
	return nil
}

func addSubscription(subs *rtmgr.SubscriptionList, xappSubData *models.XappSubscriptionData) bool {
	var b bool = false
	sub := rtmgr.Subscription{SubID:*xappSubData.SubscriptionID, Fqdn:*xappSubData.Address, Port:*xappSubData.Port,}
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

