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
  Mnemonic:	nbi.go
  Abstract:	Contains NBI (NorthBound Interface) module definitions and generic NBI components
  Date:		12 March 2019
*/

package nbi

import (
	"errors"
	"fmt"
	"rtmgr"
)

var (
	SupportedNbis = []*NbiEngineConfig{
		&NbiEngineConfig{
			NbiEngine{
				Name:     "httpGetter",
				Version:  "v1",
				Protocol: "http",
			},
			batchFetch(fetchXappList),
			true,
		},
		&NbiEngineConfig{
			NbiEngine{
				Name:     "httpRESTful",
				Version:  "v1",
				Protocol: "http",
			},
			batchFetch(nil),
			false,
		},
		&NbiEngineConfig{
			NbiEngine{
				Name:     "gRPC",
				Version:  "v1",
				Protocol: "http2",
			},
			batchFetch(nil),
			false,
		},
	}
)

func ListNbis() {
	fmt.Printf("NBI:\n")
	for _, nbi := range SupportedNbis {
		if nbi.IsAvailable {
			rtmgr.Logger.Info(nbi.Engine.Name + "/" + nbi.Engine.Version)
		}
	}
}

func GetNbi(nbiName string) (*NbiEngineConfig, error) {
	for _, nbi := range SupportedNbis {
		if nbi.Engine.Name == nbiName && nbi.IsAvailable {
			return nbi, nil
		}
	}
	return nil, errors.New("NBI:" + nbiName + " is not supported or still not a available")
}
