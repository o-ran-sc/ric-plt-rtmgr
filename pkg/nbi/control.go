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
	//"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"net/http"
	"os"
	"routing-manager/pkg/models"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sbi"
	"routing-manager/pkg/sdl"
	"strconv"
	"strings"
	"sync"
	"time"
)

var m sync.Mutex
var EndpointLock sync.Mutex

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

	xapp.Resource.InjectRoute("/ric/v1/symptomdata", c.SymptomDataHandler, "GET")

	xapp.Run(c)
}

func (c *Control) SymptomDataHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := DumpDebugData()
	xapp.Resource.SendSymptomDataJson(w, r, resp, "platform/rttable.json")
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
	}
}

func (c *Control) recievermr(msg *xapp.RMRParams) {
	xapp_msg := sbi.RMRParams{msg}
	switch msg.Mtype {
	case xapp.RICMessageTypes["RMRRM_REQ_TABLE"]:
		if rtmgr.Rtmgr_ready == false {
			xapp.Logger.Info("Update route Table Request(RMR -> RM), message discarded as routing manager is not ready")
		} else {
			xapp.Logger.Info("Update Route Table Request(RMR -> RM)")
			go c.handleUpdateToRoutingManagerRequest(msg)
		}
	case xapp.RICMessageTypes["RMRRM_TABLE_STATE"]:
		xapp.Logger.Info("State of route table to route mgr %s,payload %s", xapp_msg.String(), msg.Payload)
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
	if data == nil {
		if err != nil {
			xapp.Logger.Error("Cannot get data from sdl interface due to: " + err.Error())
			return
		} else {
			xapp.Logger.Debug("Cannot get data from sdl interface")
			return
		}
	}

	/* hack with WA only for mcxapp in near future */
	if strings.Contains(msg.String(), "ricxapp") {
		ep := sbiEngine.CheckEndpoint(string(params.Payload))
		if ep == nil {
			xapp.Logger.Error("Update Routing Table Request, can't handle due to end point %s is not available in complete ep list: ", string(params.Payload))
			return
		}
	}

	epstr, whid := sbiEngine.CreateEndpoint(msg.String())
	if epstr == nil || whid < 0 {
		xapp.Logger.Error("Wormhole Id creation failed %d for %s", whid, msg.String())
		return
	}

	/*This is to ensure the latest routes are sent.
	Assumption is that in this time interval the routes are built for this endpoint */
	time.Sleep(100 * time.Millisecond)
	policies := rpeEngine.GeneratePolicies(rtmgr.Eps, data)
	err = sbiEngine.DistributeToEp(policies, *epstr, whid)
	if err != nil {
		xapp.Logger.Error("Not able to publish the routing table due to: " + err.Error())
		return
	}
}

func getConfigData() (*rtmgr.RicComponents, error) {
	var data *rtmgr.RicComponents
	m.Lock()
	data, err := sdlEngine.ReadAll(xapp.Config.GetString("rtfile"))

	m.Unlock()
	if data == nil {
		if err != nil {
			return nil, errors.New("Cannot get data from sdl interface due to: " + err.Error())
		} else {
			xapp.Logger.Debug("Cannot get data from sdl interface due to data is null")
			return nil, errors.New("Cannot get data from sdl interface")
		}
	}

	return data, nil
}

func updateEp() (err error) {
	data, err := getConfigData()
	if err != nil {
		return errors.New("Routing table cannot be published due to: " + err.Error())
	}
	EndpointLock.Lock()
	sbiEngine.UpdateEndpoints(data)
	EndpointLock.Unlock()

	return nil
}

func sendPartialRoutesToAll(xappSubData *models.XappSubscriptionData, updatetype rtmgr.RMRUpdateType) (err error) {
	policies := rpeEngine.GeneratePartialPolicies(rtmgr.Eps, xappSubData, updatetype)
	err = sbiEngine.DistributeAll(policies)
	if err != nil {
		return errors.New("Routing table cannot be published due to: " + err.Error())
	}

	return nil
}

func sendRoutesToAll() (err error) {

	data, err := getConfigData()
	if err != nil {
		return errors.New("Routing table cannot be published due to: " + err.Error())
	}

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
		xapp.Logger.Error("Initialization of nbi failed due to: " + nbiErr.Error())
		return
	}

	err := sbiEngine.Initialize(xapp.Config.GetString("sbiurl"))
	if err != nil {
		xapp.Logger.Info("Failed to open push socket due to: " + err.Error())
		return
	}
	defer nbiEngine.Terminate()
	defer sbiEngine.Terminate()

	/* used for rtmgr restart case to connect to Endpoints */
	go updateEp()
	time.Sleep(5 * time.Second)
	sendRoutesToAll()
	for i := 0; i <= 5; i++ {
		/* Sometimes first message  fails, retry after 5 sec */
		time.Sleep(10 * time.Second)
		sendRoutesToAll()
	}

	for {
		xapp.Logger.Debug("Periodic Routes value = %s", xapp.Config.GetString("periodicroutes"))
		if xapp.Config.GetString("periodicroutes") == "enable" {
			go updateEp()
			time.Sleep(5 * time.Second)
			sendRoutesToAll()
		}

		rtmgr.Rtmgr_ready = true
		time.Sleep(INTERVAL * time.Second)
		xapp.Logger.Debug("Periodic loop timed out. Setting triggerSBI flag to distribute updated routes.")
	}
}
