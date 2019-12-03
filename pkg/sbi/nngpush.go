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
	"routing-manager/pkg/rtmgr"
	"strconv"
)

type NngPush struct {
	Sbi
	NewSocket CreateNewNngSocketHandler
}

func NewNngPush() *NngPush {
	instance := new(NngPush)
	instance.NewSocket = createNewPushSocket
	return instance
}

func createNewPushSocket() (NngSocket, error) {
	rtmgr.Logger.Debug("Invoked: createNewPushSocket()")
	socket, err := push.NewSocket()
	if err != nil {
		return nil, errors.New("can't create new push socket due to:" + err.Error())
	}
	socket.SetPipeEventHook(pipeEventHandler)
	return socket, nil
}

func pipeEventHandler(event mangos.PipeEvent, pipe mangos.Pipe) {
	rtmgr.Logger.Debug("Invoked: pipeEventHandler()")
	rtmgr.Logger.Debug("Received pipe event for " + pipe.Address() + " address")
	for _, ep := range rtmgr.Eps {
		uri := DefaultNngPipelineSocketPrefix + ep.Ip + ":" + strconv.Itoa(DefaultNngPipelineSocketNumber)
		if uri == pipe.Address() {
			switch event {
			case 1:
				ep.IsReady = true
				rtmgr.Logger.Debug("Endpoint " + uri + " successfully attached")
			default:
				ep.IsReady = false
				rtmgr.Logger.Debug("Endpoint " + uri + " has been detached")
			}
		}
	}
}

func (c *NngPush) Initialize(ip string) error {
	return nil
}

func (c *NngPush) Terminate() error {
	return nil
}

func (c *NngPush) AddEndpoint(ep *rtmgr.Endpoint) error {
	var err error
	var socket NngSocket
	rtmgr.Logger.Debug("Invoked sbi.AddEndpoint")
	rtmgr.Logger.Debug("args: %v", *ep)
	socket, err = c.NewSocket()
	if err != nil {
		return errors.New("can't add new socket to endpoint:" + ep.Uuid + " due to: " + err.Error())
	}
	ep.Socket = socket
	err = c.dial(ep)
	if err != nil {
		return errors.New("can't dial to endpoint:" + ep.Uuid + " due to: " + err.Error())
	}
	return nil
}

func (c *NngPush) DeleteEndpoint(ep *rtmgr.Endpoint) error {
	rtmgr.Logger.Debug("Invoked sbi. DeleteEndpoint")
	rtmgr.Logger.Debug("args: %v", *ep)
	if err := ep.Socket.(NngSocket).Close(); err != nil {
		return errors.New("can't close push socket of endpoint:" + ep.Uuid + " due to: " + err.Error())
	}
	return nil
}

func (c *NngPush) UpdateEndpoints(rcs *rtmgr.RicComponents) {
	c.updateEndpoints(rcs, c)
}

/*
NOTE: Asynchronous dial starts a goroutine which keep maintains the connection to the given endpoint
*/
func (c *NngPush) dial(ep *rtmgr.Endpoint) error {
	rtmgr.Logger.Debug("Dialing to endpoint: " + ep.Uuid)
	uri := DefaultNngPipelineSocketPrefix + ep.Ip + ":" + strconv.Itoa(DefaultNngPipelineSocketNumber)
	options := make(map[string]interface{})
	options[mangos.OptionDialAsynch] = true
	if err := ep.Socket.(NngSocket).DialOptions(uri, options); err != nil {
		return errors.New("can't dial on push socket to " + uri + " due to: " + err.Error())
	}
	return nil
}

func (c *NngPush) DistributeAll(policies *[]string) error {
	rtmgr.Logger.Debug("Invoked: sbi.DistributeAll")
	rtmgr.Logger.Debug("args: %v", *policies)
	for _, ep := range rtmgr.Eps {
		if ep.IsReady {
			go c.send(ep, policies)
		} else {
			rtmgr.Logger.Warn("Endpoint " + ep.Uuid + " is not ready")
		}
	}
	return nil
}

func (c *NngPush) send(ep *rtmgr.Endpoint, policies *[]string) {
	rtmgr.Logger.Debug("Push policy to endpoint: " + ep.Uuid)
	for _, pe := range *policies {
		if err := ep.Socket.(NngSocket).Send([]byte(pe)); err != nil {
			rtmgr.Logger.Error("Unable to send policy entry due to: " + err.Error())
		}
	}
	rtmgr.Logger.Info("NNG PUSH to endpoint " + ep.Uuid + ": OK (# of Entries:" + strconv.Itoa(len(*policies)) + ")")
}
