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
	"strings"
	"routing-manager/pkg/rtmgr"
	"routing-manager/pkg/models"
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

func (f *File) WriteNewE2TInstance(file string, E2TInst *rtmgr.E2TInstance,meiddata string) error {
        xapp.Logger.Debug("Invoked sdl.WriteNewE2TInstance")
        xapp.Logger.Debug("file.WriteNewE2TInstance writes into file: " + file)
        xapp.Logger.Debug("file.WriteNewE2TInstance writes data: %v", *E2TInst)

        ricData, err := NewFile().ReadAll(file)
        if err != nil {
                xapp.Logger.Error("cannot get data from sdl interface due to: " + err.Error())
                return errors.New("cannot read full ric data to modify xApps data, due to:  " + err.Error())
        }
        ricData.E2Ts[E2TInst.Fqdn] = *E2TInst
	if (len(meiddata) > 0){
	    ricData.MeidMap = []string {meiddata}
        } else {
	    ricData.MeidMap = []string {}
	}


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

func (f *File) WriteAssRANToE2TInstance(file string, rane2tmap models.RanE2tMap) error {
        xapp.Logger.Debug("Invoked sdl.WriteAssRANToE2TInstance")
        xapp.Logger.Debug("file.WriteAssRANToE2TInstance writes into file: " + file)
        xapp.Logger.Debug("file.WriteAssRANToE2TInstance writes data: %v", rane2tmap)

        ricData, err := NewFile().ReadAll(file)
        if err != nil {
                xapp.Logger.Error("cannot get data from sdl interface due to: " + err.Error())
                return errors.New("cannot read full ric data to modify xApps data, due to:  " + err.Error())
        }

	ricData.MeidMap = []string{}
	for _, element := range rane2tmap {
		xapp.Logger.Info("data received")
		var str,meidar string
		for _, meid := range element.RanNamelist {
		    meidar += meid + " "
		}
		str = "mme_ar|" + *element.E2TAddress + "|" + strings.TrimSuffix(meidar," ")
		ricData.MeidMap = append(ricData.MeidMap,str)

		for key, _ := range ricData.E2Ts {
			if key == *element.E2TAddress {
				var estObj rtmgr.E2TInstance
				estObj = ricData.E2Ts[key]
				estObj.Ranlist = append(ricData.E2Ts[key].Ranlist, element.RanNamelist...)
				ricData.E2Ts[key]= estObj
			}
		}
	}

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

func (f *File) WriteDisAssRANFromE2TInstance(file string, disassranmap models.RanE2tMap) error {
        xapp.Logger.Debug("Invoked sdl.WriteDisAssRANFromE2TInstance")
        xapp.Logger.Debug("file.WriteDisAssRANFromE2TInstance writes into file: " + file)
        xapp.Logger.Debug("file.WriteDisAssRANFromE2TInstance writes data: %v", disassranmap)

        ricData, err := NewFile().ReadAll(file)
        if err != nil {
                xapp.Logger.Error("cannot get data from sdl interface due to: " + err.Error())
                return errors.New("cannot read full ric data to modify xApps data, due to:  " + err.Error())
        }

	var str,meiddel,meiddisdel string
	ricData.MeidMap = []string{}
	for _, element := range disassranmap {
		xapp.Logger.Info("data received")
		for _, meid := range element.RanNamelist {
		    meiddisdel += meid + " "
		}
		if ( len(element.RanNamelist) > 0 ) {
		    str = "mme_del|" + strings.TrimSuffix(meiddisdel," ")
		    ricData.MeidMap = append(ricData.MeidMap,str)
	        }
		e2taddress_key := *element.E2TAddress
		//Check whether the provided E2T Address is available in SDL as a key. 
		//If exist, proceed further to check RAN list, Otherwise move to next E2T Instance
		if _, exist := ricData.E2Ts[e2taddress_key]; exist {
			var estObj rtmgr.E2TInstance
			estObj = ricData.E2Ts[e2taddress_key]
			// If RAN list is empty, then routing manager assumes that all RANs attached associated to the particular E2T Instance to be removed.
			if len(element.RanNamelist) == 0 {
				xapp.Logger.Debug("RAN List is empty. So disassociating all RANs from the E2T Instance: %v ", *element.E2TAddress)
			for _, meid := range estObj.Ranlist {
			meiddel += meid + " "
			}
			str = "mme_del|" + strings.TrimSuffix(meiddel," ")
			ricData.MeidMap = append(ricData.MeidMap,str)

			estObj.Ranlist = []string{}
			} else {
				xapp.Logger.Debug("Remove only selected rans from E2T Instance: %v and %v ", ricData.E2Ts[e2taddress_key].Ranlist, element.RanNamelist)
				for _, disRanValue := range element.RanNamelist {
					for ranIndex, ranValue := range ricData.E2Ts[e2taddress_key].Ranlist {
						if disRanValue == ranValue {
							estObj.Ranlist[ranIndex] = estObj.Ranlist[len(estObj.Ranlist)-1]
							estObj.Ranlist[len(estObj.Ranlist)-1] = ""
							estObj.Ranlist = estObj.Ranlist[:len(estObj.Ranlist)-1]
						}
					}
				}
			}
			ricData.E2Ts[e2taddress_key]= estObj
		}
	}

	xapp.Logger.Debug("Final data after disassociate: %v", ricData)

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

func (f *File) WriteDeleteE2TInstance(file string, E2TInst *models.E2tDeleteData) error {
	xapp.Logger.Debug("Invoked sdl.WriteDeleteE2TInstance")
	xapp.Logger.Debug("file.WriteDeleteE2TInstance writes into file: " + file)
	xapp.Logger.Debug("file.WriteDeleteE2TInstance writes data: %v", *E2TInst)

	ricData, err := NewFile().ReadAll(file)
	if err != nil {
	        xapp.Logger.Error("cannot get data from sdl interface due to: " + err.Error())
	        return errors.New("cannot read full ric data to modify xApps data, due to:  " + err.Error())
	}


	ricData.MeidMap = []string {}
	var delrow,meiddel string
	if(len(E2TInst.RanNamelistTobeDissociated)>0) {
	    for _, meid := range E2TInst.RanNamelistTobeDissociated {
			meiddel += meid + " "
		}
	    delrow = "mme_del|" + strings.TrimSuffix(meiddel," ")
	    ricData.MeidMap = append(ricData.MeidMap,delrow)
	} else {
	      if(len(ricData.E2Ts[*E2TInst.E2TAddress].Ranlist) > 0) {
	          for _, meid := range ricData.E2Ts[*E2TInst.E2TAddress].Ranlist {
			meiddel += meid + " "
	          }
	          delrow = "mme_del|" + strings.TrimSuffix(meiddel," ")
	          ricData.MeidMap = append(ricData.MeidMap,delrow)
	      }
	}

	delete(ricData.E2Ts, *E2TInst.E2TAddress)

	for _, element := range E2TInst.RanAssocList {
		var str,meidar string
		xapp.Logger.Info("data received")
		for _, meid := range element.RanNamelist {
			meidar = meid + " "
		}
		str = "mme_ar|" + *element.E2TAddress + "|" + strings.TrimSuffix(meidar," ")
		ricData.MeidMap = append(ricData.MeidMap,str)
		key := *element.E2TAddress

		if val, ok := ricData.E2Ts[key]; ok {
			var estObj rtmgr.E2TInstance
			estObj = val
			estObj.Ranlist = append(ricData.E2Ts[key].Ranlist, element.RanNamelist...)
			ricData.E2Ts[key]= estObj
		} else {
			xapp.Logger.Error("file.WriteDeleteE2TInstance E2T instance is not found for provided E2TAddress : %v", errors.New(key).Error())
		}

	}

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
