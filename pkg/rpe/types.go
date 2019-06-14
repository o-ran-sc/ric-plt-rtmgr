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
  Mnemonic:	rpe/types.go
  Abstract:	Containes RPE (Route Policy Engine) specific types
  Date:		12 March 2019
*/

package rpe

import "routing-manager/pkg/rtmgr"

type generatePolicies func(rtmgr.Endpoints) *[]string

type RpeEngineConfig  struct {
	Name     string
	Version  string
  Protocol string
  Instance RpeEngine
	IsAvailable bool
}

type RpeEngine interface {
  GeneratePolicies(rtmgr.Endpoints) *[]string
}
