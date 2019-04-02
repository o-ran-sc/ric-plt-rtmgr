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
  Mnemonic:	file.go
  Abstract:	File SDL implementation. Only for testing purpose.
  Date:		16 March 2019
*/

package sdl

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"rtmgr"
)

/*
Reads the content of the rt.json file
Parses the JSON content and loads each xApp entry into an xApp object
Returns an array os xApp object
*/
func fileReadAll(file string) (*[]rtmgr.XApp, error) {
	rtmgr.Logger.Debug("Invoked file.fileReadAll")
	rtmgr.Logger.Debug("file.fileReadAll opens file: " + file)
	var xapps *[]rtmgr.XApp
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, errors.New("cannot open the file due to: " + err.Error())
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.New("cannot read the file due to: " + err.Error())
	}
	err = json.Unmarshal(byteValue, &xapps)
	if err != nil {
		return nil, errors.New("cannot parse data due to: " + err.Error())
	}
	rtmgr.Logger.Debug("file.fileReadAll returns: %v", xapps)
	return xapps, nil
}

func fileWriteAll(file string, xapps *[]rtmgr.XApp) error {
	rtmgr.Logger.Debug("Invoked file.fileWriteAll")
	rtmgr.Logger.Debug("file.fileWriteAll writes into file: " + file)
	rtmgr.Logger.Debug("file.fileWriteAll writes data: %v", (*xapps))
	byteValue, err := json.Marshal(xapps)
	if err != nil {
		return errors.New("cannot convert data due to: " + err.Error())
	}
	err = ioutil.WriteFile(file, byteValue, 0644)
	if err != nil {
		return errors.New("cannot write file due to: " + err.Error())
	}
	return nil
}
