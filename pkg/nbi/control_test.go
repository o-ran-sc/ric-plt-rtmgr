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
  Mnemonic:     nbi_test.go
  Abstract:     NBI unit tests
  Date:         21 May 2019
*/

package nbi

import (
	"testing"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

func TestControlRun(t *testing.T) {
	var c =  NewControl()
	var params xapp.RMRParams
	//rp := make(chan xapp.RMRParams)
	var rmrmeid xapp.RMRMeid
	rmrmeid.RanName = "gnb1"
	params.Payload = []byte{1, 2, 3, 4}
	params.Mtype = 1234
	params.SubId = -1
	params.Meid = &rmrmeid
	params.Src = "sender"
	params.PayloadLen = 4

	go c.Consume(&params)
	go c.controlLoop()
}
