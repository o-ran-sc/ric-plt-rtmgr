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
  Mnemonic:	nngpipe.go
  Abstract: mangos (NNG) Pipeline SBI implementation
  Date:		12 March 2019
*/

package sbi

import (
	"errors"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/push"
	_ "nanomsg.org/go/mangos/v2/transport/all"
	"rtmgr"
	"strconv"
)

func openNngPush(ip string) error {
	return nil
}

func closeNngPush() error {
	return nil
}

func createNngPushEndpointSocket(ep *rtmgr.Endpoint) error {
	rtmgr.Logger.Debug("Invoked sbi.createNngPushEndpointSocket")
	rtmgr.Logger.Debug("args: %v", (*ep))
	s, err := push.NewSocket()
	if err != nil {
		return errors.New("can't open push socket for endpoint: " + ep.Name +" due to:" + err.Error())
	}
	s.SetPipeEventHook(pipeEventHandler)
	ep.Socket = s
	dial(ep)
	return nil
}

func destroyNngPushEndpointSocket(ep *rtmgr.Endpoint) error {
	rtmgr.Logger.Debug("Invoked sbi.destroyNngPushEndpointSocket")
	rtmgr.Logger.Debug("args: %v", (*ep))
	if err:= ep.Socket.(mangos.Socket).Close(); err != nil {
			return errors.New("can't close push socket of endpoint:" + ep.Uuid + " due to:" + err.Error())
		}
	return nil
}

func pipeEventHandler(event mangos.PipeEvent, pipe mangos.Pipe) {
	for _, ep := range rtmgr.Eps {
		uri := DEFAULT_NNG_PIPELINE_SOCKET_PREFIX + ep.Ip + ":" + strconv.Itoa(DEFAULT_NNG_PIPELINE_SOCKET_NUMBER)
		if uri == pipe.Address() {
			switch event {
			case 1:
				ep.IsReady = true
				rtmgr.Logger.Debug("Endpoint " + uri + " successfully registered")
			default:
				ep.IsReady = false
				rtmgr.Logger.Debug("Endpoint " + uri + " has been deregistered")
			}
		}	
	}
}

/*
NOTE: Asynchronous dial starts a goroutine which keep maintains the connection to the given endpoint
*/
func dial(ep *rtmgr.Endpoint) {
	rtmgr.Logger.Debug("Dialing to endpoint: " + ep.Uuid)
	uri := DEFAULT_NNG_PIPELINE_SOCKET_PREFIX + ep.Ip + ":" + strconv.Itoa(DEFAULT_NNG_PIPELINE_SOCKET_NUMBER)
	options := make(map[string]interface{})
	options[mangos.OptionDialAsynch] = true
	if err := ep.Socket.(mangos.Socket).DialOptions(uri, options); err != nil {
		rtmgr.Logger.Error("can't dial on push socket to " + uri + " due to:" + err.Error())
	}
}

func pushAll(policies *[]string) error {
	rtmgr.Logger.Debug("Invoked: sbi.pushAll")
	rtmgr.Logger.Debug("args: %v", (*policies))
	for _, ep := range rtmgr.Eps {
		if ep.IsReady {
			go send(ep, policies)
		} else {
			rtmgr.Logger.Warn("Endpoint " + ep.Uuid + "is not ready")
		}
	}
	return nil
}

func send(ep *rtmgr.Endpoint, policies *[]string) {
	rtmgr.Logger.Debug("Invoked: sbi.pushAll")
	rtmgr.Logger.Debug("Push policy to endpoint: "+ ep.Uuid)
	for _, pe := range *policies {
		if err := ep.Socket.(mangos.Socket).Send([]byte(pe)); err != nil {
			rtmgr.Logger.Error("Unable to send policy entry due to: " + err.Error())
		}
	}
	rtmgr.Logger.Info("NNG PUSH to ednpoint " + ep.Uuid + ": OK (# of Entries:" + strconv.Itoa(len((*policies))) + ")")
}
