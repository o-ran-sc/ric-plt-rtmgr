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
  Mnemonic:	sbi/types.go
  Abstract:	Contains SBI (SouthBound Interface) specific types
  Date:		16 March 2019
*/

package sbi

import "routing-manager/pkg/rtmgr"

type EngineConfig struct {
	Name        string
	Version     string
	Protocol    string
	Instance    Engine
	IsAvailable bool
}

type Engine interface {
	Initialize(string) error
	Terminate() error
	DistributeAll(*[]string) error
	AddEndpoint(*rtmgr.Endpoint) error
	DeleteEndpoint(*rtmgr.Endpoint) error
	UpdateEndpoints(*rtmgr.RicComponents)
	CreateEndpoint(string) (*rtmgr.Endpoint)
	DistributeToEp(*[]string, *rtmgr.Endpoint) error
}

/*type NngSocket interface {
	Listen(string) error
	Send([]byte) error
	Close() error
	DialOptions(string, map[string]interface{}) error
}

type CreateNewNngSocketHandler func() (NngSocket, error)*/
