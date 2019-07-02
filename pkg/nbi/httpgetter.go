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
  Mnemonic:	httpgetter.go
  Abstract:	HTTPgetter NBI implementation
  		Simple HTTP getter solution. Only for testing purpose.
  Date:		15 March 2019
*/

package nbi

import (
	"encoding/json"
	"net/http"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/sdl"
	"time"
)

type HttpGetter struct {
	NbiEngine
	FetchAllXapps FetchAllXappsHandler
}

func NewHttpGetter() *HttpGetter {
	instance := new(HttpGetter)
	instance.FetchAllXapps = fetchAllXapps
	return instance
}

var myClient = &http.Client{Timeout: 1 * time.Second}

func fetchAllXapps(xmurl string) (*[]rtmgr.XApp, error) {
	rtmgr.Logger.Info("Invoked httpgetter.fetchXappList: " + xmurl)
	r, err := myClient.Get(xmurl)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == 200 {
		rtmgr.Logger.Debug("http client raw response: %v", r)
		var xapps []rtmgr.XApp
		err = json.NewDecoder(r.Body).Decode(&xapps)
		if err != nil {
			rtmgr.Logger.Warn("Json decode failed: " + err.Error())
		}
		rtmgr.Logger.Info("HTTP GET: OK")
		rtmgr.Logger.Debug("httpgetter.fetchXappList returns: %v", xapps)
		return &xapps, err
	}
	rtmgr.Logger.Warn("httpgetter got an unexpected http status code: %v", r.StatusCode)
	return nil, nil
}

func (g *HttpGetter) Initialize(xmurl string, nbiif string, fileName string, configfile string,
				sdlEngine sdl.SdlEngine, rpeEngine rpe.RpeEngine, triggerSBI chan<- bool) error {
	return nil
}

func (g *HttpGetter) Terminate() error {
	return nil
}
