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
  Abstract:	Containes NBI (NorthBound Interface) specific types
  Date:		12 March 2019
*/

package nbi

import (
	"routing-manager/pkg/models"
	"routing-manager/pkg/rpe"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/sdl"
)

type FetchAllXappsHandler func(string) (*[]rtmgr.XApp, error)
type RecvXappCallbackDataHandler func(<-chan *models.XappCallbackData) (*[]rtmgr.XApp, error)
type LaunchRestHandler func(*string, chan<- *models.XappCallbackData, chan<- *models.XappSubscriptionData, chan<- *models.XappSubscriptionData)
type ProvideXappHandleHandlerImpl func(chan<- *models.XappCallbackData, *models.XappCallbackData) error
type RetrieveStartupDataHandler func(string, string, string, string, sdl.SdlEngine) error

type NbiEngineConfig struct {
	Name        string
	Version     string
	Protocol    string
	Instance    NbiEngine
	IsAvailable bool
}

type NbiEngine interface {
	Initialize(string, string, string, string, sdl.SdlEngine, rpe.RpeEngine, chan<- bool) error
	Terminate() error
}
