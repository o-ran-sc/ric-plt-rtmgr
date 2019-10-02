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
	"flag"
	"os"
	"os/signal"
	"routing-manager/pkg/nbi"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sbi"
	"routing-manager/pkg/sdl"
	"syscall"
	"time"
)

const SERVICENAME = "rtmgr"
const INTERVAL time.Duration = 2

var (
	args map[string]*string
)

func parseArgs() {
	// TODO: arguments should be validated (filename; xm-url; sbi-if; rest-url; rest-port)
	args = make(map[string]*string)
	args["configfile"] = flag.String("configfile", "/etc/rtmgrcfg.json", "Routing manager's configuration file path")
	args["nbi"] = flag.String("nbi", "httpRESTful", "Northbound interface module to be used. Valid values are: 'httpGetter | httpRESTful'")
	args["sbi"] = flag.String("sbi", "nngpush", "Southbound interface module to be used. Valid values are: 'nngpush")
	args["rpe"] = flag.String("rpe", "rmrpush", "Route Policy Engine to be used. Valid values are: 'rmrpush'")
	args["sdl"] = flag.String("sdl", "file", "Data store engine to be used. Valid values are: 'file'")
	args["xm-url"] = flag.String("xm-url", "http://localhost:3000/xapps", "HTTP URL where xApp Manager exposes the entire xApp List")
	args["nbi-if"] = flag.String("nbi-if", "http://localhost:8888", "Base HTTP URL where routing manager will be listening on")
	args["sbi-if"] = flag.String("sbi-if", "0.0.0.0", "IPv4 address of interface where Southbound socket to be opened")
	args["filename"] = flag.String("filename", "/db/rt.json", "Absolute path of file where the route information to be stored")
	args["loglevel"] = flag.String("loglevel", "INFO", "INFO | WARN | ERROR | DEBUG | TRACE")
	flag.Parse()
}

func initRtmgr() (nbiEngine nbi.Engine, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine, err error) {
	if nbiEngine, err = nbi.GetNbi(*args["nbi"]); err == nil && nbiEngine != nil {
		if sbiEngine, err = sbi.GetSbi(*args["sbi"]); err == nil && sbiEngine != nil {
			if sdlEngine, err = sdl.GetSdl(*args["sdl"]); err == nil && sdlEngine != nil {
				if rpeEngine, err = rpe.GetRpe(*args["rpe"]); err == nil && rpeEngine != nil {
					return nbiEngine, sbiEngine, sdlEngine, rpeEngine, nil
				}
			}
		}
	}
	return nil, nil, nil, nil, err
}

func serveSBI(triggerSBI <-chan bool, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine) {
	for {
		if <-triggerSBI {
			data, err := sdlEngine.ReadAll(*args["filename"])
			if err != nil || data == nil {
				rtmgr.Logger.Error("Cannot get data from sdl interface due to: " + err.Error())
				continue
			}
			sbiEngine.UpdateEndpoints(data)
			policies := rpeEngine.GeneratePolicies(rtmgr.Eps)
			err = sbiEngine.DistributeAll(policies)
			if err != nil {
				rtmgr.Logger.Error("Routing table cannot be published due to: " + err.Error())
			}
		}
	}
}

func serve(nbiEngine nbi.Engine, sbiEngine sbi.Engine, sdlEngine sdl.Engine, rpeEngine rpe.Engine) {

	triggerSBI := make(chan bool)

	nbiErr := nbiEngine.Initialize(*args["xm-url"], *args["nbi-if"], *args["filename"], *args["configfile"],
		sdlEngine, rpeEngine, triggerSBI)
	if nbiErr != nil {
		rtmgr.Logger.Error("Failed to initialize nbi due to: " + nbiErr.Error())
		return
	}

	err := sbiEngine.Initialize(*args["sbi-if"])
	if err != nil {
		rtmgr.Logger.Info("Failed to open push socket due to: " + err.Error())
		return
	}
	defer nbiEngine.Terminate()
	defer sbiEngine.Terminate()

	// This SBI Go routine is trtiggered by periodic main loop and when data is recieved on REST interface.
	go serveSBI(triggerSBI, sbiEngine, sdlEngine, rpeEngine)

	for {
		time.Sleep(INTERVAL * time.Second)
		if *args["nbi"] == "httpGetter" {
			data, err := nbiEngine.(*nbi.HttpGetter).FetchAllXApps(*args["xm-url"])
			if err != nil {
				rtmgr.Logger.Error("Cannot fetch xapp data due to: " + err.Error())
			} else if data != nil {
				sdlEngine.WriteXApps(*args["filename"], data)
			}
		}

		triggerSBI <- true
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		rtmgr.Logger.Info("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func main() {
	parseArgs()
	rtmgr.SetLogLevel(*args["loglevel"])
	nbiEngine, sbiEngine, sdlEngine, rpeEngine, err := initRtmgr()
	if err != nil {
		rtmgr.Logger.Error(err.Error())
		os.Exit(1)
	}
	SetupCloseHandler()
	rtmgr.Logger.Info("Start " + SERVICENAME + " service")
	rtmgr.Eps = make(rtmgr.Endpoints)
	serve(nbiEngine, sbiEngine, sdlEngine, rpeEngine)
	os.Exit(0)
}
