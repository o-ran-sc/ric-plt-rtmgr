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
package nbi

import "C"

import (
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sbi"
	"routing-manager/pkg/sdl"
	"strconv"
	"sync"
)

func NewControl() Control {

	return Control{make(chan *xapp.RMRParams)}
}

type Control struct {
	rcChan chan *xapp.RMRParams
}

func (c *Control) Run(sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, m *sync.Mutex) {
	go c.controlLoop(sbiEngine, sdlEngine, rpeEngine, m)
	xapp.Run(c)
}

func (c *Control) Consume(rp *xapp.RMRParams) (err error) {
	c.rcChan <- rp
	return
}

func (c *Control) controlLoop(sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, m *sync.Mutex) {
	for {
		msg := <-c.rcChan
		xapp_msg := sbi.RMRParams{msg}
		switch msg.Mtype {
		case xapp.RICMessageTypes["RMRRM_REQ_TABLE"]:
			if rtmgr.Rtmgr_ready == false {
				xapp.Logger.Info("Update Route Table Request(RMR to RM), message discarded as routing manager is not ready")
			} else {
				xapp.Logger.Info("Update Route Table Request(RMR to RM)")
				go c.handleUpdateToRoutingManagerRequest(msg, sbiEngine, sdlEngine, rpeEngine, m)
			}
		case xapp.RICMessageTypes["RMRRM_TABLE_STATE"]:
			xapp.Logger.Info("state of table to route mgr %s,payload %s", xapp_msg.String(), msg.Payload)

		default:
			err := errors.New("Message Type " + strconv.Itoa(msg.Mtype) + " is discarded")
			xapp.Logger.Error("Unknown message type: %v", err)
		}
		xapp.Rmr.Free(msg.Mbuf)
	}
}

func (c *Control) handleUpdateToRoutingManagerRequest(params *xapp.RMRParams, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, m *sync.Mutex) {

	msg := sbi.RMRParams{params}

	xapp.Logger.Info("Update Route Table Request, msg.String() : %s", msg.String())
	xapp.Logger.Info("Update Route Table Request, params.Payload : %s", string(params.Payload))

	m.Lock()
	data, err := sdlEngine.ReadAll(xapp.Config.GetString("rtfile"))
	m.Unlock()
	if err != nil || data == nil {
		xapp.Logger.Error("Cannot get data from sdl interface due to: " + err.Error())
		return
	}

	ep := sbiEngine.CreateEndpoint(string(params.Payload))
	if ep == nil {
		xapp.Logger.Error("Update Routing Table Request can't handle due to end point %s is not avail in complete ep list: ", string(params.Payload))
		return
	}

	policies := rpeEngine.GeneratePolicies(rtmgr.Eps, data)
	err = sbiEngine.DistributeToEp(policies, ep)
	if err != nil {
		xapp.Logger.Error("Routing table cannot be published due to: " + err.Error())
		return
	}
}
