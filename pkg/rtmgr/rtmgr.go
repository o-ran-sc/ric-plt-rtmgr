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
  Mnemonic:	rtmgr/rtmgr.go
  Abstract:	Contains RTMGR (Routing Manager) module's generic variables and functions
  Date:		26 March 2019
*/

package rtmgr

import (
	"encoding/json"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
)

var (
	Eps  Endpoints
	Subs SubscriptionList
	PrsCfg  *PlatformRoutes
	Mtype MessageTypeList
	RmrEp ProcessMultipleRMR
	DynamicRouteList []string
)

func GetPlatformComponents(configfile string) (*PlatformComponents, error) {
	xapp.Logger.Debug("Invoked rtmgr.GetPlatformComponents(" + configfile + ")")
	var rcfg ConfigRtmgr
	var rtroutes RtmgrRoutes
	var mtypes MessageTypeIdentifier
	yamlFile, err := os.Open(configfile)
	if err != nil {
		return nil, errors.New("cannot open the file due to: " + err.Error())
	}
	defer yamlFile.Close()
	byteValue, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		return nil, errors.New("cannot read the file due to: " + err.Error())
	}
	jsonByteValue, err := yaml.YAMLToJSON(byteValue)
	if err != nil {
		return nil, errors.New("cannot read the file due to: " + err.Error())
	}
	err = json.Unmarshal(jsonByteValue,&rtroutes)
        if err != nil {
               return nil, errors.New("cannot parse data due to: " + err.Error())
        }
        PrsCfg = &(rtroutes.Prs)

	err = json.Unmarshal(jsonByteValue,&mtypes)
        if err != nil {
               return nil, errors.New("cannot parse data due to: " + err.Error())
        } else {
		xapp.Logger.Debug("Messgaetypes = %v", mtypes)
		for _,m := range mtypes.Mit {
			splitstr := strings.Split(m,"=")
			Mtype[splitstr[0]] = splitstr[1]
		}
	}
	err = json.Unmarshal(jsonByteValue, &rcfg)
	if err != nil {
		return nil, errors.New("cannot parse data due to: " + err.Error())
	}
	xapp.Logger.Debug("Platform components read from the configfile:  %v", rcfg.Pcs)
	return &(rcfg.Pcs), nil
}
