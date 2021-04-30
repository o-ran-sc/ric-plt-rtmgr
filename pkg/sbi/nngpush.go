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
  Mnemonic:	rmrpipe.go
  Abstract: mangos (RMR) Pipeline SBI implementation
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
	"strings"
	"sync"
	"time"
)

var rmrcallid = 1
var rmrdynamiccallid = 201
var addendpointct = 1

var conn sync.Mutex

type RmrPush struct {
	Sbi
	rcChan chan *xapp.RMRParams
}

type EPStatus struct {
	endpoint string
	status   bool
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

func NewRmrPush() *RmrPush {
	instance := new(RmrPush)
	return instance
}

func (c *RmrPush) Initialize(ip string) error {
	return nil
}

func (c *RmrPush) Terminate() error {
	return nil
}

func (c *RmrPush) AddEndpoint(ep *rtmgr.Endpoint) error {
	count := addendpointct + 1
	xapp.Logger.Debug("Invoked sbi.AddEndpoint for %s with count = %d", ep.Ip, count)
	endpoint := ep.Ip + ":" + strconv.Itoa(DefaultRmrPipelineSocketNumber)
	ep.Whid = int(xapp.Rmr.Openwh(endpoint))
	if ep.Whid < 0 {
		time.Sleep(time.Duration(10) * time.Second)
		ep.Whid = int(xapp.Rmr.Openwh(endpoint))
		if ep.Whid < 0 {
			return errors.New("can't open warmhole connection for endpoint:" + ep.Uuid + " due to invalid Wormhole ID: " + string(ep.Whid) + " count: " + strconv.Itoa(count))
		}
	} else {
		xapp.Logger.Debug("Wormhole ID is %v and EP is %v", ep.Whid, endpoint)
	}

	return nil
}

func (c *RmrPush) DeleteEndpoint(ep *rtmgr.Endpoint) error {
	xapp.Logger.Debug("Invoked sbi. DeleteEndpoint")
	xapp.Logger.Debug("args: %v", *ep)

	xapp.Rmr.Closewh(ep.Whid)
	return nil
}

func (c *RmrPush) UpdateEndpoints(rcs *rtmgr.RicComponents) {
	c.updateEndpoints(rcs, c)
}

func (c *RmrPush) DistributeAll(policies *[]string) error {
	xapp.Logger.Debug("Invoked: sbi.DistributeAll")
	xapp.Logger.Debug("args: %v", *policies)

	/*for _, ep := range rtmgr.Eps {
		go c.send(ep, policies)
	}*/
	//channel := make(chan EPStatus)

	if rmrcallid == 200 {
		rmrcallid = 1
	}

	for _, ep := range rtmgr.Eps {
		go c.send_sync(ep, policies, rmrcallid)
	}

	rmrcallid++

	/*
				count := 0
		        result := make([]EPStatus, len(rtmgr.Eps))
		        for i, _ := range result {
		                result[i] = <-channel
		                if result[i].status == true {
		                        count++
		                } else {
		                        xapp.Logger.Error("RMR send failed for endpoint %v", result[i].endpoint)
		                }
		        }

		        if count < len(rtmgr.Eps) {
		                return errors.New(" RMR response count " + string(count) + " is less than half of endpoint list " + string(len(rtmgr.Eps)))
		        }*/

	return nil
}

//func (c *RmrPush) send_sync(ep *rtmgr.Endpoint, policies *[]string, channel chan EPStatus, call_id int) {
func (c *RmrPush) send_sync(ep *rtmgr.Endpoint, policies *[]string, call_id int) {
	xapp.Logger.Debug("Push policy to endpoint: " + ep.Uuid)

	ret := c.send_data(ep, policies, call_id)
	xapp.Logger.Debug("return value is %v", ret)
	conn.Lock()
	rtmgr.RMRConnStatus[ep.Uuid] = ret
	conn.Unlock()
	// Handling per connection .. may be updating global map

	//channel <- EPStatus{ep.Uuid, ret}

}

func (c *RmrPush) send_data(ep *rtmgr.Endpoint, policies *[]string, call_id int) bool {
	xapp.Logger.Debug("Invoked send_data to endpoint: " + ep.Uuid + " call_id: " + strconv.Itoa(call_id))
	var state int
	var retstr string

	var policy = []byte{}

	for _, pe := range *policies {
		b := []byte(pe)
		for j := 0; j < len(b); j++ {
			policy = append(policy, b[j])
		}
	}
	params := &RMRParams{&xapp.RMRParams{}}
	params.Mtype = 20
	params.PayloadLen = len(policy)
	params.Payload = []byte(policy)
	params.Mbuf = nil
	params.Whid = ep.Whid
	params.Callid = call_id
	params.Timeout = 200
	state, retstr = xapp.Rmr.SendCallMsg(params.RMRParams)
	routestatus := strings.Split(retstr, " ")
	if state != C.RMR_OK && routestatus[0] != "OK" {
		xapp.Logger.Error("Updating Routes to Endpoint: " + ep.Uuid + " failed, call_id: " + strconv.Itoa(call_id) + " for xapp.Rmr.SendCallMsg " + " Route Update Status: " + routestatus[0])
		return false
	} else {
		xapp.Logger.Info("Update Routes to Endpoint: " + ep.Uuid + " successful, call_id: " + strconv.Itoa(call_id) + ", Payload length: " + strconv.Itoa(params.PayloadLen) + ", Route Update Status: " + routestatus[0] + "(# of Entries:" + strconv.Itoa(len(*policies)))
		return true
	}
}

func (c *RmrPush) CheckEndpoint(payload string) (ep *rtmgr.Endpoint) {
	return c.checkEndpoint(payload)
}

func (c *RmrPush) CreateEndpoint(rmrsrc string) (ep *string, whid int) {
	return c.createEndpoint(rmrsrc)
}

func (c *RmrPush) DistributeToEp(policies *[]string, ep string, whid int) error {
	xapp.Logger.Debug("Invoked: sbi.DistributeToEp")
	xapp.Logger.Debug("args: %v", *policies)

	if rmrdynamiccallid == 255 {
		rmrdynamiccallid = 201
	}

	go c.sendDynamicRoutes(ep, whid, policies, rmrdynamiccallid)
	rmrdynamiccallid++

	return nil
}

func (c *RmrPush) sendDynamicRoutes(ep string, whid int, policies *[]string, call_id int) bool {
	xapp.Logger.Debug("Invoked send_rt_process to endpoint: " + ep + " call_id: " + strconv.Itoa(call_id) + "whid: " + strconv.Itoa(whid))
	var state int
	var retstr string

	var policy = []byte{}

	for _, pe := range *policies {
		b := []byte(pe)
		for j := 0; j < len(b); j++ {
			policy = append(policy, b[j])
		}
	}
	params := &RMRParams{&xapp.RMRParams{}}
	params.Mtype = 20
	params.PayloadLen = len(policy)
	params.Payload = []byte(policy)
	params.Mbuf = nil
	params.Whid = whid
	params.Callid = call_id
	params.Timeout = 200
	state, retstr = xapp.Rmr.SendCallMsg(params.RMRParams)
	routestatus := strings.Split(retstr, " ")
	if state != C.RMR_OK && routestatus[0] != "OK" {
		xapp.Logger.Error("Updating Routes to Endpoint: " + ep + " failed, call_id: " + strconv.Itoa(call_id) + ",whi_id: " + strconv.Itoa(whid) + " for xapp.Rmr.SendCallMsg " + " Route Update Status: " + routestatus[0])
		return false
	} else {
		xapp.Logger.Info("Update Routes to Endpoint: " + ep + " successful, call_id: " + strconv.Itoa(call_id) + ", Payload length: " + strconv.Itoa(params.PayloadLen) + ",whid: " + strconv.Itoa(whid) + ", Route Update Status: " + routestatus[0] + "(# of Entries:" + strconv.Itoa(len(*policies)))
		return true
	}
}
