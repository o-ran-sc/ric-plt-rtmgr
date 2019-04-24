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
  Mnemonic:	rtmgr/rtmgr.go
  Abstract:	Containes RTMGR (Routing Manager) module's generic variables and functions
  Date:		26 March 2019
*/

package rtmgr

import (
	"github.com/jcelliott/lumber"
)

var (
	//TODO: temporary solution
	// CamelCase Message Types are for being able to test with old fashioned admin controll xApps
	MESSAGETYPES = map[string]string{
		"HandoverPreparation":              "0",
		"HandoverCancel":                   "1",
		"LoadIndication":                   "2",
		"ErrorIndication":                  "3",
		"SNStatusTransfer":                 "4",
		"UEContextRelease":                 "5",
		"X2Setup":                          "6",
		"Reset":                            "7",
		"RIC_X2_SETUP":                     "10000",
		"RIC_X2_RESPONSE":                  "10001",
		"RIC_X2_RESOURCE_STATUS_REQUEST":   "10002",
		"RIC_X2_RESOURCE_STATUS_RESPONSE":  "10003",
		"RIC_X2_LOAD_INFORMATION":          "10004",
		"RIC_E2_TERMINATION_HC_REQUEST":    "10005",
		"RIC_E2_TERMINATION_HC_RESPONSE":   "10006",
		"RIC_E2_MANAGER_HC_REQUEST":        "10007",
		"RIC_E2_MANAGER_HC_RESPONSE":       "10008",
		"RIC_CONTROL_XAPP_CONFIG_REQUEST":  "100000",
		"RIC_CONTROL_XAPP_CONFIG_RESPONSE": "100001",
	}
	Logger = lumber.NewConsoleLogger(lumber.INFO)
	Eps Endpoints
)

func SetLogLevel(loglevel string) {
	switch loglevel {
	case "INFO":
		Logger.Level(lumber.INFO)
	case "WARN":
		Logger.Level(lumber.WARN)
	case "ERROR":
		Logger.Level(lumber.ERROR)
	case "DEBUG":
		Logger.Info("debugmode")
		Logger.Level(lumber.DEBUG)
	}
}

