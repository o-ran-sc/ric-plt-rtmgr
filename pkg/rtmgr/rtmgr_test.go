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
	Mnemonic:	rtmgr_test.go
	Abstract:
	Date:		14 May 2019
*/

package rtmgr

import (
	"testing"
)

func TestGetPlatformComponents(t *testing.T) {
	//Check epty file
	_, err := GetPlatformComponents("")
	t.Log(err)

	//Valid JSON file
	_, err = GetPlatformComponents("/tmp/go/src/routing-manager/manifests/rtmgr/rtmgr-cfg.yaml")

	//Invalid JSON file
	_, err = GetPlatformComponents("./pkg/rtmg/rtmgr.go")
}
