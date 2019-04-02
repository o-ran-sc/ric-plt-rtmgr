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
  Mnemonic:	rtmgr/types.go
  Abstract:	Containes RTMGR (Routing Manager) specific types
  Date:		12 March 2019
*/

package rtmgr

type Endpoint struct {
	Name     string
	Type     string
	IpSocket string
}

type XApps struct {
	XApplist []XApp
}

type RouteTable []RouteTableEntry

type EndpointList []Endpoint

type RouteTableEntry struct {
	MessageType string
	TxList      EndpointList
	RxGroups    []EndpointList
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
	Port       int      `json:"port"`
	TxMessages []string `json:"txMessages"`
	RxMessages []string `json:"rxMessages"`
}
