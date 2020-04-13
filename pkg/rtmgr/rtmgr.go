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
)

var (
	//TODO: temporary solution
	// CamelCase Message Types are for being able to test with old fashioned admin control xApps
	// TODO: Add a separate message definition file (Not using the one from RMR to not create dependency on that library).
	MessageTypes = map[string]string{
		"HandoverPreparation":              "0",
		"HandoverCancel":                   "1",
		"LoadIndication":                   "2",
		"ErrorIndication":                  "3",
		"SNStatusTransfer":                 "4",
		"UEContextRelease":                 "5",
		"X2Setup":                          "6",
		"Reset":                            "7",
		"E2_TERM_INIT":                     "1100",
		"E2_TERM_KEEP_ALIVE_REQ":            "1101",
		"E2_TERM_KEEP_ALIVE_RESP":           "1102",
		"RAN_CONNECTED":                    "1200",
		"RAN_RESTARTED":                    "1210",
		"RAN_RECONFIGURED":                 "1220",
		"RIC_SCTP_CLEAR_ALL":               "1090",
		"RIC_SCTP_CONNECTION_FAILURE":      "1080",
		"RIC_X2_SETUP":                     "10000",
		"RIC_X2_RESPONSE":                  "10001",
		"RIC_X2_RESOURCE_STATUS_REQUEST":   "10002",
		"RIC_X2_RESOURCE_STATUS_RESPONSE":  "10003",
		"RIC_X2_LOAD_INFORMATION":          "10004",
		"RIC_E2_TERMINATION_HC_REQUEST":    "10005",
		"RIC_E2_TERMINATION_HC_RESPONSE":   "10006",
		"RIC_E2_MANAGER_HC_REQUEST":        "10007",
		"RIC_E2_MANAGER_HC_RESPONSE":       "10008",
		"RIC_ENB_LOAD_INFORMATION":         "10020",
		"RIC_ERROR_INDICATION":             "10030",
		"RIC_X2_SETUP_REQ":                 "10060",
		"RIC_X2_SETUP_RESP":                "10061",
		"RIC_X2_SETUP_FAILURE":             "10062",
		"RIC_X2_RESET_REQ":                 "10070",
		"RIC_X2_RESET_RESP":                "10071",
		"RIC_ENB_CONF_UPDATE":              "10080",
		"RIC_ENB_CONF_UPDATE_ACK":          "10081",
		"RIC_ENB_CONF_UPDATE_FAILURE":      "10082",
		"RIC_RES_STATUS_REQ":               "10090",
		"RIC_RES_STATUS_RESP":              "10091",
		"RIC_RES_STATUS_FAILURE":           "10092",
		"RIC_RESOURCE_STATUS_UPDATE":       "10100",
		"RIC_ENDC_X2_SETUP_REQ":            "10360",
		"RIC_ENDC_X2_SETUP_RESP":           "10361",
		"RIC_ENDC_X2_SETUP_FAILURE":        "10362",
		"RIC_ENDC_CONF_UPDATE":             "10370",
		"RIC_ENDC_CONF_UPDATE_ACK":         "10371",
		"RIC_ENDC_CONF_UPDATE_FAILURE":     "10372",
		"RIC_GNB_STATUS_INDICATION":        "10450",
		"RIC_SUB_REQ":                      "12010",
		"RIC_SUB_RESP":                     "12011",
		"RIC_SUB_FAILURE":                  "12012",
		"RIC_SUB_DEL_REQ":                  "12020",
		"RIC_SUB_DEL_RESP":                 "12021",
		"RIC_SUB_DEL_FAILURE":              "12022",
		"RIC_CONTROL_REQ":                  "12040",
		"RIC_CONTROL_ACK":                  "12041",
		"RIC_CONTROL_FAILURE":              "12042",
		"RIC_INDICATION":                   "12050",
		"DC_ADM_INT_CONTROL":               "20000",
		"DC_ADM_INT_CONTROL_ACK":           "20001",
		"A1_POLICY_REQ":                    "20010",
		"A1_POLICY_RESPONSE":               "20011",
		"A1_POLICY_QUERY":                  "20012",
		"RIC_CONTROL_XAPP_CONFIG_REQUEST":  "100000",
		"RIC_CONTROL_XAPP_CONFIG_RESPONSE": "100001",
	}

	// Messagetype mappings for the platform components.
	// This implements static default routes needed by the RIC. Needs to be changed in case new components/message types needes to be added/updated.
	// Representation : {"componentName1": {"tx": <tx message type list>, "rx": <rx message type list>}}
	PLATFORMMESSAGETYPES = map[string]map[string][]string{
		"E2TERM":     {"tx": []string{"RIC_X2_SETUP_REQ", "RIC_X2_SETUP_RESP", "RIC_X2_SETUP_FAILURE", "RIC_X2_RESET", "RIC_X2_RESET_RESP", "RIC_ENDC_X2_SETUP_REQ", "RIC_ENDC_X2_SETUP_RESP", "RIC_ENDC_X2_SETUP_FAILURE", "RIC_SUB_RESP", "RIC_SUB_FAILURE", "RIC_SUB_DEL_RESP", "RIC_SUB_DEL_FAILURE"}, "rx": []string{"RIC_X2_SETUP_REQ", "RIC_X2_SETUP_RESP", "RIC_X2_SETUP_FAILURE", "RIC_X2_RESET", "RIC_X2_RESET_RESP", "RIC_ENDC_X2_SETUP_REQ", "RIC_ENDC_X2_SETUP_RESP", "RIC_ENDC_X2_SETUP_FAILURE", "RIC_SUB_REQ", "RIC_SUB_DEL_REQ", "RIC_CONTROL_REQ"}},
		"E2MAN":      {"tx": []string{"RIC_X2_SETUP_REQ", "RIC_X2_SETUP_RESP", "RIC_X2_SETUP_FAILURE", "RIC_X2_RESET", "RIC_X2_RESET_RESP", "RIC_ENDC_X2_SETUP_REQ", "RIC_ENDC_X2_SETUP_RESP", "RIC_ENDC_X2_SETUP_FAILURE"}, "rx": []string{"RIC_X2_SETUP_REQ", "RIC_X2_SETUP_RESP", "RIC_X2_SETUP_FAILURE", "RIC_X2_RESET", "RIC_X2_RESET_RESP", "RIC_ENDC_X2_SETUP_REQ", "RIC_ENDC_X2_SETUP_RESP", "RIC_ENDC_X2_SETUP_FAILURE"}},
		"SUBMAN":     {"tx": []string{"RIC_SUB_REQ", "RIC_SUB_DEL_REQ"}, "rx": []string{"RIC_SUB_RESP", "RIC_SUB_FAILURE", "RIC_SUB_DEL_RESP", "RIC_SUB_DEL_FAILURE"}},
		"UEMAN":      {"tx": []string{"RIC_CONTROL_REQ"}, "rx": []string{}},
		"RSM":        {"tx": []string{"RIC_RES_STATUS_REQ"}, "rx": []string{"RAN_CONNECTED", "RAN_RESTARTED", "RAN_RECONFIGURED"}},
		"A1MEDIATOR": {"tx": []string{}, "rx": []string{"A1_POLICY_QUERY", "A1_POLICY_RESPONSE"}},
	}

	Eps  Endpoints
	Sessions  SessionMap
	Subs SubscriptionList
	PrsCfg  *PlatformRoutes
)

func GetPlatformComponents(configfile string) (*PlatformComponents, error) {
	xapp.Logger.Debug("Invoked rtmgr.GetPlatformComponents(" + configfile + ")")
	var rcfg ConfigRtmgr
	var rtroutes RtmgrRoutes
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

	err = json.Unmarshal(jsonByteValue, &rcfg)
	if err != nil {
		return nil, errors.New("cannot parse data due to: " + err.Error())
	}
	xapp.Logger.Debug("Platform components read from the configfile:  %v", rcfg.Pcs)
	return &(rcfg.Pcs), nil
}
