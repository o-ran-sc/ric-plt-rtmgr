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
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"net/url"
	"net"
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
	"sync"
	"strings"
)

//var myClient = &http.Client{Timeout: 1 * time.Second}

type HttpRestful struct {
	Engine
	LaunchRest                   LaunchRestHandler
	RecvXappCallbackData         RecvXappCallbackDataHandler
        RecvNewE2Tdata               RecvNewE2TdataHandler 
	ProvideXappHandleHandlerImpl ProvideXappHandleHandlerImpl
	RetrieveStartupData          RetrieveStartupDataHandler
}

func NewHttpRestful() *HttpRestful {
	instance := new(HttpRestful)
	instance.LaunchRest = launchRest
	instance.RecvXappCallbackData = recvXappCallbackData
        instance.RecvNewE2Tdata = recvNewE2Tdata
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
	//xapp.Logger.Info(length)
	for i := 0; i <= length; i++ {
		xapp.Logger.Info("data received")
		// If no data received from the REST, it blocks.
		xappData = <-dataChannel
	}
	if nil != xappData {
		var xapps []rtmgr.XApp
		err := json.Unmarshal([]byte(xappData.XApps), &xapps)
		return &xapps, err
	} else {
		xapp.Logger.Info("No data")
	}

	xapp.Logger.Debug("Nothing received on the Http interface")
	return nil, nil
}

func recvNewE2Tdata(dataChannel <-chan *models.E2tData) (*rtmgr.E2TInstance, error) {
        var e2tData *models.E2tData
        xapp.Logger.Info("data received")

        e2tData = <-dataChannel

        if nil != e2tData {

			e2tinst :=  rtmgr.E2TInstance {
				 Ranlist : make([]string, len(e2tData.RanNamelist)),
			}

            e2tinst.Fqdn = *e2tData.E2TAddress
            e2tinst.Name = "E2TERMINST"
		    copy(e2tinst.Ranlist, e2tData.RanNamelist)

            return &e2tinst,nil

        } else {
                xapp.Logger.Info("No data")
        }

        xapp.Logger.Debug("Nothing received on the Http interface")
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
		xapp.Logger.Debug("Received callback data")
	}
	err := validateXappCallbackData(data)
	if err != nil {
		xapp.Logger.Warn("XApp callback data validation failed: " + err.Error())
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

func validateE2tData(data *models.E2tData) error {

	e2taddress_key := *data.E2TAddress
        if (e2taddress_key == "") {
                return fmt.Errorf("E2TAddress is empty!!!")
        }
	stringSlice := strings.Split(e2taddress_key, ":")
	if (len(stringSlice) == 1) {
		return fmt.Errorf("E2T E2TAddress is not a proper format like ip:port, %v", e2taddress_key )
	}

	_, err := net.LookupIP(stringSlice[0])
	if err != nil {
		return fmt.Errorf("E2T E2TAddress DNS look up failed, E2TAddress: %v", stringSlice[0])
        }

	if checkValidaE2TAddress(e2taddress_key) {
		return fmt.Errorf("E2TAddress already exist!!!, E2TAddress: %v",e2taddress_key)
	}

	return nil
}

func validateDeleteE2tData(data *models.E2tDeleteData) error {

        if (*data.E2TAddress == "") {
                return fmt.Errorf("E2TAddress is empty!!!")
        }

	for _, element := range data.RanAssocList {
		e2taddress_key := *element.E2TAddress
		stringSlice := strings.Split(e2taddress_key, ":")

		if (len(stringSlice) == 1) {
			return fmt.Errorf("E2T Delete - RanAssocList E2TAddress is not a proper format like ip:port, %v", e2taddress_key)
		}


		if !checkValidaE2TAddress(e2taddress_key) {
				return fmt.Errorf("E2TAddress doesn't exist!!!, E2TAddress: %v",e2taddress_key)
		}

	}
	return nil
}

func checkValidaE2TAddress(e2taddress string) bool {

	_, exist := rtmgr.Eps[e2taddress]
	return exist

}

func provideXappSubscriptionHandleImpl(subchan chan<- *models.XappSubscriptionData,
	data *models.XappSubscriptionData) error {
	xapp.Logger.Debug("Invoked provideXappSubscriptionHandleImpl")
	err := validateXappSubscriptionData(data)
	if err != nil {
		xapp.Logger.Error(err.Error())
		return err
	}
	subchan <- data
	//var val = string(*data.Address + ":" + strconv.Itoa(int(*data.Port)))
	xapp.Logger.Debug("Endpoints: %v", rtmgr.Eps)
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
	xapp.Logger.Debug("Invoked deleteXappSubscriptionHandleImpl")
	err := validateXappSubscriptionData(data)
	if err != nil {
		xapp.Logger.Error(err.Error())
		return err
	}

	if !subscriptionExists(data) {
		xapp.Logger.Warn("subscription not found: %d", *data.SubscriptionID)
		err := fmt.Errorf("subscription not found: %d", *data.SubscriptionID)
		return err
	}

	subdelchan <- data
	return nil
}

func createNewE2tHandleHandlerImpl(e2taddchan chan<- *models.E2tData,
        data *models.E2tData) error {
        xapp.Logger.Debug("Invoked createNewE2tHandleHandlerImpl")
        err := validateE2tData(data)
        if err != nil {
                xapp.Logger.Error(err.Error())
                return err
        }
        e2taddchan <- data
        return nil
}

func validateE2TAddressRANListData(assRanE2tData models.RanE2tMap) error {

	xapp.Logger.Debug("Invoked.validateE2TAddressRANListData : %v", assRanE2tData)

	for _, element := range assRanE2tData {
		if *element.E2TAddress == "" {
			return fmt.Errorf("E2T Instance - E2TAddress is empty!!!")
		}

		e2taddress_key := *element.E2TAddress
		if !checkValidaE2TAddress(e2taddress_key) {
			return fmt.Errorf("E2TAddress doesn't exist!!!, E2TAddress: %v",e2taddress_key)
		}

	}
	return nil
}

func associateRanToE2THandlerImpl(assranchan chan<- models.RanE2tMap,
        data models.RanE2tMap) error {
        xapp.Logger.Debug("Invoked associateRanToE2THandlerImpl")
	err := validateE2TAddressRANListData(data)
	if err != nil {
		xapp.Logger.Warn(" Association of RAN to E2T Instance data validation failed: " + err.Error())
		return err
	}
	assranchan <- data
        return nil
}

func disassociateRanToE2THandlerImpl(disassranchan chan<- models.RanE2tMap,
        data models.RanE2tMap) error {
        xapp.Logger.Debug("Invoked disassociateRanToE2THandlerImpl")
	err := validateE2TAddressRANListData(data)
	if err != nil {
		xapp.Logger.Warn(" Disassociation of RAN List from E2T Instance data validation failed: " + err.Error())
		return err
	}
	disassranchan <- data
        return nil
}

func deleteE2tHandleHandlerImpl(e2tdelchan chan<- *models.E2tDeleteData,
        data *models.E2tDeleteData) error {
        xapp.Logger.Debug("Invoked deleteE2tHandleHandlerImpl")

        err := validateDeleteE2tData(data)
        if err != nil {
                xapp.Logger.Error(err.Error())
                return err
        }

        e2tdelchan <- data
        return nil
}

func launchRest(nbiif *string, datach chan<- *models.XappCallbackData, subchan chan<- *models.XappSubscriptionData,
	subdelchan chan<- *models.XappSubscriptionData, e2taddchan chan<- *models.E2tData, assranchan chan<- models.RanE2tMap, disassranchan chan<- models.RanE2tMap, e2tdelchan chan<- *models.E2tDeleteData) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		//log.Fatalln(err)
		xapp.Logger.Error(err.Error())
		os.Exit(1)
	}
	nbiUrl, err := url.Parse(*nbiif)
	if err != nil {
		xapp.Logger.Error(err.Error())
		os.Exit(1)
	}
	api := operations.NewRoutingManagerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port, err = strconv.Atoi(nbiUrl.Port())
	if err != nil {
		xapp.Logger.Error("Invalid NBI RestAPI port")
		os.Exit(1)
	}
	server.Host = "0.0.0.0"
	// set handlers
	api.HandleProvideXappHandleHandler = handle.ProvideXappHandleHandlerFunc(
		func(params handle.ProvideXappHandleParams) middleware.Responder {
			xapp.Logger.Info("Data received on Http interface")
			err := provideXappHandleHandlerImpl(datach, params.XappCallbackData)
			if err != nil {
				xapp.Logger.Error("Invalid XApp callback data: " + err.Error())
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
       api.HandleCreateNewE2tHandleHandler = handle.CreateNewE2tHandleHandlerFunc(
                func(params handle.CreateNewE2tHandleParams) middleware.Responder {
                        err := createNewE2tHandleHandlerImpl(e2taddchan, params.E2tData)
                        if err != nil {
                                return handle.NewCreateNewE2tHandleBadRequest()
                        } else {
				time.Sleep(1 * time.Second)
                                return handle.NewCreateNewE2tHandleCreated()
                        }
                })

       api.HandleAssociateRanToE2tHandleHandler = handle.AssociateRanToE2tHandleHandlerFunc(
		func(params handle.AssociateRanToE2tHandleParams) middleware.Responder {
                        err := associateRanToE2THandlerImpl(assranchan, params.RanE2tList)
			if err != nil {
                                return handle.NewAssociateRanToE2tHandleBadRequest()
                        } else {
				time.Sleep(1 * time.Second)
                                return handle.NewAssociateRanToE2tHandleCreated()
                        }
                })

       api.HandleDissociateRanHandler = handle.DissociateRanHandlerFunc(
	        func(params handle.DissociateRanParams) middleware.Responder {
			err := disassociateRanToE2THandlerImpl(disassranchan, params.DissociateList)
			if err != nil {
                                return handle.NewDissociateRanBadRequest()
                        } else {
				time.Sleep(1 * time.Second)
                                return handle.NewDissociateRanCreated()
                        }
                })

       api.HandleDeleteE2tHandleHandler = handle.DeleteE2tHandleHandlerFunc(
                func(params handle.DeleteE2tHandleParams) middleware.Responder {
                        err := deleteE2tHandleHandlerImpl(e2tdelchan, params.E2tData)
                        if err != nil {
                                return handle.NewDeleteE2tHandleBadRequest()
                        } else {
				time.Sleep(1 * time.Second)
                                return handle.NewDeleteE2tHandleCreated()
                        }
                })
	// start to serve API
	xapp.Logger.Info("Starting the HTTP Rest service")
	if err := server.Serve(); err != nil {
		xapp.Logger.Error(err.Error())
	}
}

func httpGetXApps(xmurl string) (*[]rtmgr.XApp, error) {
	xapp.Logger.Info("Invoked httprestful.httpGetXApps: " + xmurl)
	r, err := myClient.Get(xmurl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 200 {
		xapp.Logger.Debug("http client raw response: %v", r)
		var xapps []rtmgr.XApp
		err = json.NewDecoder(r.Body).Decode(&xapps)
		if err != nil {
			xapp.Logger.Warn("Json decode failed: " + err.Error())
		}
		xapp.Logger.Info("HTTP GET: OK")
		xapp.Logger.Debug("httprestful.httpGetXApps returns: %v", xapps)
		return &xapps, err
	}
	xapp.Logger.Warn("httprestful got an unexpected http status code: %v", r.StatusCode)
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
				xapp.Logger.Error(confErr.Error())
				return confErr
			}
			xapp.Logger.Info("Recieved intial xapp data and platform data, writing into SDL.")
			// Combine the xapps data and platform data before writing to the SDL
			ricData := &rtmgr.RicComponents{XApps: *xappData, Pcs: *pcData, E2Ts:  make(map[string]rtmgr.E2TInstance)}
			writeErr := sdlEngine.WriteAll(fileName, ricData)
			if writeErr != nil {
				xapp.Logger.Error(writeErr.Error())
			}
			// post subscription req to appmgr
			readErr = PostSubReq(xmurl, nbiif)
			if readErr == nil {
				return nil
			}
		} else if err == nil {
			readErr = errors.New("unexpected HTTP status code")
		} else {
			xapp.Logger.Warn("cannot get xapp data due to: " + err.Error())
			readErr = err
		}
	}
	return readErr
}

func (r *HttpRestful) Initialize(xmurl string, nbiif string, fileName string, configfile string,
	sdlEngine sdl.Engine, rpeEngine rpe.Engine, triggerSBI chan<- bool, m *sync.Mutex) error {
	err := r.RetrieveStartupData(xmurl, nbiif, fileName, configfile, sdlEngine)
	if err != nil {
		xapp.Logger.Error("Exiting as nbi failed to get the initial startup data from the xapp manager: " + err.Error())
		return err
	}

	datach := make(chan *models.XappCallbackData, 10)
	subschan := make(chan *models.XappSubscriptionData, 10)
	subdelchan := make(chan *models.XappSubscriptionData, 10)
	e2taddchan := make(chan *models.E2tData, 10)
	associateranchan := make(chan models.RanE2tMap, 10)
	disassociateranchan := make(chan models.RanE2tMap, 10)
	e2tdelchan := make(chan *models.E2tDeleteData, 10)
	xapp.Logger.Info("Launching Rest Http service")
	go func() {
		r.LaunchRest(&nbiif, datach, subschan, subdelchan, e2taddchan, associateranchan, disassociateranchan, e2tdelchan)
	}()

	go func() {
		for {
			data, err := r.RecvXappCallbackData(datach)
			if err != nil {
				xapp.Logger.Error("cannot get data from rest api dute to: " + err.Error())
			} else if data != nil {
				xapp.Logger.Debug("Fetching all xApps deployed in xApp Manager through GET operation.")
				alldata, err1 := httpGetXApps(xmurl)
				if alldata != nil && err1 == nil {
					m.Lock()
					sdlEngine.WriteXApps(fileName, alldata)
					m.Unlock()
					triggerSBI <- true
				}
			}
		}
	}()

	go func() {
		for {
			data := <-subschan
			xapp.Logger.Debug("received XApp subscription data")
			addSubscription(&rtmgr.Subs, data)
			triggerSBI <- true
		}
	}()

	go func() {
		for {
			data := <-subdelchan
			xapp.Logger.Debug("received XApp subscription delete data")
			delSubscription(&rtmgr.Subs, data)
			triggerSBI <- true
		}
	}()

        go func() {
                for {
                        xapp.Logger.Debug("received create New E2T data")

                        data, _ := r.RecvNewE2Tdata(e2taddchan)
                        if data != nil {
				m.Lock()
                                sdlEngine.WriteNewE2TInstance(fileName, data)
				m.Unlock()
                                triggerSBI <- true
                        }
                }
        }()

        go func() {
                for {
			data := <-associateranchan
                        xapp.Logger.Debug("received associate RAN list to E2T instance mapping from E2 Manager")
			m.Lock()
                        sdlEngine.WriteAssRANToE2TInstance(fileName, data)
			m.Unlock()
                        triggerSBI <- true
                }
        }()

        go func() {
                for {

			data := <-disassociateranchan
                        xapp.Logger.Debug("received disassociate RANs from E2T instance")
			m.Lock()
                        sdlEngine.WriteDisAssRANFromE2TInstance(fileName, data)
			m.Unlock()
                        triggerSBI <- true
                }
        }()

        go func() {
                for {
                        xapp.Logger.Debug("received Delete E2T data")

			data := <-e2tdelchan
                        if data != nil {
				m.Lock()
                                sdlEngine.WriteDeleteE2TInstance(fileName, data)
				m.Unlock()
                                triggerSBI <- true
                        }
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
			xapp.Logger.Warn("rtmgr.addSubscription: Subscription already present: %v", elem)
			b = true
		}
	}
	if b == false {
		*subs = append(*subs, sub)
	}
	return b
}

func delSubscription(subs *rtmgr.SubscriptionList, xappSubData *models.XappSubscriptionData) bool {
	xapp.Logger.Debug("Deleteing the subscription from the subscriptions list")
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
		xapp.Logger.Warn("rtmgr.delSubscription: Subscription = %v, not present in the existing subscriptions", xappSubData)
	}
	return present
}
