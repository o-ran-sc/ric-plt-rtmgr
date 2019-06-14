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
  Mnemonic:	NngPub.go
  Abstract:	mangos (NNG) Pub/Sub SBI implementation
  Date:		12 March 2019
*/

package sbi

import (
	"errors"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	_ "nanomsg.org/go/mangos/v2/transport/all"
	"routing-manager/pkg/rtmgr"
	"strconv"
)

type NngPub struct {
	Sbi
	socket NngSocket
	NewSocket CreateNewNngSocketHandler
}

func NewNngPub() *NngPub {
	instance := new(NngPub)
	instance.NewSocket = createNewPubSocket
	return instance
}

func createNewPubSocket() (NngSocket, error) {
	rtmgr.Logger.Debug("Invoked createNewPubSocket()")
	s, err := pub.NewSocket()
	if err != nil {
		return nil, errors.New("can't create new pub socket due to: " + err.Error())
	}
	return s, nil
}

func (c *NngPub) Initialize(ip string) error {
	rtmgr.Logger.Debug("Invoked sbi.Initialize("+ ip +")")
	var err error
	c.socket, err = c.NewSocket()
	if err != nil {
		return errors.New("create socket error due to: " + err.Error())
	}
	if err = c.listen(ip); err != nil {
		return errors.New("can't listen on socket due to: " + err.Error())
	}
	return nil
}

func (c *NngPub) Terminate() error {
	rtmgr.Logger.Debug("Invoked sbi.Terminate()")
	return c.closeSocket()
}

func (c *NngPub) AddEndpoint(ep *rtmgr.Endpoint) error {
	return nil
}

func (c *NngPub) DeleteEndpoint(ep *rtmgr.Endpoint) error {
	return nil
}

func (c *NngPub) UpdateEndpoints(xapps *[]rtmgr.XApp) {
	c.updateEndpoints(xapps, c)
}

func (c *NngPub) listen(ip string) error {
	rtmgr.Logger.Debug("Start listening on: " + ip)
	uri := DEFAULT_NNG_PUBSUB_SOCKET_PREFIX + ip + ":" + strconv.Itoa(DEFAULT_NNG_PUBSUB_SOCKET_NUMBER)
	rtmgr.Logger.Info("publishing on: " + uri)
	if err := c.socket.(NngSocket).Listen(uri); err != nil {
		return errors.New("can't publish on socket " + uri + " due to: " + err.Error())
	}
	return nil
}

func (c *NngPub) closeSocket() error {
	rtmgr.Logger.Debug("Close NngPub Socket")
	if err := c.socket.(NngSocket).Close(); err != nil {
		return errors.New("can't close socket due to: " + err.Error())
	}
	return nil
}

func (c *NngPub) DistributeAll(policies *[]string) error {
	rtmgr.Logger.Debug("Invoked: sbi.DistributeAll(), args: %v",(*policies))
	for _, pe := range *policies {
		if err := c.socket.(NngSocket).Send([]byte(pe)); err != nil {
			return errors.New("Unable to send policy entry due to: " + err.Error())
		}
	}
	rtmgr.Logger.Info("NNG PUB: OK (# of Entries: " + strconv.Itoa(len((*policies))) + ")")
	return nil
}
