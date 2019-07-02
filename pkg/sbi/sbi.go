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
        "strconv"
	"routing-manager/pkg/rtmgr"
)

const DEFAULT_NNG_PUBSUB_SOCKET_PREFIX = "tcp://"
const DEFAULT_NNG_PUBSUB_SOCKET_NUMBER = 4560
const DEFAULT_NNG_PIPELINE_SOCKET_PREFIX = "tcp://"
const DEFAULT_NNG_PIPELINE_SOCKET_NUMBER = 4561
const PLATFORMTYPE = "platformcomponenttype"

var (
	SupportedSbis = []*SbiEngineConfig{
		&SbiEngineConfig{
                        Name:     "nngpush",
                        Version:  "v1",
                        Protocol: "nngpipeline",
                        Instance: NewNngPush(),
                        IsAvailable: true,
                },
		&SbiEngineConfig{
                        Name:     "nngpub",
                        Version:  "v1",
                        Protocol: "nngpubsub",
                        Instance: NewNngPub(),
                        IsAvailable: true,
                },
	}
)

func GetSbi(sbiName string) (SbiEngine, error) {
	for _, sbi := range SupportedSbis {
		if sbi.Name == sbiName && sbi.IsAvailable {
			return sbi.Instance, nil
		}
	}
	return nil, errors.New("SBI:" + sbiName + " is not supported or still not available")
}

type Sbi struct {

}

func (s *Sbi) pruneEndpointList(sbii SbiEngine) {
        for _, ep := range rtmgr.Eps {
                if !ep.Keepalive {
			rtmgr.Logger.Debug("deleting %v",ep)
			sbii.DeleteEndpoint(ep)
                        delete(rtmgr.Eps, ep.Uuid)
                } else {
                        rtmgr.Eps[ep.Uuid].Keepalive = false
                }
        }
}

func (s *Sbi) updateEndpoints(rcs *rtmgr.RicComponents, sbii SbiEngine) {
	for _, xapp := range (*rcs).Xapps {
                for _, instance := range xapp.Instances {
                        uuid := instance.Ip + ":" + strconv.Itoa(int(instance.Port))
                        if _, ok := rtmgr.Eps[uuid]; ok {
                                rtmgr.Eps[uuid].Keepalive = true
                        } else {
                                ep := &rtmgr.Endpoint{
                                        uuid,
                                        instance.Name,
                                        xapp.Name,
                                        instance.Ip,
                                        instance.Port,
                                        instance.TxMessages,
                                        instance.RxMessages,
                                        nil,
                                        false,
                                        true,
                                }
                                if err := sbii.AddEndpoint(ep); err != nil {
                                        rtmgr.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
                                        continue
                                }
                                rtmgr.Eps[uuid] = ep
                        }
                }
        }
	s.updatePlatformEndpoints(&((*rcs).Pcs), sbii)
        s.pruneEndpointList(sbii)
}

func (s *Sbi ) updatePlatformEndpoints(pcs *rtmgr.PlatformComponents, sbii SbiEngine) {
	rtmgr.Logger.Debug("updatePlatformEndpoints invoked. PCS: %v", *pcs)
        for _, pc := range *pcs {
                uuid := pc.Fqdn + ":" + strconv.Itoa(int(pc.Port))
                if _, ok := rtmgr.Eps[uuid]; ok {
                        rtmgr.Eps[uuid].Keepalive = true
                } else {
                        ep := &rtmgr.Endpoint{
                                uuid,
                                pc.Name,
                                PLATFORMTYPE,
                                pc.Fqdn,
                                pc.Port,
                                rtmgr.PLATFORMMESSAGETYPES[pc.Name]["tx"],
                                rtmgr.PLATFORMMESSAGETYPES[pc.Name]["rx"],
                                nil,
                                false,
                                true,
                        }
			rtmgr.Logger.Debug("ep created: %v",ep)
                        if err := sbii.AddEndpoint(ep); err != nil {
                                rtmgr.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
                                continue
                        }
                        rtmgr.Eps[uuid] = ep
                }
        }
}

