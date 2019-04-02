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
	"rtmgr"
	"time"
)

var myClient = &http.Client{Timeout: 1 * time.Second}

func fetchXappList(url string) (*[]rtmgr.XApp, error) {
	rtmgr.Logger.Debug("Invoked httpgetter.fetchXappList")
	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	rtmgr.Logger.Debug("http client raw response: %v", r)
	var xapps []rtmgr.XApp
	json.NewDecoder(r.Body).Decode(&xapps)
	rtmgr.Logger.Info("HTTP GET: OK")
	rtmgr.Logger.Debug("httpgetter.fetchXappList returns: %v", xapps)
	return &xapps, err
}
