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
  Mnemonic:	nbi.go
  Abstract:	Contains NBI (NorthBound Interface) specific types
  Date:		12 March 2019
*/

package nbi

import (
	"routing-manager/pkg/models"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sdl"
	"sync"
)

type FetchAllXAppsHandler func(string) (*[]rtmgr.XApp, error)
type RecvXappCallbackDataHandler func(<-chan *models.XappCallbackData) (*[]rtmgr.XApp, error)
type RecvNewE2TdataHandler func(<-chan *models.E2tData) (*rtmgr.E2TInstance, string, error)
type LaunchRestHandler func(*string)

//type ProvideXappHandleHandlerImpl func(chan<- *models.XappCallbackData, *models.XappCallbackData) error
type RetrieveStartupDataHandler func(string, string, string, string, string, sdl.Engine) error

type EngineConfig struct {
	Name        string
	Version     string
	Protocol    string
	Instance    Engine
	IsAvailable bool
}

type Engine interface {
	Initialize(string, string, string, string, string, sdl.Engine, rpe.Engine, *sync.Mutex) error
	Terminate() error
}
