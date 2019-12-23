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
	        if (tx != nil) {
	                txList = rtmgr.EndpointList{*tx}
	        }
	        if (rx != nil) {
	                rxList = []rtmgr.EndpointList{[]rtmgr.Endpoint{*rx}}
	        }
	        messageId := rtmgr.MessageTypes[messageType]
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

	if (tx != nil) {
	        txList = rtmgr.EndpointList{*tx}
	}

	if (rx != nil) {
	        rxList = []rtmgr.EndpointList{rx}
	}

	messageId := rtmgr.MessageTypes[messageType]
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
	r.addRoute("RIC_SUB_REQ", subManEp, nil, routeTable, -1, "%meid")
	r.addRoute("RIC_SUB_DEL_REQ", subManEp, nil, routeTable, -1, "%meid")
	//E2 Termination -> Subscription Manager
        r.addRoute("RIC_SUB_RESP", nil, subManEp, routeTable, -1, "")
        r.addRoute("RIC_SUB_DEL_RESP", nil, subManEp, routeTable, -1, "")
        r.addRoute("RIC_SUB_FAILURE", nil, subManEp, routeTable, -1, "")
        r.addRoute("RIC_SUB_DEL_FAILURE", nil, subManEp, routeTable, -1, "")

	//TODO: UE Man Routes removed (since it is not existing)
	//UE Manager -> Subscription Manager
	//r.addRoute("RIC_SUB_REQ", ueManEp, subManEp, routeTable)
	//r.addRoute("RIC_SUB_DEL_REQ", ueManEp, subManEp, routeTable)
	////UE Manager -> E2 Termination
	//r.addRoute("RIC_CONTROL_REQ", ueManEp, e2TermEp, routeTable)

	//Platform Routes --- X2 Routes
	//E2 Manager -> E2 Termination
        r.addRoute("RIC_X2_SETUP_REQ", e2ManEp, nil, routeTable, -1, "%meid")
        r.addRoute("RIC_X2_RESET_REQ", e2ManEp, nil, routeTable, -1, "%meid")
        r.addRoute("RIC_X2_RESET_RESP", e2ManEp, nil, routeTable, -1, "%meid")
        r.addRoute("RIC_ENDC_X2_SETUP_REQ", e2ManEp, nil, routeTable, -1, "%meid")
        r.addRoute("RIC_ENB_CONF_UPDATE_ACK", e2ManEp, nil, routeTable, -1, "%meid")
        r.addRoute("RIC_ENB_CONF_UPDATE_FAILURE", e2ManEp, nil, routeTable, -1, "%meid")
        r.addRoute("RIC_ENDC_CONF_UPDATE_ACK", e2ManEp, nil, routeTable, -1, "%meid")
        r.addRoute("RIC_ENDC_CONF_UPDATE_FAILURE", e2ManEp, nil, routeTable, -1, "%meid")

        if len(e2TermEp) > 0 {
                r.addRoute_rx_list("RIC_SCTP_CLEAR_ALL", e2ManEp, e2TermEp, routeTable, -1, "")
                r.addRoute_rx_list("E2_TERM_KEEP_ALIVE_REQ", e2ManEp, e2TermEp, routeTable, -1, "")
        }

	//E2 Termination -> E2 Manager
        r.addRoute("E2_TERM_INIT", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_X2_SETUP_RESP", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_X2_SETUP_FAILURE", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_X2_RESET_REQ", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_X2_RESET_RESP", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_ENDC_X2_SETUP_RESP", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_ENDC_X2_SETUP_FAILURE", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_ENDC_CONF_UPDATE", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_SCTP_CONNECTION_FAILURE", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_ERROR_INDICATION", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_ENB_CONF_UPDATE", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("RIC_ENB_LOAD_INFORMATION", nil, e2ManEp, routeTable, -1, "")
        r.addRoute("E2_TERM_KEEP_ALIVE_RESP", nil, e2ManEp, routeTable, -1, "")



	//E2 Manager -> Resource Status Manager
        r.addRoute("RAN_CONNECTED", e2ManEp, rsmEp, routeTable, -1, "")
        r.addRoute("RAN_RESTARTED", e2ManEp, rsmEp, routeTable, -1, "")
        r.addRoute("RAN_RECONFIGURED", e2ManEp, rsmEp, routeTable, -1, "")

	//Resource Status Manager -> E2 Termination
	r.addRoute("RIC_RES_STATUS_REQ", rsmEp, nil, routeTable, -1, "%meid")
	//E2 Termination -> Resource Status Manager
        r.addRoute("RIC_RES_STATUS_RESP", nil, rsmEp, routeTable, -1, "")
        r.addRoute("RIC_RES_STATUS_FAILURE", nil, rsmEp, routeTable, -1, "")

	//ACxapp -> A1 Mediator
	r.addRoute("A1_POLICY_QUERY", nil, a1mediatorEp, routeTable, -1, "")
	r.addRoute("A1_POLICY_RESPONSE", nil, a1mediatorEp, routeTable, -1, "")
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
