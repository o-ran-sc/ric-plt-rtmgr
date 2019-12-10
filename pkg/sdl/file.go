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
  Mnemonic:	file.go
  Abstract:	File SDL implementation. Only for testing purpose.
  Date:		16 March 2019
*/

package sdl

import (
	"encoding/json"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"io/ioutil"
	"os"
	"routing-manager/pkg/rtmgr"
)

/*
Reads the content of the rt.json file
Parses the JSON content and loads each xApp entry into an xApp object
Returns an array os xApp object
*/

type File struct {
	Sdl
}

func NewFile() *File {
	instance := new(File)
	return instance
}

func (f *File) ReadAll(file string) (*rtmgr.RicComponents, error) {
	xapp.Logger.Debug("Invoked sdl.ReadAll(" + file + ")")
	var rcs *rtmgr.RicComponents
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, errors.New("cannot open the file due to: " + err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.New("cannot read the file due to: " + err.Error())
	}
	err = json.Unmarshal(byteValue, &rcs)
	if err != nil {
		return nil, errors.New("cannot parse data due to: " + err.Error())
	}
	xapp.Logger.Debug("file.fileReadAll returns: %v", rcs)
	return rcs, nil
}

func (f *File) WriteAll(file string, rcs *rtmgr.RicComponents) error {
	xapp.Logger.Debug("Invoked sdl.WriteAll")
	xapp.Logger.Debug("file.fileWriteAll writes into file: " + file)
	xapp.Logger.Debug("file.fileWriteAll writes data: %v", *rcs)
	byteValue, err := json.Marshal(rcs)
	if err != nil {
		return errors.New("cannot convert data due to: " + err.Error())
	}
	err = ioutil.WriteFile(file, byteValue, 0644)
	if err != nil {
		return errors.New("cannot write file due to: " + err.Error())
	}
	return nil
}

func (f *File) WriteXApps(file string, xApps *[]rtmgr.XApp) error {
	xapp.Logger.Debug("Invoked sdl.WriteXApps")
	xapp.Logger.Debug("file.fileWriteXApps writes into file: " + file)
	xapp.Logger.Debug("file.fileWriteXApps writes data: %v", *xApps)

	ricData, err := NewFile().ReadAll(file)
	if err != nil {
		xapp.Logger.Error("cannot get data from sdl interface due to: " + err.Error())
		return errors.New("cannot read full ric data to modify xApps data, due to:  " + err.Error())
	}
	ricData.XApps = *xApps

	byteValue, err := json.Marshal(ricData)
	if err != nil {
		return errors.New("cannot convert data due to: " + err.Error())
	}
	err = ioutil.WriteFile(file, byteValue, 0644)
	if err != nil {
		return errors.New("cannot write file due to: " + err.Error())
	}
	return nil
}

func (f *File) WriteNewE2TInstance(file string, E2TInst *rtmgr.E2TInstance) error {
        xapp.Logger.Debug("Invoked sdl.WriteNewE2TInstance")
        xapp.Logger.Debug("file.WriteNewE2TInstance writes into file: " + file)
        xapp.Logger.Debug("file.WriteNewE2TInstance writes data: %v", *E2TInst)

        ricData, err := NewFile().ReadAll(file)
        if err != nil {
                xapp.Logger.Error("cannot get data from sdl interface due to: " + err.Error())
                return errors.New("cannot read full ric data to modify xApps data, due to:  " + err.Error())
        }
        ricData.E2Ts[E2TInst.Fqdn] = *E2TInst

        byteValue, err := json.Marshal(ricData)
        if err != nil {
                return errors.New("cannot convert data due to: " + err.Error())
        }
        err = ioutil.WriteFile(file, byteValue, 0644)
        if err != nil {
                return errors.New("cannot write file due to: " + err.Error())
        }
        return nil
}
