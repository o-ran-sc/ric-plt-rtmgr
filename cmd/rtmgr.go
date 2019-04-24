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
	args map[string]*string
)

func parseArgs() {
	// TODO: arguments should be validated (filename; xm-url; sbi-if)
	args = make(map[string]*string)
	args["nbi"] = flag.String("nbi", "httpGetter", "Northbound interface module to be used. Valid values are: 'httpGetter'")
	args["sbi"] = flag.String("sbi", "nngpush", "Southbound interface module to be used. Valid values are: 'nngpush | nngpub'")
	args["rpe"] = flag.String("rpe", "rmrpush", "Route Policy Engine to be used. Valid values are: 'rmrpush | rmrpub'")
	args["sdl"] = flag.String("sdl", "file", "Datastore enginge to be used. Valid values are: 'file'")
	args["xm-url"] = flag.String("xm-url", "http://localhost:3000/xapps", "HTTP URL where xApp Manager exposes the entire xApp List")
	args["sbi-if"] = flag.String("sbi-if", "0.0.0.0", "IPv4 address of interface where Southbound socket to be opened")
	args["filename"] = flag.String("filename", "/db/rt.json", "Absolute path of file where the route information to be stored")
	args["loglevel"] = flag.String("loglevel", "INFO", "INFO | WARN | ERROR | DEBUG")
	flag.Parse()
}

func initRtmgr() (*nbi.NbiEngineConfig, *sbi.SbiEngineConfig, *sdl.SdlEngineConfig, *rpe.RpeEngineConfig, error) {
	var err error
	var nbii *nbi.NbiEngineConfig
	var sbii *sbi.SbiEngineConfig
	var sdli *sdl.SdlEngineConfig
	var rpei *rpe.RpeEngineConfig
	if nbii, err = nbi.GetNbi(*args["nbi"]); err == nil && nbii != nil {
		if sbii, err = sbi.GetSbi(*args["sbi"]); err == nil && sbii != nil {
			if sdli, err = sdl.GetSdl(*args["sdl"]); err == nil && sdli != nil {
				if rpei, err = rpe.GetRpe(*args["rpe"]); err == nil && rpei != nil {
					return nbii, sbii, sdli, rpei, nil
				}
			}
		}
	}
	return nil, nil, nil, nil, err
}

func serve(nbii *nbi.NbiEngineConfig, sbii *sbi.SbiEngineConfig, sdli *sdl.SdlEngineConfig, rpei *rpe.RpeEngineConfig) {
	err := sbii.OpenSocket(*args["sbi-if"])
	if err != nil {
		rtmgr.Logger.Info("fail to open pub socket due to: " + err.Error())
		return
	}
	defer sbii.CloseSocket()
	for {
		time.Sleep(INTERVAL * time.Second)
		data, err := nbii.BatchFetch(*args["xm-url"])
		if err != nil {
			rtmgr.Logger.Error("cannot get data from " + nbii.Engine.Name + " interface dute to: " + err.Error())
		} else {
			sdli.WriteAll(*args["filename"], data)
		}
		data, err = sdli.ReadAll(*args["filename"])
		if err != nil || data == nil {
			rtmgr.Logger.Error("cannot get data from " + sdli.Engine.Name + " interface dute to: " + err.Error())
			continue
		}
		sbi.UpdateEndpointList(data, sbii)
		policies := rpei.GeneratePolicies(rtmgr.Eps)
		err = sbii.DistributeAll(policies)
		if err != nil {
			rtmgr.Logger.Error("routing rable cannot be published due to: " + err.Error())
		}
	}
}

func main() {
	parseArgs()
	rtmgr.SetLogLevel(*args["loglevel"])
	nbii, sbii, sdli, rpei, err := initRtmgr()
	if err != nil {
		rtmgr.Logger.Error(err.Error())
		os.Exit(1)
	}
	rtmgr.Logger.Info("Start " + SERVICENAME + " service")
	rtmgr.Eps = make(rtmgr.Endpoints)
	serve(nbii, sbii, sdli, rpei)
	os.Exit(0)
}
