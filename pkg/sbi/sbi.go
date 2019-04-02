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
  Mnemonic:	sbi.go
  Abstract:	Contains SBI (SouthBound Interface) module definitions and generic SBI components
  Date:		16 March 2019
*/

package sbi

import (
	"errors"
	"fmt"
	"rtmgr"
)

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
			distributeAll(publishAll),
			true,
		},
		&SbiEngineConfig{
			SbiEngine{
				Name:     "nngpush",
				Version:  "v1",
				Protocol: "nngpipeline",
			},
			openSocket(nil),
			closeSocket(nil),
			distributeAll(nil),
			false,
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
		if sbi.Engine.Name == sbiName && sbi.IsAvailable {
			return sbi, nil
		}
	}
	return nil, errors.New("SBI:" + sbiName + "is not supported or still not a available")
}
