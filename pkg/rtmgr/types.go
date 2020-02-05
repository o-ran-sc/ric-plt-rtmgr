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
  Mnemonic:	rtmgr/types.go
  Abstract:	Contains RTMGR (Routing Manager) specific types
  Date:		12 March 2019
*/

package rtmgr

type XApps struct {
	XAppList []XApp
}

type RouteTable []RouteTableEntry
type EndpointList []Endpoint

type Endpoints map[string]*Endpoint

type SubscriptionList []Subscription

//TODO: uuid is not a real UUID but a string of "ip:port"
// this should be changed to real UUID later on which should come from xApp Manager // petszila
type Endpoint struct {
	Uuid       string
	Name       string
	XAppType   string
	Ip         string
	Port       uint16
	TxMessages []string
	RxMessages []string
	Policies   []int32
	Socket     interface{}
	IsReady    bool
	Keepalive  bool
}

type RouteTableEntry struct {
	MessageType string
	TxList      EndpointList
	RxGroups    []EndpointList
	SubID       int32
	RouteType   string
}

type XApp struct {
	Name      string         `json:"name"`
	Status    string         `json:"status"`
	Version   string         `json:"version"`
	Instances []XAppInstance `json:"instances"`
}

type XAppInstance struct {
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	Ip         string   `json:"ip"`
	Port       uint16   `json:"port"`
	TxMessages []string `json:"txMessages"`
	RxMessages []string `json:"rxMessages"`
	Policies   []int32  `json:"policies"`
}

type PlatformComponents []struct {
	Name string `json:"name"`
	Fqdn string `json:"fqdn"`
	Port uint16 `json:"port"`
}

type E2TInstance struct {
        Name string `json:"name"`
        Fqdn string `json:"fqdn"`
        Ranlist []string `json:"ranlist"`
}

type ConfigRtmgr struct {
	Pcs PlatformComponents `json:"PlatformComponents"`
}

type RicComponents struct {
	XApps []XApp
	E2Ts  map [string]E2TInstance
	MeidMap  []string
	Pcs   PlatformComponents
}

type Subscription struct {
	SubID int32
	Fqdn string
	Port uint16
}

type PlatformRoutes []struct {
       MessageType     string `json:"messagetype"`
       SenderEndPoint  string `json:"senderendpoint"`
       SubscriptionId  int32  `json:"subscriptionid"`
       EndPoint        string `json:"endpoint"`
       Meid            string `json:"meid"`
}

type RtmgrRoutes struct {
       Prs PlatformRoutes      `json:"PlatformRoutes"`
}

type FqDn struct {
	Address *string
        Port *uint16
}

type XappList struct {
        SubscriptionID  uint16
	FqdnList []FqDn
}
