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
  Mnemonic:	sbi/types.go
  Abstract:	Containes SBI (SouthBound Interface) specific types
  Date:		16 March 2019
*/

package sbi

import "rtmgr"

type distributeAll func(*[]string) error
type openSocket func(string) error
type closeSocket func() error
type createEndpointSocket func(*rtmgr.Endpoint) error
type destroyEndpointSocket func(*rtmgr.Endpoint) error


type SbiEngine struct {
	Name     string
	Version  string
	Protocol string
}

type SbiEngineConfig struct {
	Engine        SbiEngine
	OpenSocket    openSocket
	CloseSocket   closeSocket
	CreateEndpointSocket createEndpointSocket
	DestroyEndpointSocket destroyEndpointSocket
	DistributeAll distributeAll
	IsAvailable   bool
}
