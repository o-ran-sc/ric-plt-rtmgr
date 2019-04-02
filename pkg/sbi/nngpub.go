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
  Mnemonic:	nngpub.go
  Abstract:	mangos (NNG) Pub/Sub SBI implementation
  Date:		12 March 2019
*/

package sbi

import (
	"errors"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	_ "nanomsg.org/go/mangos/v2/transport/all"
	"rtmgr"
	"strconv"
)

var sock mangos.Socket

/*
Creates the NNG publication channel
*/
func openNngPub(url string) error {
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		return errors.New("can't get new pub socket due to:" + err.Error())
	}
	rtmgr.Logger.Info("publishing on: " + url)
	if err = sock.Listen(url); err != nil {
		return errors.New("can't publish on socket " + url + " due to:" + err.Error())
	}
	return nil
}

func closeNngPub() error {
	if err := sock.Close(); err != nil {
		return errors.New("can't close socket due to:" + err.Error())
	}
	return nil
}

func publishAll(policies *[]string) error {
	for _, pe := range *policies {
		if err := sock.Send([]byte(pe)); err != nil {
			return errors.New("Unable to send policy entry due to: " + err.Error())
		}
	}
	rtmgr.Logger.Info("NNG PUB: OK (# of Entries:" + strconv.Itoa(len((*policies))) + ")")
	return nil
}
