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
package sbi

import (
	"testing"
	"reflect"
	"errors"
)

func TestGetSbi(t *testing.T) {
	var errtype = errors.New("")
	var sbitype = new(NngPub)
	var invalids = []string{"nngpus", ""}

	sbii, err := GetSbi("nngpub")
	if err != nil {
		t.Errorf("GetSbi(nngpub) was incorrect, got: %v, want: %v.", reflect.TypeOf(err), nil)
	}
	if reflect.TypeOf(sbii) != reflect.TypeOf(sbitype) {
		t.Errorf("GetSbi(nngpub) was incorrect, got: %v, want: %v.", reflect.TypeOf(sbii), reflect.TypeOf(sbitype))
	}

	for _, arg := range invalids {
		_, err := GetSbi(arg)
		if err == nil {
			t.Errorf("GetSbi("+arg+") was incorrect, got: %v, want: %v.", reflect.TypeOf(err), reflect.TypeOf(errtype))
		}
	}
}
