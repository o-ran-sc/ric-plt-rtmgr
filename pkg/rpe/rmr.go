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
  Mnemonic:	rmr.go
  Abstract:	RMR Route Policy implementation
		Produces RMR (RIC Management Routing) formatted route messages
  Date:		16 March 2019
*/

package rpe

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"routing-manager/pkg/models"
	"routing-manager/pkg/rtmgr"
	"strconv"
	"strings"
)

type Rmr struct {
	Rpe
}

type RmrPush struct {
	Rmr
}

func NewRmrPush() *RmrPush {
	instance := new(RmrPush)
	return instance
}

/*
Produces the raw route message consumable by RMR
*/
func (r *Rmr) generateRMRPolicies(eps rtmgr.Endpoints, rcs *rtmgr.RicComponents, key string) *[]string {
	rawrt := []string{key + "newrt|start\n"}
	rt := r.generateRouteTable(eps)
	for _, rte := range *rt {
		rawrte := key + "mse|" + rte.MessageType
		for _, tx := range rte.TxList {
			rawrte += "," + tx.Ip + ":" + strconv.Itoa(int(tx.Port))
		}
		rawrte += "|" + strconv.Itoa(int(rte.SubID)) + "|"
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

		if rte.RouteType == "%meid" {
			rawrte += "%" + "meid"
		} else {
			rawrte += group
		}

		rawrt = append(rawrt, rawrte+"\n")
	}
	for _, val := range rtmgr.DynamicRouteList {
		rawrt = append(rawrt, val)
	}

	rawrt = append(rawrt, key+"newrt|end\n")
	count := 0

	rawrt = append(rawrt, key+"meid_map|start\n")
	
	keys := make(map[string]MeidEntry)
	MEID := ""
	E2TIP := ""
	RECTYP := ""
	for _, value := range rcs.MeidMap {
		if strings.Contains(value, "mme_ar") {
			tmpstr := strings.Split(value, "|")
			//MEID entry for mme_ar must always contain 3 strings speartred by | i.e "mme_ar|<string1>|<string2>"
			MEID = strings.TrimSuffix(tmpstr[2], "\n")
			E2TIP = strings.TrimSuffix(tmpstr[1], "\n")
			RECTYP = "mme_ar"
		} else if strings.Contains(value, "mme_del") {
			tmpstr := strings.Split(value, "|")
			MEID = strings.TrimSuffix(tmpstr[1], "\n")
			E2TIP = ""
			RECTYP = "mme_del"
		}
		keys[MEID] = MeidEntry{RECTYP, E2TIP}
	}

	for k, v := range keys {
		if v.recordtype == "mme_ar" {
			rawrt = append(rawrt, key+v.recordtype+"|"+v.e2tip+"|"+k+"\n")
		}
		count++
	}

	rawrt = removeEmptyStrings(rawrt)
	rawrt = append(rawrt, key+"meid_map|end|"+strconv.Itoa(count)+"\n")

	xapp.Logger.Debug("rmr.GeneratePolicies returns: %v", rawrt)
	xapp.Logger.Debug("rmr.GeneratePolicies returns: %v", rcs)
	return &rawrt
}

/*
Produces the raw route message consumable by RMR
*/
func (r *Rmr) generatePartialRMRPolicies(eps rtmgr.Endpoints, xappSubData *models.XappSubscriptionData, key string, updatetype rtmgr.RMRUpdateType) *[]string {
	rawrt := []string{key + "updatert|start\n"}
	rt := r.generatePartialRouteTable(eps, xappSubData, updatetype)
	for _, rte := range *rt {
		rawrte := key + "mse|" + rte.MessageType
		for _, tx := range rte.TxList {
			rawrte += "," + tx.Ip + ":" + strconv.Itoa(int(tx.Port))
		}
		rawrte += "|" + strconv.Itoa(int(rte.SubID)) + "|"
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

		if rte.RouteType == "%meid" {
			rawrte += "%" + "meid"
		} else {
			rawrte += group
		}

		rawrt = append(rawrt, rawrte+"\n")
	}

	rawrt = append(rawrt, key+"updatert|end\n")
	//count := 0

	xapp.Logger.Debug("rmr.GeneratePolicies returns: %v", rawrt)
	return &rawrt
}
func (r *RmrPush) GeneratePolicies(eps rtmgr.Endpoints, rcs *rtmgr.RicComponents) *[]string {
	xapp.Logger.Debug("Invoked rmr.GeneratePolicies, args: %v: ", eps)
	return r.generateRMRPolicies(eps, rcs, "")
}

func (r *RmrPush) GenerateRouteTable(eps rtmgr.Endpoints) *rtmgr.RouteTable {
	return r.generateRouteTable(eps)
}

func (r *RmrPush) GeneratePartialPolicies(eps rtmgr.Endpoints, xappSubData *models.XappSubscriptionData, updatetype rtmgr.RMRUpdateType) *[]string {
	xapp.Logger.Debug("Invoked rmr.GeneratePartialPolicies, args: %v: ", eps)
	return r.generatePartialRMRPolicies(eps, xappSubData, "", updatetype)
}

func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
