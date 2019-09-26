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
  Mnemonic:	sdl/types.go
  Abstract:	Containes SDL (Shared Data Layer) specific types
  Date:		16 March 2019
*/
package sdl

import "routing-manager/pkg/rtmgr"

type readAll func(string) (*rtmgr.RicComponents, error)
type writeAll func(string, *rtmgr.RicComponents) error

type SdlEngineConfig struct {
	Name        string
	Version     string
	Protocol    string
	Instance    SdlEngine
	IsAvailable bool
}

type SdlEngine interface {
	ReadAll(string) (*rtmgr.RicComponents, error)
	WriteAll(string, *rtmgr.RicComponents) error
	WriteXapps(string, *[]rtmgr.XApp) error
}
