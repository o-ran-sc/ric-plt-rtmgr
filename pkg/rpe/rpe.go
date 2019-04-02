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
	"fmt"
	"rtmgr"
	"strconv"
)

var (
	SupportedRpes = []*RpeEngineConfig{
		&RpeEngineConfig{
			RpeEngine{
				Name:     "rmr",
				Version:  "v1",
				Protocol: "rmruta",
			},
			generatePolicies(generateRMRPolicies),
			true,
		},
	}
)

func ListRpes() {
	fmt.Printf("RPE:\n")
	for _, rpe := range SupportedRpes {
		if rpe.IsAvailable {
			rtmgr.Logger.Info(rpe.Engine.Name + "/" + rpe.Engine.Version)
		}
	}
}

func GetRpe(rpeName string) (*RpeEngineConfig, error) {
	for _, rpe := range SupportedRpes {
		if rpe.Engine.Name == rpeName && rpe.IsAvailable {
			return rpe, nil
		}
	}
	return nil, errors.New("SBI:" + rpeName + "is not supported or still not a available")
}

/*
Gets the raw xApp list and generates the list of sender endpoints and receiver endpoint groups
Returns the Tx EndpointList map where the key is the messge type and also returns the nested map of Rx EndpointList's map where keys are message type and xapp type
Endpoint object's message type already transcoded to integer id
*/
func getEndpointLists(xapps *[]rtmgr.XApp) (*map[string]rtmgr.EndpointList, *map[string]map[string]rtmgr.EndpointList) {
	txlist := make(map[string]rtmgr.EndpointList)
	rxgroups := make(map[string]map[string]rtmgr.EndpointList)
	for _, xapp := range *xapps {
		for _, instance := range xapp.Instances {
			ep := rtmgr.Endpoint{
				instance.Name,
				xapp.Name,
				instance.Ip + ":" + strconv.Itoa(instance.Port),
			}
			for _, message := range instance.RxMessages {
				messageid := rtmgr.MESSAGETYPES[message]
				if _, ok := rxgroups[messageid]; !ok {
					rxgroups[messageid] = make(map[string]rtmgr.EndpointList)
				}
				rxgroups[messageid][xapp.Name] = append(rxgroups[messageid][xapp.Name], ep)
			}
			for _, message := range instance.TxMessages {
				messageid := rtmgr.MESSAGETYPES[message]
				txlist[messageid] = append(txlist[messageid], ep)
			}
		}
	}
	return &txlist, &rxgroups
}

/*
Gets the raw xapp list and creates a route table for
Returns the array of route table entries
*/
func getRouteTable(xapps *[]rtmgr.XApp) *rtmgr.RouteTable {
	tx, rx := getEndpointLists(xapps)
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
		}
		rt = append(rt, rte)
	}
	return &rt
}
