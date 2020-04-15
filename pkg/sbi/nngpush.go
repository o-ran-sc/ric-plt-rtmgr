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

import (
	"bytes"
	"crypto/md5"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"routing-manager/pkg/rtmgr"
	"strconv"
	//"time"
	"fmt"
)

type NngPush struct {
	Sbi
	rcChan chan *xapp.RMRParams
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

	var policy = []byte{}
	cumulative_policy := 0
	count := 0
	maxrecord := xapp.Config.GetInt("maxrecord")
	if maxrecord == 0 {
		maxrecord = 10
	}

	for _, pe := range *policies {
		b := []byte(pe)
		for j := 0; j < len(b); j++ {
			policy = append(policy, b[j])
		}
		count++
		cumulative_policy++
		if count == maxrecord || cumulative_policy == len(*policies) {
			params := &RMRParams{&xapp.RMRParams{}}
			params.Mtype = 20
			params.PayloadLen = len(policy)
			params.Payload = []byte(policy)
			params.Mbuf = nil
			params.Whid = ep.Whid
			xapp.Rmr.SendMsg(params.RMRParams)
			count = 0
			policy = nil
			xapp.Logger.Debug("Sent message with payload len = %d", params.PayloadLen)
		}
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
