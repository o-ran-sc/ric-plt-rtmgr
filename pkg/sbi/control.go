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
package sbi

import "C"

import (
        "errors"
        "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
        "strconv"
)


func NewControl() Control {

        return Control{make(chan *xapp.RMRParams)}
}


type Control struct {
        rcChan      chan *xapp.RMRParams
}


func (c *Control) Run() {
        go c.controlLoop()
        xapp.Run(c)
}

func (c *Control) Consume(rp *xapp.RMRParams) (err error) {
        c.rcChan <- rp
        return
}

func (c *Control) controlLoop() {
        for {
                msg := <-c.rcChan
                switch msg.Mtype {
                case xapp.RICMessageTypes["RIC_SUB_REQ"]:
                       xapp.Logger.Info("Message handling when RMR instance queries for Routes")
                default:
                        err := errors.New("Message Type " + strconv.Itoa(msg.Mtype) + " is discarded")
                        xapp.Logger.Error("Unknown message type: %v", err)
                }
        }
}

