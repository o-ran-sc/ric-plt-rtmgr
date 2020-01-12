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
  Mnemonic:	sdl/types.go
  Abstract:	Contains SDL (Shared Data Layer) specific types
  Date:		16 March 2019
*/
package sdl

import "routing-manager/pkg/rtmgr"
import "routing-manager/pkg/models" 

//type readAll func(string) (*rtmgr.RicComponents, error)
//type writeAll func(string, *rtmgr.RicComponents) error

type EngineConfig struct {
	Name        string
	Version     string
	Protocol    string
	Instance    Engine
	IsAvailable bool
}

type Engine interface {
	ReadAll(string) (*rtmgr.RicComponents, error)
	WriteAll(string, *rtmgr.RicComponents) error
	WriteXApps(string, *[]rtmgr.XApp) error
	WriteNewE2TInstance(string, *rtmgr.E2TInstance,string) error
	WriteAssRANToE2TInstance(string, models.RanE2tMap) error
	WriteDisAssRANFromE2TInstance(string, models.RanE2tMap) error
	WriteDeleteE2TInstance(string, *models.E2tDeleteData) error
}
