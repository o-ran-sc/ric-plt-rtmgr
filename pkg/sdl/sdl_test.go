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
	Mnemonic:	sbi_test.go
	Abstract:
	Date:		25 April 2019
*/
package sdl

import (
	"routing-manager/pkg/stub"
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
}
