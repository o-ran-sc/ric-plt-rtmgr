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
  Mnemonic:	stub.go
  Abstract:
  Date:		27 April 2019
*/

package stub

import "routing-manager/pkg/rtmgr"

var ValidXApps = &[]rtmgr.XApp{
	{Name: "app1", Status: "", Version: "", Instances: []rtmgr.XAppInstance{{Name: "E2TERM", Status: "unknown", Ip: "10.0.0.1", Port: 0, TxMessages: []string{"HandoverPreparation", "HandoverCancel"}, RxMessages: []string{"HandoverPreparation", "HandoverCancel"}}}},
	{Name: "app2", Status: "", Version: "", Instances: []rtmgr.XAppInstance{{Name: "SUBMAN", Status: "unknown", Ip: "192.168.0.1", Port: 0, TxMessages: []string{"HandoverCancel", "HandoverPreparation"}, RxMessages: []string{"HandoverPreparation", "HandoverCancel"}}}},
	{Name: "app3", Status: "", Version: "", Instances: []rtmgr.XAppInstance{{Name: "E2MAN", Status: "unknown", Ip: "10.1.1.1", Port: 0, TxMessages: []string{"X2Setup"}, RxMessages: []string{"Reset", "UEContextRelease"}}}},
	{Name: "app4", Status: "", Version: "", Instances: []rtmgr.XAppInstance{{Name: "UEMAN", Status: "unknown", Ip: "10.2.2.1", Port: 0, TxMessages: []string{"Reset", "UEContextRelease"}, RxMessages: []string{"", ""}}}},
}

var ValidPlatformComponents = &rtmgr.PlatformComponents{
	{Name: "E2TERM", Fqdn: "e2term", Port: 4561},
	{Name: "SUBMAN", Fqdn: "subman", Port: 4561},
	{Name: "E2MAN", Fqdn: "e2man", Port: 4561},
	{Name: "UEMAN", Fqdn: "ueman", Port: 4561},
}

var ValidEndpoints = []rtmgr.Endpoint{
	{Uuid: "10.0.0.1:0", Name: "E2TERM", XAppType: "app1", Ip: "", Port: 0, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: true, Keepalive: true},
	{Uuid: "192.168.0.1:0", Name: "SUBMAN", XAppType: "app2", Ip: "", Port: 0, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: false, Keepalive: false},
	{Uuid: "10.1.1.1:0", Name: "E2MAN", XAppType: "app3", Ip: "", Port: 0, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: true, Keepalive: false},
	{Uuid: "10.2.2.1:0", Name: "UEMAN", XAppType: "app4", Ip: "", Port: 0, TxMessages: []string{"", ""}, RxMessages: []string{"", ""}, Socket: nil, IsReady: false, Keepalive: true},
}

var ValidRicComponents = rtmgr.RicComponents{
	XApps: *ValidXApps, Pcs: *ValidPlatformComponents,
}

var ValidPolicies = &[]string{"", ""}
