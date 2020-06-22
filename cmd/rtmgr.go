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
	//"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	//"routing-manager/pkg/sbi"
	//"routing-manager/pkg/sdl"
	"syscall"
	"time"
)

const SERVICENAME = "rtmgr"

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

	SetupCloseHandler()

	xapp.Logger.Info("Start " + SERVICENAME + " service")
	rtmgr.Eps = make(rtmgr.Endpoints)
	rtmgr.Mtype = make(rtmgr.MessageTypeList)
	rtmgr.Rtmgr_ready = false

// RMR thread is starting port: 4560
	c := nbi.NewControl()
	go c.Run()

// Waiting for RMR to be ready
	time.Sleep(time.Duration(2) * time.Second)
	for xapp.Rmr.IsReady() == false {
	        time.Sleep(time.Duration(2) * time.Second)
	}

	dummy_whid := int(xapp.Rmr.Openwh("localhost:4560"))
	xapp.Logger.Info("created dummy Wormhole ID for routingmanager and dummy_whid :%d", dummy_whid)

	nbi.Serve()
	os.Exit(0)
}
