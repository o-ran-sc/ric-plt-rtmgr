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
  Mnemonic:	rpe.go
  Abstract:	Contains RPE (Route Policy Engine) module definitions and generic RPE components
  Date:		16 March 2019
*/

package rpe

import (
	"errors"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sbi"
	"strconv"
)

var (
	SupportedRpes = []*EngineConfig{
		{
			Name:        "rmrpush",
			Version:     "pubsush",
			Protocol:    "rmruta",
			Instance:    NewRmrPush(),
			IsAvailable: true,
		},
	}
)

func GetRpe(rpeName string) (Engine, error) {
	for _, rpe := range SupportedRpes {
		if rpe.Name == rpeName && rpe.IsAvailable {
			return rpe.Instance, nil
		}
	}
	return nil, errors.New("SBI:" + rpeName + " is not supported or still not a available")
}

type Rpe struct {
}

func getEndpointByName(eps *rtmgr.Endpoints, name string) *rtmgr.Endpoint {
	for _, ep := range *eps {
		if ep.Name == name {
			rtmgr.Logger.Debug("name: %s", ep.Name)
			rtmgr.Logger.Debug("ep: %v", ep)
			return ep
		}
	}
	return nil
}

func getEndpointByUuid(uuid string) *rtmgr.Endpoint {
	endPoints := rtmgr.Eps
	for _, ep := range endPoints {
		if ep.Uuid == uuid {
			rtmgr.Logger.Debug("name: %s", ep.Uuid)
			rtmgr.Logger.Debug("ep: %v", ep)
			return ep
		}
	}
	return nil
}

func (r *Rpe) addRoute(messageType string, tx *rtmgr.Endpoint, rx *rtmgr.Endpoint, routeTable *rtmgr.RouteTable, subId int32) {
	txList := rtmgr.EndpointList{*tx}
	rxList := []rtmgr.EndpointList{[]rtmgr.Endpoint{*rx}}
	messageId := rtmgr.MessageTypes[messageType]
	route := rtmgr.RouteTableEntry{
		MessageType: messageId,
		TxList:      txList,
		RxGroups:    rxList,
		SubID:       subId}
	*routeTable = append(*routeTable, route)
	rtmgr.Logger.Debug("Route added: MessageTyp: %v, Tx: %v, Rx: %v, SubId: %v", messageId, tx.Uuid, rx.Uuid, subId)
	rtmgr.Logger.Trace("Route added: MessageTyp: %v, Tx: %v, Rx: %v, SubId: %v", messageId, tx, rx, subId)
}

func (r *Rpe) generateXappRoutes(xAppEp *rtmgr.Endpoint, e2TermEp *rtmgr.Endpoint, subManEp *rtmgr.Endpoint, routeTable *rtmgr.RouteTable) {
	rtmgr.Logger.Debug("rpe.generateXappRoutes invoked")
	rtmgr.Logger.Debug("Endpoint: %v, xAppType: %v", xAppEp.Name, xAppEp.XAppType)
	if xAppEp.XAppType != sbi.PlatformType && len(xAppEp.TxMessages) > 0 && len(xAppEp.RxMessages) > 0 {
		//xApp -> Subscription Manager
		r.addRoute("RIC_SUB_REQ", xAppEp, subManEp, routeTable, -1)
		r.addRoute("RIC_SUB_DEL_REQ", xAppEp, subManEp, routeTable, -1)
		//xApp -> E2 Termination
		r.addRoute("RIC_CONTROL_REQ", xAppEp, e2TermEp, routeTable, -1)
	}
}

func (r *Rpe) generateSubscriptionRoutes(e2TermEp *rtmgr.Endpoint, subManEp *rtmgr.Endpoint, routeTable *rtmgr.RouteTable) {
	rtmgr.Logger.Debug("rpe.addSubscriptionRoutes invoked")
	subscriptionList := &rtmgr.Subs
	for _, subscription := range *subscriptionList {
		rtmgr.Logger.Debug("Subscription: %v", subscription)
		xAppUuid := subscription.Fqdn + ":" + strconv.Itoa(int(subscription.Port))
		rtmgr.Logger.Debug("xApp UUID: %v", xAppUuid)
		xAppEp := getEndpointByUuid(xAppUuid)
		//Subscription Manager -> xApp
		r.addRoute("RIC_SUB_RESP", subManEp, xAppEp, routeTable, subscription.SubID)
		r.addRoute("RIC_SUB_FAILURE", subManEp, xAppEp, routeTable, subscription.SubID)
		r.addRoute("RIC_SUB_DEL_RESP", subManEp, xAppEp, routeTable, subscription.SubID)
		r.addRoute("RIC_SUB_DEL_FAILURE", subManEp, xAppEp, routeTable, subscription.SubID)
		//E2 Termination -> xApp
		r.addRoute("RIC_INDICATION", e2TermEp, xAppEp, routeTable, subscription.SubID)
		r.addRoute("RIC_CONTROL_ACK", e2TermEp, xAppEp, routeTable, subscription.SubID)
		r.addRoute("RIC_CONTROL_FAILURE", e2TermEp, xAppEp, routeTable, subscription.SubID)
	}
}

func (r *Rpe) generatePlatformRoutes(e2TermEp *rtmgr.Endpoint, subManEp *rtmgr.Endpoint, e2ManEp *rtmgr.Endpoint, ueManEp *rtmgr.Endpoint, routeTable *rtmgr.RouteTable) {
	rtmgr.Logger.Debug("rpe.generatePlatformRoutes invoked")
	//Platform Routes --- Subscription Routes
	//Subscription Manager -> E2 Termination
	r.addRoute("RIC_SUB_REQ", subManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_SUB_DEL_REQ", subManEp, e2TermEp, routeTable, -1)
	//E2 Termination -> Subscription Manager
	r.addRoute("RIC_SUB_RESP", e2TermEp, subManEp, routeTable, -1)
	r.addRoute("RIC_SUB_DEL_RESP", e2TermEp, subManEp, routeTable, -1)
	r.addRoute("RIC_SUB_FAILURE", e2TermEp, subManEp, routeTable, -1)
	r.addRoute("RIC_SUB_DEL_FAILURE", e2TermEp, subManEp, routeTable, -1)
	//TODO: UE Man Routes removed (since it is not existing)
	//UE Manager -> Subscription Manager
	//r.addRoute("RIC_SUB_REQ", ueManEp, subManEp, routeTable)
	//r.addRoute("RIC_SUB_DEL_REQ", ueManEp, subManEp, routeTable)
	////UE Manager -> E2 Termination
	//r.addRoute("RIC_CONTROL_REQ", ueManEp, e2TermEp, routeTable)

	//Platform Routes --- X2 Routes
	//E2 Manager -> E2 Termination
	r.addRoute("RIC_X2_SETUP_REQ", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_X2_SETUP_RESP", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_X2_SETUP_FAILURE", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_X2_RESET_RESP", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_ENDC_X2_SETUP_REQ", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_ENDC_X2_SETUP_RESP", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_ENDC_X2_SETUP_FAILURE", e2ManEp, e2TermEp, routeTable, -1)
	//E2 Termination -> E2 Manager
	r.addRoute("RIC_X2_SETUP_REQ", e2TermEp, e2ManEp, routeTable, -1)
	r.addRoute("RIC_X2_SETUP_RESP", e2TermEp, e2ManEp, routeTable, -1)
	r.addRoute("RIC_X2_RESET", e2TermEp, e2ManEp, routeTable, -1)
	r.addRoute("RIC_X2_RESOURCE_STATUS_RESPONSE", e2TermEp, e2ManEp, routeTable, -1)
	r.addRoute("RIC_X2_RESET_RESP", e2TermEp, e2ManEp, routeTable, -1)
	r.addRoute("RIC_ENDC_X2_SETUP_REQ", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_ENDC_X2_SETUP_RESP", e2ManEp, e2TermEp, routeTable, -1)
	r.addRoute("RIC_ENDC_X2_SETUP_FAILURE", e2ManEp, e2TermEp, routeTable, -1)
}

func (r *Rpe) generateRouteTable(endPointList rtmgr.Endpoints) *rtmgr.RouteTable {
	rtmgr.Logger.Debug("rpe.generateRouteTable invoked")
	rtmgr.Logger.Debug("Endpoint List:  %v", endPointList)
	routeTable := &rtmgr.RouteTable{}
	e2TermEp := getEndpointByName(&endPointList, "E2TERM")
	if e2TermEp == nil {
		rtmgr.Logger.Error("Platform component not found: %v", "E2 Termination")
		rtmgr.Logger.Debug("Endpoints: %v", endPointList)
	}
	subManEp := getEndpointByName(&endPointList, "SUBMAN")
	if subManEp == nil {
		rtmgr.Logger.Error("Platform component not found: %v", "Subscription Manager")
		rtmgr.Logger.Debug("Endpoints: %v", endPointList)
	}
	e2ManEp := getEndpointByName(&endPointList, "E2MAN")
	if e2ManEp == nil {
		rtmgr.Logger.Error("Platform component not found: %v", "E2 Manager")
		rtmgr.Logger.Debug("Endpoints: %v", endPointList)
	}
	ueManEp := getEndpointByName(&endPointList, "UEMAN")
	if ueManEp == nil {
		rtmgr.Logger.Error("Platform component not found: %v", "UE Manger")
		rtmgr.Logger.Debug("Endpoints: %v", endPointList)
	}
	r.generatePlatformRoutes(e2TermEp, subManEp, e2ManEp, ueManEp, routeTable)

	for _, endPoint := range endPointList {
		rtmgr.Logger.Debug("Endpoint: %v, xAppType: %v", endPoint.Name, endPoint.XAppType)
		if endPoint.XAppType != sbi.PlatformType && len(endPoint.TxMessages) > 0 && len(endPoint.RxMessages) > 0 {
			r.generateXappRoutes(endPoint, e2TermEp, subManEp, routeTable)
			r.generateSubscriptionRoutes(e2TermEp, subManEp, routeTable)
		}
	}
	return routeTable
}
