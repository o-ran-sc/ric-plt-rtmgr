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
  Mnemonic:	rmr.go
  Abstract:	RMR Route Policy implementation
		Produces RMR (RIC Management Routing) formatted route messages
  Date:		16 March 2019
*/

package rpe

import (
	"rtmgr"
	"strconv"
)

/*
Produces the raw route message consumable by RMR
*/
func generateRMRPolicies(eps rtmgr.Endpoints, key string) *[]string {
	rtmgr.Logger.Debug("Invoked rmr.generateRMRPolicies")
	rtmgr.Logger.Debug("args: %v", eps)
	rawrt := []string{key + "newrt|start\n"}
	rt := getRouteTable(eps)
	for _, rte := range *rt {
		rawrte := key + "rte|" + rte.MessageType
		for _, tx := range rte.TxList {
			rawrte += "," + tx.Ip + ":" + strconv.Itoa(int(tx.Port))
		}
		rawrte += "|"
		group := ""
		for _, rxg := range rte.RxGroups {
			member := ""
			for _, rx := range rxg {
				if member == "" {
					member += rx.Ip + ":" + strconv.Itoa(int(rx.Port))
				} else {
					member += "," + rx.Ip + ":" + strconv.Itoa(int(rx.Port))
				}
			}
			if group == "" {
				group += member
			} else {
				group += ";" + member
			}
		}
		rawrte += group
		rawrt = append(rawrt, rawrte+"\n")
	}
	rawrt = append(rawrt, key+"newrt|end\n")
	rtmgr.Logger.Debug("rmr.generateRMRPolicies returns: %v", rawrt)
	return &rawrt
}

func generateRMRPubPolicies(eps rtmgr.Endpoints) *[]string {
	return generateRMRPolicies(eps, "00000           ")
}

func generateRMRPushPolicies(eps rtmgr.Endpoints) *[]string {
	return generateRMRPolicies(eps, "")
}
