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
	Mnemonic:	rtmgr.go
	Abstract:	Routing Manager Main file. Implemets the following functions:
			- parseArgs: reading command line arguments
			- init:Rtmgr initializing the service modules
			- serve: running the loop
	Date:		12 March 2019
*/
package main

//TODO: change flag to pflag (won't need any argument parse)

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"os"
	"os/signal"
	"routing-manager/pkg/nbi"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sbi"
	"routing-manager/pkg/sdl"
	"syscall"
	"time"
	"sync"
)

const SERVICENAME = "rtmgr"
const INTERVAL time.Duration = 60

func initRtmgr() (nbiEngine nbi.Engine, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, err error) {
	if nbiEngine, err = nbi.GetNbi(xapp.Config.GetString("nbi")); err == nil && nbiEngine != nil {
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



func serveSBI(triggerSBI <-chan bool, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, m *sync.Mutex) {
	for {
		if <-triggerSBI {
			m.Lock()
			data, err := sdlEngine.ReadAll(xapp.Config.GetString("rtfile"))
			m.Unlock()
			if err != nil || data == nil {
				xapp.Logger.Error("Cannot get data from sdl interface due to: " + err.Error())
				continue
			}
			sbiEngine.UpdateEndpoints(data)
			policies := rpeEngine.GeneratePolicies(rtmgr.Eps, data)
			err = sbiEngine.DistributeAll(policies)
			if err != nil {
				xapp.Logger.Error("Routing table cannot be published due to: " + err.Error())
			}
		}
	}
}

func sendRoutesToAll(sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine) {

	data, err := sdlEngine.ReadAll(xapp.Config.GetString("rtfile"))
	if err != nil || data == nil {
		xapp.Logger.Error("Cannot get data from sdl interface due to: " + err.Error())
		return
	}
	sbiEngine.UpdateEndpoints(data)
	policies := rpeEngine.GeneratePolicies(rtmgr.Eps, data)
	err = sbiEngine.DistributeAll(policies)
	if err != nil {
		xapp.Logger.Error("Routing table cannot be published due to: " + err.Error())
		return
	}
}


func serve(nbiEngine nbi.Engine, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, m *sync.Mutex) {

	triggerSBI := make(chan bool)

	nbiErr := nbiEngine.Initialize(xapp.Config.GetString("xmurl"), xapp.Config.GetString("nbiurl"), xapp.Config.GetString("rtfile"), xapp.Config.GetString("cfgfile"), xapp.Config.GetString("e2murl"), 
		sdlEngine, rpeEngine, triggerSBI, m)
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

	// This SBI Go routine is trtiggered by periodic main loop and when data is recieved on REST interface.
	go serveSBI(triggerSBI, sbiEngine, sdlEngine, rpeEngine, m)

	for {
		if xapp.Config.GetString("nbi") == "httpGetter" {
			data, err := nbiEngine.(*nbi.HttpGetter).FetchAllXApps(xapp.Config.GetString("xmurl"))
			if err != nil {
				xapp.Logger.Error("Cannot fetch xapp data due to: " + err.Error())
			} else if data != nil {
				sdlEngine.WriteXApps(xapp.Config.GetString("rtfile"), data)
			}
		}

		sendRoutesToAll(sbiEngine, sdlEngine, rpeEngine)

		rtmgr.Rtmgr_ready = true
		time.Sleep(INTERVAL * time.Second)
		xapp.Logger.Debug("Periodic loop timed out. Setting triggerSBI flag to distribute updated routes.")
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		xapp.Logger.Info("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func main() {

	nbiEngine, sbiEngine, sdlEngine, rpeEngine, err := initRtmgr()
	if err != nil {
		xapp.Logger.Error(err.Error())
		os.Exit(1)
	}

	SetupCloseHandler()

	xapp.Logger.Info("Start " + SERVICENAME + " service")
	rtmgr.Eps = make(rtmgr.Endpoints)
	rtmgr.Rtmgr_ready = false

	var m sync.Mutex

// RMR thread is starting port: 4560
	c := nbi.NewControl()
	go c.Run(sbiEngine, sdlEngine, rpeEngine, &m)

// Waiting for RMR to be ready
	time.Sleep(time.Duration(2) * time.Second)
	for xapp.Rmr.IsReady() == false {
	        time.Sleep(time.Duration(2) * time.Second)
	}

	dummy_whid := int(xapp.Rmr.Openwh("localhost:4560"))
	xapp.Logger.Info("created dummy Wormhole ID for routingmanager and dummy_whid :%d", dummy_whid)

	serve(nbiEngine, sbiEngine, sdlEngine, rpeEngine, &m)
	os.Exit(0)
}
