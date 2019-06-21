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
	"strconv"
)

var (
	SupportedRpes = []*RpeEngineConfig{
		&RpeEngineConfig{
			Name:        "rmrpub",
			Version:     "pubsub",
			Protocol:    "rmruta",
			Instance:    NewRmrPub(),
			IsAvailable: true,
		},
		&RpeEngineConfig{
			Name:        "rmrpush",
			Version:     "pubsush",
			Protocol:    "rmruta",
			Instance:    NewRmrPush(),
			IsAvailable: true,
		},
	}
)

func GetRpe(rpeName string) (RpeEngine, error) {
	for _, rpe := range SupportedRpes {
		if rpe.Name == rpeName && rpe.IsAvailable {
			return rpe.Instance, nil
		}
	}
	return nil, errors.New("SBI:" + rpeName + " is not supported or still not a available")
}

type Rpe struct {
}

/*
Gets the raw xApp list and generates the list of sender endpoints and receiver endpoint groups
Returns the Tx EndpointList map where the key is the messge type and also returns the nested map of Rx EndpointList's map where keys are message type and xapp type
Endpoint object's message type already transcoded to integer id
*/

func (r *Rpe) getRouteRxTxLists(eps rtmgr.Endpoints) (*map[string]rtmgr.EndpointList, *map[string]map[string]rtmgr.EndpointList) {
	txlist := make(map[string]rtmgr.EndpointList)
	rxgroups := make(map[string]map[string]rtmgr.EndpointList)
	for _, ep := range eps {
		for _, message := range ep.RxMessages {
			messageid := rtmgr.MESSAGETYPES[message]
			if _, ok := rxgroups[messageid]; !ok {
				rxgroups[messageid] = make(map[string]rtmgr.EndpointList)
			}
			rxgroups[messageid][ep.XAppType] = append(rxgroups[messageid][ep.XAppType], (*ep))
		}
		for _, message := range ep.TxMessages {
			messageid := rtmgr.MESSAGETYPES[message]
			txlist[messageid] = append(txlist[messageid], (*ep))
		}
	}
	return &txlist, &rxgroups
}

/*
Gets the raw xapp list and creates a route table for
Returns the array of route table entries
*/
func (r *Rpe) getRouteTable(eps rtmgr.Endpoints) *rtmgr.RouteTable {
	tx, rx := r.getRouteRxTxLists(eps)
	var rt rtmgr.RouteTable
	for _, messagetype := range rtmgr.MESSAGETYPES {
		if _, ok := (*tx)[messagetype]; !ok {
			continue
		}
		if _, ok := (*rx)[messagetype]; !ok {
			continue
		}
		var rxgroups []rtmgr.EndpointList
		for _, endpointlist := range (*rx)[messagetype] {
			rxgroups = append(rxgroups, endpointlist)
		}
		rte := rtmgr.RouteTableEntry{
			messagetype,
			(*tx)[messagetype],
			rxgroups,
			-1,
		}
		rt = append(rt, rte)
	}
	r.addStaticRoutes(eps, &rt)
	r.addSubscriptionRoutes(eps, &rt, &rtmgr.Subs)
	return &rt
}

/*
Adds specific static routes to the route table
which cannot be calculated with endpoint tx/rx message types.
*/
func (r *Rpe) addStaticRoutes(eps rtmgr.Endpoints, rt *rtmgr.RouteTable) {
	var uemanep, submanep *rtmgr.Endpoint
	for _, ep := range eps {
		if ep.Name == "UEMAN" {
			uemanep = ep
		}
		if ep.Name == "SUBMAN" {
			submanep = ep
		}
	}

	if uemanep != nil && submanep != nil {
		txlist := rtmgr.EndpointList{*uemanep}
		rxlist := []rtmgr.EndpointList{[]rtmgr.Endpoint{*submanep}}
		rte1 := rtmgr.RouteTableEntry{
			rtmgr.MESSAGETYPES["RIC_SUB_REQ"],
			txlist,
			rxlist,
			-1,
		}
		rte2 := rtmgr.RouteTableEntry{
			rtmgr.MESSAGETYPES["RIC_SUB_DEL_REQ"],
			txlist,
			rxlist,
			-1,
		}
		*rt = append(*rt, rte1)
		*rt = append(*rt, rte2)
	} else {
		rtmgr.Logger.Warn("Cannot get the static route details of the platform components UEMAN/SUBMAN")
	}
}

func getEndpointByName(eps *rtmgr.Endpoints, name string) *rtmgr.Endpoint {
	for _, ep := range *eps {
		if ep.Name == name {
			rtmgr.Logger.Debug("name: %s", ep.Name)
			rtmgr.Logger.Debug("ep: &v",ep)
			return ep
		}
	}
	return nil
}

func getEndpointByUuid(eps *rtmgr.Endpoints, uuid string) *rtmgr.Endpoint {
	for _, ep := range *eps {
		if ep.Uuid == uuid {
			rtmgr.Logger.Debug("name: %s", ep.Uuid)
			rtmgr.Logger.Debug("ep: &v",ep)
			return ep
		}
	}
	return nil
}
func (r *Rpe) addSubscriptionRoutes(eps rtmgr.Endpoints, rt *rtmgr.RouteTable, subs *rtmgr.SubscriptionList) {
	rtmgr.Logger.Debug("rpe.addSubscriptionRoutes invoked")
	rtmgr.Logger.Debug("params: %v", eps)
	var e2termep, submanep, xappEp *rtmgr.Endpoint
	var xappName string
	e2termep = getEndpointByName(&eps, "E2TERM")
	submanep = getEndpointByName(&eps, "SUBMAN")
	if e2termep != nil && submanep != nil {
		// looping through the subscription list, add routes one by one
		for _, sub := range *subs {
			// SubMan -> XApp
			xappName = sub.Fqdn + ":" + strconv.Itoa(int(sub.Port))
			xappEp = getEndpointByUuid(&eps, xappName)
			if xappEp == nil {
				rtmgr.Logger.Error("XApp not found: %s", xappName)
				rtmgr.Logger.Debug("Endpoints: %v", eps)
			} else {
				txlist := rtmgr.EndpointList{*submanep}
				rxlist := []rtmgr.EndpointList{[]rtmgr.Endpoint{*xappEp}}
				subManMsgs := []string{"RIC_SUB_RESP", "RIC_SUB_FAILURE", "RIC_SUB_DEL_RESP", "RIC_SUB_DEL_FAILURE"}
				for _, entry := range subManMsgs {
					rte := rtmgr.RouteTableEntry{
						rtmgr.MESSAGETYPES[entry],
						txlist,
						rxlist,
						sub.SubID,
					}
					*rt = append(*rt, rte)
				}
				// E2Term -> XApp
				txlist = rtmgr.EndpointList{*e2termep}
				rxlist = []rtmgr.EndpointList{[]rtmgr.Endpoint{*xappEp}}
				e2apMsgs := []string{"RIC_CONTROL_ACK", "RIC_CONTROL_FAILURE", "RIC_INDICATION"}
				for _, entry := range e2apMsgs {
					rte := rtmgr.RouteTableEntry{
						rtmgr.MESSAGETYPES[entry],
						txlist,
						rxlist,
						sub.SubID,
					}
					*rt = append(*rt, rte)
				}
			}
		}
		rtmgr.Logger.Debug("addSubscriptionRoutes eps: %v", eps)
	} else {
		rtmgr.Logger.Warn("Subscription route update failure: Cannot get the static route details of the platform components E2TERM/SUBMAN")
	}

}
