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
	xfmodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	//"github.com/go-openapi/loads"
	//"github.com/go-openapi/runtime/middleware"
	"net"
	//"net/url"
	//"os"
	"routing-manager/pkg/models"
	//"routing-manager/pkg/restapi"
	//"routing-manager/pkg/restapi/operations"
	//"routing-manager/pkg/restapi/operations/debug"
	//"routing-manager/pkg/restapi/operations/handle"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sdl"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HttpRestful struct {
	Engine
	//LaunchRest          LaunchRestHandler
	RetrieveStartupData RetrieveStartupDataHandler
}

func NewHttpRestful() *HttpRestful {
	instance := new(HttpRestful)
	//instance.LaunchRest = launchRest
	instance.RetrieveStartupData = retrieveStartupData
	return instance
}

func recvXappCallbackData(xappData *models.XappCallbackData) (*[]rtmgr.XApp, error) {
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

func recvNewE2Tdata(e2tData *models.E2tData) (*rtmgr.E2TInstance, string, error) {
	var str string
	xapp.Logger.Info("Data received")

	if nil != e2tData {

		e2tinst := rtmgr.E2TInstance{
			Ranlist: make([]string, len(e2tData.RanNamelist)),
		}

		e2tinst.Fqdn = *e2tData.E2TAddress
		e2tinst.Name = "E2TERMINST"
		copy(e2tinst.Ranlist, e2tData.RanNamelist)
		if len(e2tData.RanNamelist) > 0 {
			var meidar string
			for _, meid := range e2tData.RanNamelist {
				meidar += meid + " "
			}
			str = "mme_ar|" + *e2tData.E2TAddress + "|" + strings.TrimSuffix(meidar, " ")
		}
		return &e2tinst, str, nil

	} else {
		xapp.Logger.Info("No data")
	}

	xapp.Logger.Debug("Nothing received on the Http interface")
	return nil, str, nil
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

func ProvideXappHandleHandlerImpl(data *models.XappCallbackData) error {
	if data == nil {
		xapp.Logger.Debug("Received callback data")
		return nil
	}
	err := validateXappCallbackData(data)
	if err != nil {
		xapp.Logger.Warn("XApp callback data validation failed: " + err.Error())
		return err
	} else {
		appdata, err := recvXappCallbackData(data)
		if err != nil {
			xapp.Logger.Error("Cannot get data from rest api dute to: " + err.Error())
		} else if appdata != nil {
			xapp.Logger.Debug("Fetching all xApps deployed in xApp Manager through GET operation")
			alldata, err1 := httpGetXApps(xapp.Config.GetString("xmurl"))
			if alldata != nil && err1 == nil {
				m.Lock()
				sdlEngine.WriteXApps(xapp.Config.GetString("rtfile"), alldata)
				m.Unlock()
				updateEp()
				return sendRoutesToAll()
				//return sendPartialRoutesToAll(nil, rtmgr.XappType)
			}
		}

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

func validateE2tData(data *models.E2tData) (error, bool) {

	e2taddress_key := *data.E2TAddress
	if e2taddress_key == "" {
		return fmt.Errorf("E2TAddress is empty!!!"), false
	}
	stringSlice := strings.Split(e2taddress_key, ":")
	if len(stringSlice) == 1 {
		return fmt.Errorf("E2T E2TAddress is not a proper format like ip:port, %v", e2taddress_key), false
	}

	_, err := net.LookupIP(stringSlice[0])
	if err != nil {
		return fmt.Errorf("E2T E2TAddress DNS look up failed, E2TAddress: %v", stringSlice[0]), false
	}

	if checkValidaE2TAddress(e2taddress_key) {
		return fmt.Errorf("E2TAddress already exist!!!, E2TAddress: %v", e2taddress_key), true
	}

	return nil, false
}

func validateDeleteE2tData(data *models.E2tDeleteData) error {

	if *data.E2TAddress == "" {
		return fmt.Errorf("E2TAddress is empty!!!")
	}

	for _, element := range data.RanAssocList {
		e2taddress_key := *element.E2TAddress
		stringSlice := strings.Split(e2taddress_key, ":")

		if len(stringSlice) == 1 {
			return fmt.Errorf("E2T Delete - RanAssocList E2TAddress is not a proper format like ip:port, %v", e2taddress_key)
		}

		if !checkValidaE2TAddress(e2taddress_key) {
			return fmt.Errorf("E2TAddress doesn't exist!!!, E2TAddress: %v", e2taddress_key)
		}

	}
	return nil
}

func checkValidaE2TAddress(e2taddress string) bool {

	_, exist := rtmgr.Eps[e2taddress]
	return exist

}

func ProvideXappSubscriptionHandleImpl(data *models.XappSubscriptionData) error {
	xapp.Logger.Debug("Invoked provideXappSubscriptionHandleImpl")
	err := validateXappSubscriptionData(data)
	if err != nil {
		xapp.Logger.Error(err.Error())
		return err
	}
	xapp.Logger.Debug("Received xApp subscription data")
	addSubscription(&rtmgr.Subs, data)
	xapp.Logger.Debug("Endpoints: %v", rtmgr.Eps)
	updateEp()
	return sendPartialRoutesToAll(data, rtmgr.SubsType)
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

func DeleteXappSubscriptionHandleImpl(data *models.XappSubscriptionData) error {
	xapp.Logger.Debug("Invoked deleteXappSubscriptionHandleImpl")
	err := validateXappSubscriptionData(data)
	if err != nil {
		xapp.Logger.Error(err.Error())
		return err
	}

	if !subscriptionExists(data) {
		xapp.Logger.Warn("Subscription not found: %d", *data.SubscriptionID)
		err := fmt.Errorf("Subscription not found: %d", *data.SubscriptionID)
		return err
	}

	xapp.Logger.Debug("Received xApp subscription delete data")
	delSubscription(&rtmgr.Subs, data)
	updateEp()
	return sendRoutesToAll()

}

func UpdateXappSubscriptionHandleImpl(data *models.XappList, subid uint16) error {
	xapp.Logger.Debug("Invoked updateXappSubscriptionHandleImpl")

	var fqdnlist []rtmgr.FqDn
	for _, item := range *data {
		fqdnlist = append(fqdnlist, rtmgr.FqDn(*item))
	}
	xapplist := rtmgr.XappList{SubscriptionID: subid, FqdnList: fqdnlist}
	var subdata models.XappSubscriptionData
	var id int32
	id = int32(subid)
	subdata.SubscriptionID = &id
	for _, items := range fqdnlist {
		subdata.Address = items.Address
		subdata.Port = items.Port
		err := validateXappSubscriptionData(&subdata)
		if err != nil {
			xapp.Logger.Error(err.Error())
			return err
		}
	}
	xapp.Logger.Debug("Received XApp subscription Merge data")
	updateSubscription(&xapplist)
	updateEp()
	return sendRoutesToAll()
}

func CreateNewE2tHandleHandlerImpl(data *models.E2tData) error {
	xapp.Logger.Debug("Invoked createNewE2tHandleHandlerImpl")
	err, IsDuplicate := validateE2tData(data)
	if IsDuplicate == true {
		updateEp()
		//return sendPartialRoutesToAll(nil, rtmgr.E2Type)
		return sendRoutesToAll()
	}

	if err != nil {
		xapp.Logger.Error(err.Error())
		return err
	}
	//e2taddchan <- data
	e2data, meiddata, _ := recvNewE2Tdata(data)
	xapp.Logger.Debug("Received create New E2T data")
	m.Lock()
	sdlEngine.WriteNewE2TInstance(xapp.Config.GetString("rtfile"), e2data, meiddata)
	m.Unlock()
	updateEp()
	//sendPartialRoutesToAll(nil, rtmgr.E2Type)
	sbi.Conn.Lock()
	rtmgr.RMRConnStatus[*data.E2TAddress] = false
	sbi.Conn.Unlock()

	sendRoutesToAll()
	time.Sleep(10 * time.Second)
	if rtmgr.RMRConnStatus[*data.E2TAddress] == true {
    	sbi.Conn.Lock()
    	delete(rtmgr.RMRConnStatus,*data.E2TAddress)
    	sbi.Conn.Unlock()
    	xapp.Logger.Debug("RMRConnStatus Map after = %v",rtmgr.RMRConnStatus)
    	return nil
    }

	}

	return errors.New("Error while adding new E2T " + *data.E2TAddress)

}

func validateE2TAddressRANListData(assRanE2tData models.RanE2tMap) error {

	xapp.Logger.Debug("Invoked.validateE2TAddressRANListData : %v", assRanE2tData)

	for _, element := range assRanE2tData {
		if *element.E2TAddress == "" {
			return fmt.Errorf("E2T Instance - E2TAddress is empty!!!")
		}

		e2taddress_key := *element.E2TAddress
		if !checkValidaE2TAddress(e2taddress_key) {
			return fmt.Errorf("E2TAddress doesn't exist!!!, E2TAddress: %v", e2taddress_key)
		}

	}
	return nil
}

func AssociateRanToE2THandlerImpl(data models.RanE2tMap) error {
	xapp.Logger.Debug("Invoked associateRanToE2THandlerImpl")
	err := validateE2TAddressRANListData(data)
	if err != nil {
		xapp.Logger.Warn(" Association of RAN to E2T Instance data validation failed: " + err.Error())
		return err
	}
	xapp.Logger.Debug("Received associate RAN list to E2T instance mapping from E2 Manager")
	m.Lock()
	sdlEngine.WriteAssRANToE2TInstance(xapp.Config.GetString("rtfile"), data)
	m.Unlock()
	updateEp()
	return sendRoutesToAll()

}

func DisassociateRanToE2THandlerImpl(data models.RanE2tMap) error {
	xapp.Logger.Debug("Invoked disassociateRanToE2THandlerImpl")
	err := validateE2TAddressRANListData(data)
	if err != nil {
		xapp.Logger.Warn(" Disassociation of RAN List from E2T Instance data validation failed: " + err.Error())
		return err
	}
	xapp.Logger.Debug("Received disassociate RANs from E2T instance")
	m.Lock()
	sdlEngine.WriteDisAssRANFromE2TInstance(xapp.Config.GetString("rtfile"), data)
	m.Unlock()
	updateEp()
	return sendRoutesToAll()

}

func DeleteE2tHandleHandlerImpl(data *models.E2tDeleteData) error {
	xapp.Logger.Debug("Invoked deleteE2tHandleHandlerImpl")

	err := validateDeleteE2tData(data)
	if err != nil {
		xapp.Logger.Error(err.Error())
		return err
	}
	m.Lock()
	sdlEngine.WriteDeleteE2TInstance(xapp.Config.GetString("rtfile"), data)
	m.Unlock()
	updateEp()
	return sendRoutesToAll()

}

func DumpDebugData() (models.Debuginfo, error) {
	var response models.Debuginfo
	sdlEngine, _ := sdl.GetSdl("file")
	rpeEngine, _ := rpe.GetRpe("rmrpush")
	data, err := sdlEngine.ReadAll(xapp.Config.GetString("rtfile"))
	if data == nil {
		if err != nil {
			xapp.Logger.Error("Cannot get data from sdl interface due to: " + err.Error())
			return response, err
		} else {
			xapp.Logger.Debug("Cannot get data from sdl interface")
			return response, errors.New("Cannot get data from sdl interface")
		}
	}
	response.RouteTable = *rpeEngine.GeneratePolicies(rtmgr.Eps, data)

	prettyJSON, err := json.MarshalIndent(data, "", "")
	response.RouteConfigs = string(prettyJSON)

	return response, err
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

func httpGetE2TList(e2murl string) (*[]rtmgr.E2tIdentity, error) {
	xapp.Logger.Info("Invoked httprestful.httpGetE2TList: " + e2murl)
	r, err := myClient.Get(e2murl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 200 {
		xapp.Logger.Debug("http client raw response: %v", r)
		var E2Tlist []rtmgr.E2tIdentity
		err = json.NewDecoder(r.Body).Decode(&E2Tlist)
		if err != nil {
			xapp.Logger.Warn("Json decode failed: " + err.Error())
		}
		xapp.Logger.Info("HTTP GET: OK")
		xapp.Logger.Debug("httprestful.httpGetE2TList returns: %v", E2Tlist)
		return &E2Tlist, err
	}
	xapp.Logger.Warn("httprestful got an unexpected http status code: %v", r.StatusCode)
	return nil, nil
}

func PopulateE2TMap(e2tDataList *[]rtmgr.E2tIdentity, e2ts map[string]rtmgr.E2TInstance, meids *[]string) {
	xapp.Logger.Info("Invoked httprestful.PopulateE2TMap ")

	for _, e2tData := range *e2tDataList {
		var str string

		e2tinst := rtmgr.E2TInstance{
			Ranlist: make([]string, len(e2tData.Rannames)),
		}

		e2tinst.Fqdn = e2tData.E2taddress
		e2tinst.Name = "E2TERMINST"
		copy(e2tinst.Ranlist, e2tData.Rannames)

		if len(e2tData.Rannames) > 0 {
			var meidar string
			for _, meid := range e2tData.Rannames {
				meidar += meid + " "
			}
			str = "mme_ar|" + e2tData.E2taddress + "|" + strings.TrimSuffix(meidar, " ")
			*meids = append(*meids, str)
		}

		e2ts[e2tinst.Fqdn] = e2tinst
	}
	xapp.Logger.Info("MEID's retrieved are %v", *meids)
}

func retrieveStartupData(xmurl string, nbiif string, fileName string, configfile string, e2murl string, sdlEngine sdl.Engine) error {
	xapp.Logger.Info("Invoked retrieveStartupData ")
	var readErr error
	var err error
	var maxRetries = 10
	var xappData *[]rtmgr.XApp
	xappData = new([]rtmgr.XApp)
	xapp.Logger.Info("Trying to fetch XApps data from XAPP manager")
	for i := 1; i <= maxRetries; i++ {
		time.Sleep(2 * time.Second)

		readErr = nil
		xappData, err = httpGetXApps(xmurl)
		if xappData != nil && err == nil {
			break
		} else if err == nil {
			readErr = errors.New("unexpected HTTP status code")
		} else {
			xapp.Logger.Warn("Cannot get xapp data due to: " + err.Error())
			readErr = err
		}
	}

	if readErr != nil {
		return readErr
	}

	var meids []string
	e2ts := make(map[string]rtmgr.E2TInstance)
	xapp.Logger.Info("Trying to fetch E2T data from E2manager")
	for i := 1; i <= maxRetries; i++ {

		readErr = nil
		e2tDataList, err := httpGetE2TList(e2murl)
		if e2tDataList != nil && err == nil {
			PopulateE2TMap(e2tDataList, e2ts, &meids)
			break
		} else if err == nil {
			readErr = errors.New("unexpected HTTP status code")
		} else {
			xapp.Logger.Warn("Cannot get E2T data from E2M due to: " + err.Error())
			readErr = err
		}
		time.Sleep(2 * time.Second)
	}

	if readErr != nil {
		return readErr
	}

	pcData, confErr := rtmgr.GetPlatformComponents(configfile)
	if confErr != nil {
		xapp.Logger.Error(confErr.Error())
		return confErr
	}
	xapp.Logger.Info("Recieved intial xapp data, E2T data and platform data, writing into SDL.")
	// Combine the xapps data and platform data before writing to the SDL
	ricData := &rtmgr.RicComponents{XApps: *xappData, Pcs: *pcData, E2Ts: e2ts, MeidMap: meids}
	writeErr := sdlEngine.WriteAll(fileName, ricData)
	if writeErr != nil {
		xapp.Logger.Error(writeErr.Error())
	}

	xapp.Logger.Info("Trying to fetch Subscriptions data from Subscription manager")
	for i := 1; i <= maxRetries; i++ {
		readErr = nil
		sub_list, err := xapp.Subscription.QuerySubscriptions()

		if sub_list != nil && err == nil {
			PopulateSubscription(sub_list)
			break
		} else {
			readErr = err
			xapp.Logger.Warn("Cannot get xapp data due to: " + readErr.Error())
		}
		time.Sleep(2 * time.Second)
	}

	if readErr != nil {
		return readErr
	}

	// post subscription req to appmgr
	readErr = PostSubReq(xmurl, nbiif)
	if readErr != nil {
		return readErr
	}

	//rlist := make(map[string]string)
	xapp.Logger.Info("Reading SDL for any routes")
	rlist, sdlerr := xapp.SdlStorage.Read(rtmgr.RTMGR_SDL_NS, "routes")
	readErr = sdlerr
	if readErr == nil {
		xapp.Logger.Info("Value is %s", rlist["routes"])
		if rlist["routes"] != nil {
			formstring := fmt.Sprintf("%s", rlist["routes"])
			xapp.Logger.Info("Value of formed string = %s", formstring)
			newstring := strings.Split(formstring, " ")
			for i, _ := range newstring {
				xapp.Logger.Info("Value of formed string in loop = %s", newstring)
				rtmgr.DynamicRouteList = append(rtmgr.DynamicRouteList, newstring[i])
			}
		}

		return nil
	}

	return readErr
}

func (r *HttpRestful) Initialize(xmurl string, nbiif string, fileName string, configfile string, e2murl string,
	sdlEngine sdl.Engine, rpeEngine rpe.Engine, m *sync.Mutex) error {
	err := r.RetrieveStartupData(xmurl, nbiif, fileName, configfile, e2murl, sdlEngine)
	if err != nil {
		xapp.Logger.Error("Exiting as nbi failed to get the initial startup data from the xapp manager: " + err.Error())
		return err
	}

	return nil
}

func (r *HttpRestful) Terminate() error {
	return nil
}

func addSubscription(subs *rtmgr.SubscriptionList, xappSubData *models.XappSubscriptionData) bool {
	xapp.Logger.Debug("Adding the subscription into the subscriptions list")
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

func updateSubscription(data *rtmgr.XappList) {

	var subdata models.XappSubscriptionData
	var id int32
	var matchingsubid, deletecount uint8
	id = int32(data.SubscriptionID)
	subdata.SubscriptionID = &id
	for _, subs := range rtmgr.Subs {
		if int32(data.SubscriptionID) == subs.SubID {
			matchingsubid++
		}
	}

	for deletecount < matchingsubid {
		for _, subs := range rtmgr.Subs {
			if int32(data.SubscriptionID) == subs.SubID {
				subdata.SubscriptionID = &subs.SubID
				subdata.Address = &subs.Fqdn
				subdata.Port = &subs.Port
				xapp.Logger.Debug("Deleting Subscription List has %v", subdata)
				delSubscription(&rtmgr.Subs, &subdata)
				break
			}
		}
		deletecount++
	}

	for _, items := range data.FqdnList {
		subdata.Address = items.Address
		subdata.Port = items.Port
		xapp.Logger.Debug("Adding Subscription List has %v", subdata)
		addSubscription(&rtmgr.Subs, &subdata)
	}

}

func PopulateSubscription(sub_list xfmodel.SubscriptionList) {
	for _, sub_row := range sub_list {
		var subdata models.XappSubscriptionData
		id := int32(sub_row.SubscriptionID)
		subdata.SubscriptionID = &id
		for _, ep := range sub_row.ClientEndpoint {
			stringSlice := strings.Split(ep, ":")
			subdata.Address = &stringSlice[0]
			intportval, _ := strconv.Atoi(stringSlice[1])
			value := uint16(intportval)
			subdata.Port = &value
			xapp.Logger.Debug("Adding Subscription List has Address :%v, port :%v, SubscriptionID :%v ", subdata.Address, subdata.Address, subdata.SubscriptionID)
			addSubscription(&rtmgr.Subs, &subdata)
		}
	}
}

func Adddelrmrroute(routelist models.Routelist, rtflag bool) error {
	xapp.Logger.Info("Updating rmr route with Route list: %v,flag: %v", routelist, rtflag)
	for _, rlist := range routelist {
		var subid int32
		var data string
		if rlist.SubscriptionID == 0 {
			subid = -1
		} else {
			subid = rlist.SubscriptionID
		}
		if rlist.SenderEndPoint == "" && rlist.SubscriptionID != 0 {
			data = fmt.Sprintf("mse|%d|%d|%s\n", *rlist.MessageType, rlist.SubscriptionID, *rlist.TargetEndPoint)
		} else if rlist.SenderEndPoint == "" && rlist.SubscriptionID == 0 {
			data = fmt.Sprintf("mse|%d|-1|%s\n", *rlist.MessageType, *rlist.TargetEndPoint)
		} else {
			data = fmt.Sprintf("mse|%d,%s|%d|%s\n", *rlist.MessageType, rlist.SenderEndPoint, subid, *rlist.TargetEndPoint)
		}
		err := checkrepeatedroute(data)

		if rtflag == true {
			if err == true {
				xapp.Logger.Info("Given route %s is a duplicate", data)
			}
			rtmgr.DynamicRouteList = append(rtmgr.DynamicRouteList, data)
			routearray := strings.Join(rtmgr.DynamicRouteList, " ")
			xapp.SdlStorage.Store(rtmgr.RTMGR_SDL_NS, "routes", routearray)
		} else {
			if err == true {
				xapp.Logger.Info("route %s deleted successfully", data)
				routearray := strings.Join(rtmgr.DynamicRouteList, " ")
				xapp.SdlStorage.Store(rtmgr.RTMGR_SDL_NS, "routes", routearray)
			} else {
				xapp.Logger.Info("No such route: %s", data)
				return errors.New("No such route: " + data)
			}

		}
	}
	updateEp()
	return sendRoutesToAll()
}

func checkrepeatedroute(data string) bool {
	for i := 0; i < len(rtmgr.DynamicRouteList); i++ {
		if rtmgr.DynamicRouteList[i] == data {
			rtmgr.DynamicRouteList[i] = rtmgr.DynamicRouteList[len(rtmgr.DynamicRouteList)-1]
			rtmgr.DynamicRouteList[len(rtmgr.DynamicRouteList)-1] = ""
			rtmgr.DynamicRouteList = rtmgr.DynamicRouteList[:len(rtmgr.DynamicRouteList)-1]
			return true
		}
	}
	return false
}
