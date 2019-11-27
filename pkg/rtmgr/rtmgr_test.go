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

func TestSetLogLevel(t *testing.T) {
	modeIsOk := []string{"info", "warn", "debug", "error"}
	modeOsNok := []string{"inValId", "LogLEVEL", "Provided"}
	for _, value := range modeIsOk {
		if SetLogLevel(value) != nil {
			t.Error("Invalid log level: " + value)
		}
	}

	for _, value := range modeOsNok {
		if SetLogLevel(value) == nil {
			t.Error("Invalid log level: " + value)
		}
	}
}
