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

import (
	"flag"
	"nbi"
	"os"
	"rpe"
	"rtmgr"
	"sbi"
	"sdl"
	"time"
)

const SERVICENAME = "rtmgr"
const INTERVAL time.Duration = 2

var (
	args *map[string]string
)

func parseArgs() {
	a := make(map[string]string)
	xmgeturl := flag.String("nbi-httpget", "http://localhost:3000/xapps", "xApp Manager URL")
	nngpubsock := flag.String("sbi-nngsub", "tcp://0.0.0.0:4560", "NNG Subsciption Socket URI")
	file := flag.String("sdl-file", "/db/rt.json", "Local file store location")
	rpename := flag.String("rpe", "rmr", "Policy Engine Module name")
	loglevel := flag.String("loglevel", "INFO", "INFO | WARN | ERROR | DEBUG")
	flag.Parse()
	if (*xmgeturl) != "" {
		a["xmurl"] = (*xmgeturl)
		a["nbiname"] = "httpGetter"
	}
	if (*nngpubsock) != "" {
		a["socketuri"] = (*nngpubsock)
		a["sbiname"] = "nngpub"
	}
	if (*file) != "" {
		a["file"] = (*file)
		a["sdlname"] = "file"
	}
	a["rpename"] = (*rpename)
	a["loglevel"] = (*loglevel)
	args = &a
}

func initRtmgr() (*nbi.NbiEngineConfig, *sbi.SbiEngineConfig, *sdl.SdlEngineConfig, *rpe.RpeEngineConfig, error) {
	var err error
	if nbi, err := nbi.GetNbi((*args)["nbiname"]); err == nil && nbi != nil {
		if sbi, err := sbi.GetSbi((*args)["sbiname"]); err == nil && sbi != nil {
			if sdl, err := sdl.GetSdl((*args)["sdlname"]); err == nil && sdl != nil {
				if rpe, err := rpe.GetRpe((*args)["rpename"]); err == nil && rpe != nil {
					return nbi, sbi, sdl, rpe, nil
				}
			}
		}
	}
	return nil, nil, nil, nil, err
}

func serve(nbi *nbi.NbiEngineConfig, sbi *sbi.SbiEngineConfig, sdl *sdl.SdlEngineConfig, rpe *rpe.RpeEngineConfig) {
	err := sbi.OpenSocket((*args)["socketuri"])
	if err != nil {
		rtmgr.Logger.Info("fail to open pub socket due to: " + err.Error())
		return
	}
	defer sbi.CloseSocket()
	for {
		time.Sleep(INTERVAL * time.Second)
		data, err := nbi.BatchFetch((*args)["xmurl"])
		if err != nil {
			rtmgr.Logger.Error("cannot get data from " + nbi.Engine.Name + " interface dute to: " + err.Error())
		} else {
			sdl.WriteAll((*args)["file"], data)
		}
		data, err = sdl.ReadAll((*args)["file"])
		if err != nil || data == nil {
			rtmgr.Logger.Error("cannot get data from " + sdl.Engine.Name + " interface dute to: " + err.Error())
			continue
		}
		policies := rpe.GeneratePolicies(data)
		err = sbi.DistributeAll(policies)
		if err != nil {
			rtmgr.Logger.Error("routing rable cannot be published due to: " + err.Error())
		}
	}
}

func main() {
	parseArgs()
	rtmgr.SetLogLevel((*args)["loglevel"])
	nbi, sbi, sdl, rpe, err := initRtmgr()
	if err != nil {
		rtmgr.Logger.Error(err.Error())
		os.Exit(1)
	}
	rtmgr.Logger.Info("Start " + SERVICENAME + " service")
	serve(nbi, sbi, sdl, rpe)
	os.Exit(0)
}
