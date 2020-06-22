/*
w
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
  Mnemonic:	sbi.go
  Abstract:	Contains SBI (SouthBound Interface) module definitions and generic SBI components
  Date:		16 March 2019
*/

package sbi

import (
	"errors"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"net"
	"routing-manager/pkg/rtmgr"
	"strconv"
	"strings"
)

const DefaultRmrPipelineSocketPrefix = "tcp://"
const DefaultRmrPipelineSocketNumber = 4561
const PlatformType = "platform"

var (
	SupportedSbis = []*EngineConfig{
		{
			Name:        "rmrpush",
			Version:     "v1",
			Protocol:    "rmrpipeline",
			Instance:    NewRmrPush(),
			IsAvailable: true,
		},
	}
)

func GetSbi(sbiName string) (Engine, error) {
	for _, sbi := range SupportedSbis {
		if sbi.Name == sbiName && sbi.IsAvailable {
			return sbi.Instance, nil
		}
	}
	return nil, errors.New("SBI:" + sbiName + " is not supported or still not available")
}

type Sbi struct {
}

func (s *Sbi) pruneEndpointList(sbi Engine) {
	xapp.Logger.Debug("pruneEndpointList invoked.")
	for _, ep := range rtmgr.Eps {
		if !ep.Keepalive {
			xapp.Logger.Debug("deleting %v", ep)
			sbi.DeleteEndpoint(ep)
			delete(rtmgr.Eps, ep.Uuid)
		} else {
			rtmgr.Eps[ep.Uuid].Keepalive = false
		}
	}
}

func (s *Sbi) updateEndpoints(rcs *rtmgr.RicComponents, sbi Engine) {
	for _, xapps := range (*rcs).XApps {
		for _, instance := range xapps.Instances {
			uuid := instance.Ip + ":" + strconv.Itoa(int(instance.Port))
			if _, ok := rtmgr.Eps[uuid]; ok {
				rtmgr.Eps[uuid].Keepalive = true
			} else {
				ep := &rtmgr.Endpoint{
					Uuid:       uuid,
					Name:       instance.Name,
					XAppType:   xapps.Name,
					Ip:         instance.Ip,
					Port:       instance.Port,
					TxMessages: instance.TxMessages,
					RxMessages: instance.RxMessages,
					Policies:   instance.Policies,
					Socket:     nil,
					IsReady:    false,
					Keepalive:  true,
				}
				if err := sbi.AddEndpoint(ep); err != nil {
					xapp.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
					continue
				}
				rtmgr.Eps[uuid] = ep
			}
		}
	}
	s.updatePlatformEndpoints(&((*rcs).Pcs), sbi)
	s.updateE2TEndpoints(&((*rcs).E2Ts), sbi)
	s.pruneEndpointList(sbi)
}

func (s *Sbi) updatePlatformEndpoints(pcs *rtmgr.PlatformComponents, sbi Engine) {
	xapp.Logger.Debug("updatePlatformEndpoints invoked. PCS: %v", *pcs)
	for _, pc := range *pcs {
		uuid := pc.Fqdn + ":" + strconv.Itoa(int(pc.Port))
		if _, ok := rtmgr.Eps[uuid]; ok {
			rtmgr.Eps[uuid].Keepalive = true
		} else {
			ep := &rtmgr.Endpoint{
				Uuid:       uuid,
				Name:       pc.Name,
				XAppType:   PlatformType,
				Ip:         pc.Fqdn,
				Port:       pc.Port,
				TxMessages: rtmgr.PLATFORMMESSAGETYPES[pc.Name]["tx"],
				RxMessages: rtmgr.PLATFORMMESSAGETYPES[pc.Name]["rx"],
				Socket:     nil,
				IsReady:    false,
				Keepalive:  true,
			}
			xapp.Logger.Debug("ep created: %v", ep)
			if err := sbi.AddEndpoint(ep); err != nil {
				xapp.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
				continue
			}
			rtmgr.Eps[uuid] = ep
		}
	}
}

func (s *Sbi) updateE2TEndpoints(E2Ts *map[string]rtmgr.E2TInstance, sbi Engine) {
	xapp.Logger.Debug("updateE2TEndpoints invoked. E2T: %v", *E2Ts)
	for _, e2t := range *E2Ts {
		uuid := e2t.Fqdn
		stringSlice := strings.Split(e2t.Fqdn, ":")
		ipaddress := stringSlice[0]
		port, _ := strconv.Atoi(stringSlice[1])
		if _, ok := rtmgr.Eps[uuid]; ok {
			rtmgr.Eps[uuid].Keepalive = true
		} else {
			ep := &rtmgr.Endpoint{
				Uuid:       uuid,
				Name:       e2t.Name,
				XAppType:   PlatformType,
				Ip:         ipaddress,
				Port:       uint16(port),
				TxMessages: rtmgr.PLATFORMMESSAGETYPES[e2t.Name]["tx"],
				RxMessages: rtmgr.PLATFORMMESSAGETYPES[e2t.Name]["rx"],
				Socket:     nil,
				IsReady:    false,
				Keepalive:  true,
			}
			xapp.Logger.Debug("ep created: %v", ep)
			if err := sbi.AddEndpoint(ep); err != nil {
				xapp.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
				continue
			}
			rtmgr.Eps[uuid] = ep
		}
	}
}

func (s *Sbi) createEndpoint(payload string, sbi Engine) *rtmgr.Endpoint {
	xapp.Logger.Debug("CreateEndPoint %v", payload)
	stringSlice := strings.Split(payload, " ")
	uuid := stringSlice[0]
	xapp.Logger.Debug(">>> uuid %v", stringSlice[0])

	if _, ok := rtmgr.Eps[uuid]; ok {
		ep := rtmgr.Eps[uuid]
		return ep
	}

	/* incase the stored Endpoint list is in the form of IP:port*/
	stringsubsplit := strings.Split(uuid, ":")
	addr, err := net.LookupIP(stringsubsplit[0])
	if err == nil {
		convertedUuid := fmt.Sprintf("%s:%s", addr[0], stringsubsplit[1])
		xapp.Logger.Info(" IP:Port received is %s", convertedUuid)
		if _, ok := rtmgr.Eps[convertedUuid]; ok {
			ep := rtmgr.Eps[convertedUuid]
			return ep
		}
	}

	return nil
}
