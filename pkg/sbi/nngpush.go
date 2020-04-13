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
  Mnemonic:	nngpipe.go
  Abstract: mangos (NNG) Pipeline SBI implementation
  Date:		12 March 2019
*/

package sbi

/*
#include <rmr/rmr.h>
*/
import "C"




import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"routing-manager/pkg/rtmgr"
	"strconv"
	"time"
)

type EPStatus struct {
	endpoint string
	status   bool
}

type NngPush struct {
	Sbi
	NewSocket CreateNewNngSocketHandler
	rcChan    chan *xapp.RMRParams
}

type RMRParams struct {
	*xapp.RMRParams
}

func (params *RMRParams) String() string {
	var b bytes.Buffer
	sum := md5.Sum(params.Payload)
	fmt.Fprintf(&b, "params(Src=%s Mtype=%d SubId=%d Xid=%s Meid=%s Paylens=%d/%d Payhash=%x)", params.Src, params.Mtype, params.SubId, params.Xid, params.Meid.RanName, params.PayloadLen, len(params.Payload), sum)
	return b.String()
}

func NewNngPush() *NngPush {
	instance := new(NngPush)
	return instance
}

func (c *NngPush) Initialize(ip string) error {
	return nil
}

func (c *NngPush) Terminate() error {
	return nil
}

func (c *NngPush) AddEndpoint(ep *rtmgr.Endpoint) error {

	xapp.Logger.Debug("Invoked sbi.AddEndpoint")
	endpoint := ep.Ip + ":" + strconv.Itoa(DefaultNngPipelineSocketNumber)
	ep.Whid = int(xapp.Rmr.Openwh(endpoint))
	if ep.Whid < 0 {
		return errors.New("can't open warmhole connection for endpoint:" + ep.Uuid + " due to invalid Wormhole ID: " + string(ep.Whid))
	} else {
		xapp.Logger.Debug("Wormhole ID is %v and EP is %v", ep.Whid, endpoint)
	}

	return nil
}

func (c *NngPush) DeleteEndpoint(ep *rtmgr.Endpoint) error {
	xapp.Logger.Debug("Invoked sbi. DeleteEndpoint")
	xapp.Logger.Debug("args: %v", *ep)

	xapp.Rmr.Closewh(ep.Whid)
	return nil
}

func (c *NngPush) UpdateEndpoints(rcs *rtmgr.RicComponents) {
	c.updateEndpoints(rcs, c)
}

func (c *NngPush) DistributeAll(policies *[]string) error {
	xapp.Logger.Debug("Invoked: sbi.DistributeAll")
	xapp.Logger.Debug("args: %v", *policies)

	for _, ep := range rtmgr.Eps {
		go c.send(ep, policies)
	}

	return nil
}

func (c *NngPush) send(ep *rtmgr.Endpoint, policies *[]string) {
	xapp.Logger.Debug("Push policy to endpoint: " + ep.Uuid)

	for _, pe := range *policies {
		params := &RMRParams{&xapp.RMRParams{}}
		params.Mtype = 20
		params.PayloadLen = len([]byte(pe))
		params.Payload = []byte(pe)
		params.Mbuf = nil
		params.Whid = ep.Whid
		time.Sleep(1 * time.Millisecond)
		xapp.Rmr.SendMsg(params.RMRParams)
	}
	xapp.Logger.Info("NNG PUSH to endpoint " + ep.Uuid + ": OK (# of Entries:" + strconv.Itoa(len(*policies)) + ")")
}

func (c *NngPush) CreateEndpoint(payload string) *rtmgr.Endpoint {
	return c.createEndpoint(payload, c)
}

func (c *NngPush) DistributeToEp(policies *[]string, ep *rtmgr.Endpoint) error {
	xapp.Logger.Debug("Invoked: sbi.DistributeToEp")
	xapp.Logger.Debug("args: %v", *policies)

	go c.send(ep, policies)

	return nil
}

func (c *NngPush) DistributeRouteTables(route_table *[]string, meid_table *[]string) error {
	xapp.Logger.Debug("Invoked: sbi.DistributeRouteTables")
	xapp.Logger.Debug("args route_table: %v", route_table)
	xapp.Logger.Debug("args meid_table: %v", meid_table)

	channel := make(chan EPStatus)

	var i int = 2

	for _, ep := range rtmgr.Eps {
		go c.send_sync(ep, route_table, meid_table, channel, i)
		i = i + 1
	}

	count := 0
	result := make([]EPStatus, len(rtmgr.Eps))
	for i, _ := range result {
		result[i] = <-channel
		if result[i].status == true {
			count++
		} else {
			xapp.Logger.Error("RMR send is failed for endpoint %v", result[i].endpoint)
		}
	}

	if count < len(rtmgr.Eps) {
		return errors.New(" RMR response count " + string(count) + " is less than half of endpoint list " + string(len(rtmgr.Eps)))
	}

	return nil
}

func (c *NngPush) send_sync(ep *rtmgr.Endpoint, route_table *[]string, meidtable *[]string, channel chan EPStatus, call_id int) {
	xapp.Logger.Debug("Push policy to endpoint: " + ep.Uuid)

	ret := c.send_data(ep, route_table, call_id)

	if ret == true {
		ret = c.send_data(ep, meidtable, call_id)
	}
	channel <- EPStatus{ep.Uuid, ret}

}

/*

	1. first n-1 records rmr_wh_send (async send)
	2. last record rmr_wh_call (sync send)

*/

func (c *NngPush) send_data(ep *rtmgr.Endpoint, policies *[]string, call_id int) bool {
	xapp.Logger.Debug("sync send route data to endpoint: " + ep.Uuid + " call_id: " + string(call_id))
	var state int
	var retstr string

	length := len(*policies)

	for index, pe := range *policies {

		params := &RMRParams{&xapp.RMRParams{}}
		params.Mtype = 20
		params.PayloadLen = len([]byte(pe))
		params.Payload = []byte(pe)
		params.Mbuf = nil
		params.Whid = ep.Whid
		if index == length-1 {
			params.Callid = call_id
			params.Timeout = 200
			state, retstr = xapp.Rmr.SendCallMsg(params.RMRParams)
			if state != C.RMR_OK {
				xapp.Logger.Error("sync send route data to endpoint: " + ep.Uuid + " is failed,   call_id: " + string(call_id) + " for xapp.Rmr.SendCallMsg " + " return payload: " + retstr)
				return false
			} else {
				xapp.Logger.Info("sync send route data to endpoint: " + ep.Uuid + " is success,  call_id: " + string(call_id) + " return payload: " + retstr)
				return true
			}

		} else {
			if xapp.Rmr.SendMsg(params.RMRParams) != true {
				xapp.Logger.Error("sync send route data to endpoint: " + ep.Uuid + " is failed, call_id: " + string(call_id) + " for xapp.Rmr.SendMsg")
				return false
			}
		}
	}

	xapp.Logger.Error("sync send route data to endpoint: " + ep.Uuid + " is failed, call_id: " + string(call_id) + " xapp.Rmr.SendCallMsg is not called")
	return false
}
