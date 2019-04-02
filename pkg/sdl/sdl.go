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
  Mnemonic:	sdl.go
  Abstract:	Contains SDL (Shared Data Layer) module definitions and generic SDL components
  Date:		16 March 2019
*/

package sdl

import (
	"errors"
	"fmt"
	"rtmgr"
)

var (
	SupportedSdls = []*SdlEngineConfig{
		&SdlEngineConfig{
			SdlEngine{
				Name:     "file",
				Version:  "v1",
				Protocol: "rawfile",
			},
			readAll(fileReadAll),
			writeAll(fileWriteAll),
			true,
		},
		&SdlEngineConfig{
			SdlEngine{
				Name:     "redis",
				Version:  "v1",
				Protocol: "nsdl",
			},
			readAll(nil),
			writeAll(nil),
			false,
		},
	}
)

func ListSdls() {
	fmt.Printf("SDL:\n")
	for _, sdl := range SupportedSdls {
		if sdl.IsAvailable {
			rtmgr.Logger.Info(sdl.Engine.Name + "/" + sdl.Engine.Version)
		}
	}
}

func GetSdl(sdlName string) (*SdlEngineConfig, error) {
	for _, sdl := range SupportedSdls {
		if sdl.Engine.Name == sdlName && sdl.IsAvailable {
			return sdl, nil
		}
	}
	return nil, errors.New("SDL:" + sdlName + "is not supported or still not a available")
}
