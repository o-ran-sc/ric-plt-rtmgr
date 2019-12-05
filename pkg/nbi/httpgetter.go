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
  Mnemonic:	httpgetter.go
  Abstract:	HTTPgetter NBI implementation
  		Simple HTTP getter solution.
  Date:		15 March 2019
*/

package nbi

import (
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"net/http"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sdl"
	"time"
)

type HttpGetter struct {
	Engine
	FetchAllXApps FetchAllXAppsHandler
}

func NewHttpGetter() *HttpGetter {
	instance := new(HttpGetter)
	instance.FetchAllXApps = fetchAllXApps
	return instance
}

var myClient = &http.Client{Timeout: 5 * time.Second}

func fetchAllXApps(xmurl string) (*[]rtmgr.XApp, error) {
	xapp.Logger.Info("Invoked httpGetter.fetchXappList: " + xmurl)
	r, err := myClient.Get(xmurl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 200 {
		xapp.Logger.Debug("http client raw response: %v", r)
		var xApps []rtmgr.XApp
		err = json.NewDecoder(r.Body).Decode(&xApps)
		if err != nil {
			xapp.Logger.Warn("Json decode failed: " + err.Error())
		}
		xapp.Logger.Info("HTTP GET: OK")
		xapp.Logger.Debug("httpGetter.fetchXappList returns: %v", xApps)
		return &xApps, err
	}
	xapp.Logger.Warn("httpGetter got an unexpected http status code: %v", r.StatusCode)
	return nil, nil
}

func (g *HttpGetter) Initialize(xmurl string, nbiif string, fileName string, configfile string,
	sdlEngine sdl.Engine, rpeEngine rpe.Engine, triggerSBI chan<- bool) error {
	return nil
}

func (g *HttpGetter) Terminate() error {
	return nil
}
