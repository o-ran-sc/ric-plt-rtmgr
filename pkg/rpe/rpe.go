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
  Mnemonic:	rpe.go
  Abstract:	Contains RPE (Route Policy Engine) module definitions and generic RPE components
  Date:		16 March 2019
*/

package rpe

import (
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sbi"
	"runtime"
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
			xapp.Logger.Debug("name: %s", ep.Name)
			xapp.Logger.Debug("ep: %v", ep)
			return ep
		}
	}
	return nil
}

func getEndpointListByName(eps *rtmgr.Endpoints, name string) []rtmgr.Endpoint {
	var eplist []rtmgr.Endpoint

	for _, ep := range *eps {
		if ep.Name == name {
			xapp.Logger.Debug("name: %s", ep.Name)
			xapp.Logger.Debug("ep: %v", ep)
			eplist = append(eplist, *ep)
		}
	}
	return eplist
}

func getEndpointByUuid(uuid string) *rtmgr.Endpoint {
	endPoints := rtmgr.Eps
	for _, ep := range endPoints {
		if ep.Uuid == uuid {
			xapp.Logger.Debug("name: %s", ep.Uuid)
			xapp.Logger.Debug("ep: %v", ep)
			return ep
		}
	}
	return nil
}

func (r *Rpe) addRoute(messageType string, tx *rtmgr.Endpoint, rx *rtmgr.Endpoint, routeTable *rtmgr.RouteTable, subId int32, routeType string) {
	txList := rtmgr.EndpointList{}
	rxList := []rtmgr.EndpointList{}

	if tx == nil && rx == nil {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			xapp.Logger.Error("Route addition skipped: Either TX or RX endpoint not present. Caller function is %s", details.Name())
		}
	} else {
		if tx != nil {
			txList = rtmgr.EndpointList{*tx}
		}
		if rx != nil {
			rxList = []rtmgr.EndpointList{[]rtmgr.Endpoint{*rx}}
		}
		messageId := strconv.Itoa(xapp.RICMessageTypes[messageType])
		route := rtmgr.RouteTableEntry{
			MessageType: messageId,
			TxList:      txList,
			RxGroups:    rxList,
			SubID:       subId,
			RouteType:   routeType}
		*routeTable = append(*routeTable, route)
		//	        xapp.Logger.Debug("Route added: MessageTyp: %v, Tx: %v, Rx: %v, SubId: %v", messageId, tx.Uuid, rx.Uuid, subId)
		//	        xapp.Logger.Trace("Route added: MessageTyp: %v, Tx: %v, Rx: %v, SubId: %v", messageId, tx, rx, subId)
	}
}

func (r *Rpe) addRoute_rx_list(messageType string, tx *rtmgr.Endpoint, rx []rtmgr.Endpoint, routeTable *rtmgr.RouteTable, subId int32, routeType string) {
	txList := rtmgr.EndpointList{}
	rxList := []rtmgr.EndpointList{}

	if tx != nil {
		txList = rtmgr.EndpointList{*tx}
	}

	if rx != nil {
		rxList = []rtmgr.EndpointList{rx}
	}

	messageId := strconv.Itoa(xapp.RICMessageTypes[messageType])
	route := rtmgr.RouteTableEntry{
		MessageType: messageId,
		TxList:      txList,
		RxGroups:    rxList,
		SubID:       subId,
		RouteType:   routeType}
	*routeTable = append(*routeTable, route)
	//	xapp.Logger.Debug("Route added: MessageTyp: %v, Tx: %v, Rx: %v, SubId: %v", messageId, tx.Uuid, rx.Uuid, subId)
	//	xapp.Logger.Trace("Route added: MessageTyp: %v, Tx: %v, Rx: %v, SubId: %v", messageId, tx, rx, subId)
}

func (r *Rpe) generateXappRoutes(xAppEp *rtmgr.Endpoint, e2TermEp *rtmgr.Endpoint, subManEp *rtmgr.Endpoint, routeTable *rtmgr.RouteTable) {
	xapp.Logger.Debug("rpe.generateXappRoutes invoked")
	xapp.Logger.Debug("Endpoint: %v, xAppType: %v", xAppEp.Name, xAppEp.XAppType)
	if xAppEp.XAppType != sbi.PlatformType && (len(xAppEp.TxMessages) > 0 || len(xAppEp.RxMessages) > 0) {
		/// TODO ---
		//xApp -> Subscription Manager
		r.addRoute("RIC_SUB_REQ", xAppEp, subManEp, routeTable, -1, "")
		r.addRoute("RIC_SUB_DEL_REQ", xAppEp, subManEp, routeTable, -1, "")
		//xApp -> E2 Termination
		//		r.addRoute("RIC_CONTROL_REQ", xAppEp, e2TermEp, routeTable, -1, "")
		r.addRoute("RIC_CONTROL_REQ", xAppEp, nil, routeTable, -1, "%meid")
		//E2 Termination -> xApp
		///		r.addRoute("RIC_CONTROL_ACK", e2TermEp, xAppEp, routeTable, -1, "")
		///		r.addRoute("RIC_CONTROL_FAILURE", e2TermEp, xAppEp, routeTable, -1, "")
		r.addRoute("RIC_CONTROL_ACK", nil, xAppEp, routeTable, -1, "")
		r.addRoute("RIC_CONTROL_FAILURE", nil, xAppEp, routeTable, -1, "")
	}
	//xApp->A1Mediator
	if xAppEp.XAppType != sbi.PlatformType && len(xAppEp.Policies) > 0 {
		xapp.Logger.Debug("rpe.generateXappRoutes found policies section")
		for _, policy := range xAppEp.Policies {
			r.addRoute("A1_POLICY_REQ", nil, xAppEp, routeTable, policy, "")
		}
	}

}

func (r *Rpe) generateSubscriptionRoutes(selectedxAppEp *rtmgr.Endpoint, e2TermEp *rtmgr.Endpoint, subManEp *rtmgr.Endpoint, routeTable *rtmgr.RouteTable) {
	xapp.Logger.Debug("rpe.addSubscriptionRoutes invoked")
	subscriptionList := &rtmgr.Subs
	for _, subscription := range *subscriptionList {
		xapp.Logger.Debug("Subscription: %v", subscription)
		xAppUuid := subscription.Fqdn + ":" + strconv.Itoa(int(subscription.Port))
		xapp.Logger.Debug("xApp UUID: %v", xAppUuid)
		xAppEp := getEndpointByUuid(xAppUuid)
		if xAppEp.Uuid == selectedxAppEp.Uuid {
			xapp.Logger.Debug("xApp UUID is matched for selected xApp.UUID: %v and xApp.Name: %v", selectedxAppEp.Uuid, selectedxAppEp.Name)
			/// TODO
			//Subscription Manager -> xApp
			r.addRoute("RIC_SUB_RESP", subManEp, xAppEp, routeTable, subscription.SubID, "")
			r.addRoute("RIC_SUB_FAILURE", subManEp, xAppEp, routeTable, subscription.SubID, "")
			r.addRoute("RIC_SUB_DEL_RESP", subManEp, xAppEp, routeTable, subscription.SubID, "")
			r.addRoute("RIC_SUB_DEL_FAILURE", subManEp, xAppEp, routeTable, subscription.SubID, "")
			//E2 Termination -> xApp
			r.addRoute("RIC_INDICATION", e2TermEp, xAppEp, routeTable, subscription.SubID, "")
			r.addRoute("RIC_CONTROL_ACK", e2TermEp, xAppEp, routeTable, subscription.SubID, "")
			r.addRoute("RIC_CONTROL_FAILURE", e2TermEp, xAppEp, routeTable, subscription.SubID, "")
		}
	}
}

func (r *Rpe) generatePlatformRoutes(e2TermEp []rtmgr.Endpoint, subManEp *rtmgr.Endpoint, e2ManEp *rtmgr.Endpoint, ueManEp *rtmgr.Endpoint, rsmEp *rtmgr.Endpoint, a1mediatorEp *rtmgr.Endpoint, routeTable *rtmgr.RouteTable) {
	xapp.Logger.Debug("rpe.generatePlatformRoutes invoked")
	//Platform Routes --- Subscription Routes
	//Subscription Manager -> E2 Termination
	for _, routes := range *rtmgr.PrsCfg {
		var sendEp *rtmgr.Endpoint
		var Ep *rtmgr.Endpoint
		switch routes.SenderEndPoint {
		case "SUBMAN":
			sendEp = subManEp
		case "E2MAN":
			sendEp = e2ManEp
		case "UEMAN":
			sendEp = ueManEp
		case "RSM":
			sendEp = rsmEp
		case "A1MEDIATOR":
			sendEp = a1mediatorEp
		}
		switch routes.EndPoint {
		case "SUBMAN":
			Ep = subManEp
		case "E2MAN":
			Ep = e2ManEp
		case "UEMAN":
			Ep = ueManEp
		case "RSM":
			Ep = rsmEp
		case "A1MEDIATOR":
			Ep = a1mediatorEp
		}
		if routes.EndPoint == "E2TERMINST" && len(e2TermEp) > 0 {
			r.addRoute_rx_list(routes.MessageType, sendEp, e2TermEp, routeTable, routes.SubscriptionId, routes.Meid)
			continue
		}
		r.addRoute(routes.MessageType, sendEp, Ep, routeTable, routes.SubscriptionId, routes.Meid)
	}
}

func (r *Rpe) generateRouteTable(endPointList rtmgr.Endpoints) *rtmgr.RouteTable {
	xapp.Logger.Debug("rpe.generateRouteTable invoked")
	xapp.Logger.Debug("Endpoint List:  %v", endPointList)
	routeTable := &rtmgr.RouteTable{}
	e2TermEp := getEndpointByName(&endPointList, "E2TERM")
	if e2TermEp == nil {
		xapp.Logger.Error("Platform component not found: %v", "E2 Termination")
		xapp.Logger.Debug("Endpoints: %v", endPointList)
	}
	subManEp := getEndpointByName(&endPointList, "SUBMAN")
	if subManEp == nil {
		xapp.Logger.Error("Platform component not found: %v", "Subscription Manager")
		xapp.Logger.Debug("Endpoints: %v", endPointList)
	}
	e2ManEp := getEndpointByName(&endPointList, "E2MAN")
	if e2ManEp == nil {
		xapp.Logger.Error("Platform component not found: %v", "E2 Manager")
		xapp.Logger.Debug("Endpoints: %v", endPointList)
	}
	ueManEp := getEndpointByName(&endPointList, "UEMAN")
	if ueManEp == nil {
		xapp.Logger.Error("Platform component not found: %v", "UE Manger")
		xapp.Logger.Debug("Endpoints: %v", endPointList)
	}
	rsmEp := getEndpointByName(&endPointList, "RSM")
	if rsmEp == nil {
		xapp.Logger.Error("Platform component not found: %v", "Resource Status Manager")
		xapp.Logger.Debug("Endpoints: %v", endPointList)
	}
	A1MediatorEp := getEndpointByName(&endPointList, "A1MEDIATOR")
	if A1MediatorEp == nil {
		xapp.Logger.Error("Platform component not found: %v", "A1Mediator")
		xapp.Logger.Debug("Endpoints: %v", endPointList)
	}

	e2TermListEp := getEndpointListByName(&endPointList, "E2TERMINST")
	if len(e2TermListEp) == 0 {
		xapp.Logger.Error("Platform component not found: %v", "E2 Termination List")
		xapp.Logger.Debug("Endpoints: %v", endPointList)
	}
	r.generatePlatformRoutes(e2TermListEp, subManEp, e2ManEp, ueManEp, rsmEp, A1MediatorEp, routeTable)

	for _, endPoint := range endPointList {
		xapp.Logger.Debug("Endpoint: %v, xAppType: %v", endPoint.Name, endPoint.XAppType)
		if endPoint.XAppType != sbi.PlatformType && (len(endPoint.TxMessages) > 0 || len(endPoint.RxMessages) > 0) {
			r.generateXappRoutes(endPoint, e2TermEp, subManEp, routeTable)
			r.generateSubscriptionRoutes(endPoint, e2TermEp, subManEp, routeTable)
		}
	}
	return routeTable
}
