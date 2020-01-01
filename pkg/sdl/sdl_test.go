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
	Mnemonic:	sbi_test.go
	Abstract:
	Date:		25 April 2019
*/
package sdl

import (
	"routing-manager/pkg/stub"
	"github.com/go-openapi/swag"
	"routing-manager/pkg/models"
	"routing-manager/pkg/rtmgr"
	"testing"
)

/*
RmrPub.GeneratePolicies() method is tested for happy path case
*/
func TestFileWriteAll(t *testing.T) {
	var err error
	var file = File{}

	err = file.WriteAll("ut.rt", &stub.ValidRicComponents)
	t.Log(err)
	/* This test is for empty file */
	err = file.WriteAll("", &stub.ValidRicComponents)
	t.Log(err)
}

/*
RmrPush.GeneratePolicies() method is tested for happy path case
*/
func TestFileReadAll(t *testing.T) {
	var err error
	var file = File{}

	data, err := file.ReadAll("ut.rt")
	t.Log(data)
	t.Log(err)
	/* Test to read a Directory */
	data, err = file.ReadAll("/tmp")
	t.Log(data)
	t.Log(err)
}

/*
WriteXApps
*/
func TestFileWriteXApps(t *testing.T) {
	var err error
	var file = File{}

	err = file.WriteXApps("ut.rt", stub.ValidXApps)
	t.Log(err)
	/*Write data to a file that doesn't exist */
	err = file.WriteXApps("ut.rtx", stub.ValidXApps)
	t.Log(err)

}

/*
GetSdl instance with correct and incorrect arguments
*/
func TestFileGetSdl(t *testing.T) {
	var err error
	_, err = GetSdl("")
	t.Log(err)
	_, err = GetSdl("file")
	t.Log(err)
}

/*
WriteNewE2TInstance
*/
func TestWriteNewE2TInstance(t *testing.T) {
	var err error
	var file = File{}
	file.WriteNewE2TInstance("", &stub.ValidE2TInstance)
	t.Log(err)
	file.WriteNewE2TInstance("ut.rt", &stub.ValidE2TInstance)
	t.Log(err)
}

/*
WriteAssRANToE2TInstance
*/
func TestWriteAssRANToE2TInstance(t *testing.T) {
	var err error
        var file = File{}
	// File is not provided as argument
	file.WriteAssRANToE2TInstance("",stub.Rane2tmap)
	t.Log(err)
	file.WriteNewE2TInstance("ut.rt", &stub.ValidE2TInstance)
	file.WriteAssRANToE2TInstance("ut.rt",stub.Rane2tmap)
	t.Log(err)
}

/*
WriteDisAssRANFromE2TInstance 
*/
func TestWriteDisAssRANFromE2TInstance(t *testing.T) {
	var err error
        var file = File{}
	// File is not provided as argument
	file.WriteDisAssRANFromE2TInstance("",stub.Rane2tmap)
	t.Log(err)
	//RAN list is empty
	file.WriteNewE2TInstance("ut.rt", &stub.ValidE2TInstance)
        file.WriteAssRANToE2TInstance("ut.rt",stub.Rane2tmap)
	file.WriteDisAssRANFromE2TInstance("ut.rt",stub.Rane2tmaponlyE2t)
	//RAN list is present
	file.WriteNewE2TInstance("ut.rt", &stub.ValidE2TInstance)
        file.WriteAssRANToE2TInstance("ut.rt",stub.Rane2tmap)
	file.WriteDisAssRANFromE2TInstance("ut.rt",stub.Rane2tmap)
	t.Log(err)
}

/*
WriteDeleteE2TInstance E2TInst *models.E2tDeleteData) error
*/
func TestWriteDeleteE2TInstance(t *testing.T) {
	var err error
        var file = File{}
	e2deldata := &models.E2tDeleteData{}
	// File is not provided as argument
	file.WriteDeleteE2TInstance("",e2deldata)
	//Delete E2t Instance,associate new rans and dissociate some rans
	file.WriteNewE2TInstance("ut.rt", &rtmgr.E2TInstance{
		Name:    "E2Tinstance1",
		Fqdn:    "10.10.10.10:100",
		Ranlist: []string{"1", "2"},
			},
		)
	file.WriteNewE2TInstance("ut.rt", &rtmgr.E2TInstance{
		Name:    "E2Tinstance2",
		Fqdn:    "11.11.11.11:100",
		Ranlist: []string{"3", "4"},
			},
		)
	file.WriteDeleteE2TInstance("ut.rt",&models.E2tDeleteData{
		E2TAddress: swag.String("10.10.10.10:100"),
		RanAssocList: models.RanE2tMap{ 
				{E2TAddress: swag.String("11.11.11.11:100"),RanNamelist: []string{"5","6"}},
				{E2TAddress: swag.String("doesntexist"),RanNamelist: []string{}}, },
			})
	t.Log(err)

}
