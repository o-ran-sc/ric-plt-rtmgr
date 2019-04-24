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
	"fmt"
	"rtmgr"
        "strconv"
)

const DEFAULT_NNG_PUBSUB_SOCKET_PREFIX = "tcp://"
const DEFAULT_NNG_PUBSUB_SOCKET_NUMBER = 4560
const DEFAULT_NNG_PIPELINE_SOCKET_PREFIX = "tcp://"
const DEFAULT_NNG_PIPELINE_SOCKET_NUMBER = 4561

var (
	SupportedSbis = []*SbiEngineConfig{
		&SbiEngineConfig{
			SbiEngine{
				Name:     "nngpub",
				Version:  "v1",
				Protocol: "nngpubsub",
			},
			openSocket(openNngPub),
			closeSocket(closeNngPub),
			createEndpointSocket(createNngPubEndpointSocket),
			destroyEndpointSocket(createNngPubEndpointSocket),
			distributeAll(publishAll),
			true,
		},
		&SbiEngineConfig{
			SbiEngine{
				Name:     "nngpush",
				Version:  "v1",
				Protocol: "nngpipeline",
			},
			openSocket(openNngPush),
			closeSocket(closeNngPush),
			createEndpointSocket(createNngPushEndpointSocket),
			destroyEndpointSocket(destroyNngPushEndpointSocket),
			distributeAll(pushAll),
			true,
		},
	}
)

func ListSbis() {
	fmt.Printf("SBI:\n")
	for _, sbi := range SupportedSbis {
		if sbi.IsAvailable {
			rtmgr.Logger.Info(sbi.Engine.Name + "/" + sbi.Engine.Version)
		}
	}
}

func GetSbi(sbiName string) (*SbiEngineConfig, error) {
	for _, sbi := range SupportedSbis {
		if (*sbi).Engine.Name == sbiName && (*sbi).IsAvailable {
			return sbi, nil
		}
	}
	return nil, errors.New("SBI:" + sbiName + " is not supported or still not available")
}

func pruneEndpointList(sbii *SbiEngineConfig) {
        for _, ep := range rtmgr.Eps {
                if !ep.Keepalive {
			sbii.DestroyEndpointSocket(ep)
                        delete(rtmgr.Eps, ep.Uuid)
                } else {
                        rtmgr.Eps[ep.Uuid].Keepalive = false
                }
        }
}

func UpdateEndpointList(xapps *[]rtmgr.XApp, sbii *SbiEngineConfig) {
        for _, xapp := range *xapps {
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
                                if err := sbii.CreateEndpointSocket(ep); err != nil {
                                        rtmgr.Logger.Error("can't create socket for endpoint: " + ep.Name + " due to:" + err.Error())
                                        continue
                                }
                                rtmgr.Eps[uuid] = ep
                        }
                }
        }
        pruneEndpointList(sbii)
}
