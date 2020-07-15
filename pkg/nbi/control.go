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
	"time"
	"os"
	"fmt"
)

var m sync.Mutex

var nbiEngine Engine
var sbiEngine sbi.Engine
var sdlEngine sdl.Engine
var rpeEngine rpe.Engine

const INTERVAL time.Duration = 60

func NewControl() Control {
	return Control{make(chan *xapp.RMRParams)}
}

type Control struct {
	rcChan chan *xapp.RMRParams
}


func (c *Control) Run() {
	var err error
	go c.controlLoop()
	nbiEngine, sbiEngine, sdlEngine, rpeEngine, err = initRtmgr()
	if err != nil {
                xapp.Logger.Error(err.Error())
                os.Exit(1)
        }
	xapp.Run(c)
}

func (c *Control) Consume(rp *xapp.RMRParams) (err error) {
	c.rcChan <- rp
	return
}

func initRtmgr() (nbiEngine Engine, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, err error) {
        if nbiEngine, err = GetNbi(xapp.Config.GetString("nbi")); err == nil && nbiEngine != nil {
                if sbiEngine, err = sbi.GetSbi(xapp.Config.GetString("sbi")); err == nil && sbiEngine != nil {
                        if sdlEngine, err = sdl.GetSdl(xapp.Config.GetString("sdl")); err == nil && sdlEngine != nil {
                                if rpeEngine, err = rpe.GetRpe(xapp.Config.GetString("rpe")); err == nil && rpeEngine != nil {
                                        return nbiEngine, sbiEngine, sdlEngine, rpeEngine, nil
                                }
                        }
                }
        }
        return nil, nil, nil, nil, err
}

func (c *Control) controlLoop() {
	for {
		msg := <-c.rcChan
		c.recievermr(msg)
		/*
		xapp_msg := sbi.RMRParams{msg}
		switch msg.Mtype {
		case xapp.RICMessageTypes["RMRRM_REQ_TABLE"]:
			if rtmgr.Rtmgr_ready == false {
				xapp.Logger.Info("Update Route Table Request(RMR to RM), message discarded as routing manager is not ready")
			} else {
				xapp.Logger.Info("Update Route Table Request(RMR to RM)")
				go c.handleUpdateToRoutingManagerRequest(msg)
			}
		case xapp.RICMessageTypes["RMRRM_TABLE_STATE"]:
			xapp.Logger.Info("state of table to route mgr %s,payload %s", xapp_msg.String(), msg.Payload)

		default:
			err := errors.New("Message Type " + strconv.Itoa(msg.Mtype) + " is discarded")
			xapp.Logger.Error("Unknown message type: %v", err)
		}
		xapp.Rmr.Free(msg.Mbuf)*/
	}
}

func (c *Control) recievermr(msg *xapp.RMRParams) {
	xapp_msg := sbi.RMRParams{msg}
        switch msg.Mtype {
        case xapp.RICMessageTypes["RMRRM_REQ_TABLE"]:
	if rtmgr.Rtmgr_ready == false {
		xapp.Logger.Info("Update Route Table Request(RMR to RM), message discarded as routing manager is not ready")
        } else {
                xapp.Logger.Info("Update Route Table Request(RMR to RM)")
                go c.handleUpdateToRoutingManagerRequest(msg)
        }
        case xapp.RICMessageTypes["RMRRM_TABLE_STATE"]:
                xapp.Logger.Info("state of table to route mgr %s,payload %s", xapp_msg.String(), msg.Payload)
        default:
                err := errors.New("Message Type " + strconv.Itoa(msg.Mtype) + " is discarded")
                xapp.Logger.Error("Unknown message type: %v", err)
        }
        xapp.Rmr.Free(msg.Mbuf)
}

func (c *Control) handleUpdateToRoutingManagerRequest(params *xapp.RMRParams) {

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

func sendRoutesToAll() (err error) {

        m.Lock()
        data, err := sdlEngine.ReadAll(xapp.Config.GetString("rtfile"))
	fmt.Printf("data = %v,%v,%v",data,sdlEngine,sbiEngine)
        m.Unlock()
        if err != nil || data == nil {
                return errors.New("Cannot get data from sdl interface due to: " + err.Error())
        }
	if sbiEngine == nil {
		fmt.Printf("SBI is nil")
	}
        sbiEngine.UpdateEndpoints(data)
        policies := rpeEngine.GeneratePolicies(rtmgr.Eps, data)
        err = sbiEngine.DistributeAll(policies)
        if err != nil {
                return errors.New("Routing table cannot be published due to: " + err.Error())
        }

	return nil
}

func Serve() {

        nbiErr := nbiEngine.Initialize(xapp.Config.GetString("xmurl"), xapp.Config.GetString("nbiurl"), xapp.Config.GetString("rtfile"), xapp.Config.GetString("cfgfile"), xapp.Config.GetString("e2murl"), sdlEngine, rpeEngine, &m)
        if nbiErr != nil {
                xapp.Logger.Error("Failed to initialize nbi due to: " + nbiErr.Error())
                return
        }

        err := sbiEngine.Initialize(xapp.Config.GetString("sbiurl"))
        if err != nil {
                xapp.Logger.Info("Failed to open push socket due to: " + err.Error())
                return
        }
        defer nbiEngine.Terminate()
        defer sbiEngine.Terminate()

        for {
                sendRoutesToAll()

                rtmgr.Rtmgr_ready = true
                time.Sleep(INTERVAL * time.Second)
                xapp.Logger.Debug("Periodic loop timed out. Setting triggerSBI flag to distribute updated routes.")
        }
}
