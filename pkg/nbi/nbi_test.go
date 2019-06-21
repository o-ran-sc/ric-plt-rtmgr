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
  Mnemonic:     nbi_test.go
  Abstract:     NBI unit tests
  Date:         21 May 2019
*/

package nbi

import (
	"testing"
	"reflect"
	"errors"
	"routing-manager/pkg/appmgr_model"
	"github.com/go-openapi/swag"
	"net"
	"net/http"
	"net/http/httptest"
)

func TestGetNbi(t *testing.T) {
	var errtype = errors.New("")
	var nbitype = new(HttpGetter)
	var invalids = []string{"httpgetter", ""}

	nbii, err := GetNbi("httpGetter")
	if err != nil {
		t.Errorf("GetNbi(HttpGetter) was incorrect, got: %v, want: %v.", reflect.TypeOf(err), nil)
	}
	if reflect.TypeOf(nbii) != reflect.TypeOf(nbitype) {
		t.Errorf("GetNbi(HttpGetter) was incorrect, got: %v, want: %v.", reflect.TypeOf(nbii), reflect.TypeOf(nbitype))
	}

	for _, arg := range invalids {
		_, err := GetNbi(arg)
		if err == nil {
			t.Errorf("GetNbi("+arg+") was incorrect, got: %v, want: %v.", reflect.TypeOf(err), reflect.TypeOf(errtype))
		}
	}
}

func TestCreateSubReq(t *testing.T) {
	var subReq = appmgr_model.SubscriptionRequest{
                TargetURL:  swag.String("localhost:8000/ric/v1/handles/xapp-handle/"),
                EventType:  swag.String("all"),
                MaxRetries: swag.Int64(5),
                RetryTimer: swag.Int64(10),
        }
	subReq2 := CreateSubReq("localhost","8000")
	if reflect.TypeOf(subReq) != reflect.TypeOf(*subReq2) {
		t.Errorf("Invalid type, got: %v, want: %v.", reflect.TypeOf(subReq), reflect.TypeOf(*subReq2))
	}
	if *(subReq.TargetURL) != *(subReq2.TargetURL) {
		t.Errorf("Invalid TargetURL generated, got %v, want %v", *subReq.TargetURL, *subReq2.TargetURL)
	}
        if *(subReq.EventType) != *(subReq2.EventType) {
                t.Errorf("Invalid EventType generated, got %v, want %v", *subReq.EventType, *subReq2.EventType)
        }
        if *(subReq.MaxRetries) != *(subReq2.MaxRetries) {
                t.Errorf("Invalid MaxRetries generated, got %v, want %v", *subReq.MaxRetries, *subReq2.MaxRetries)
        }
        if *(subReq.RetryTimer) != *(subReq2.RetryTimer) {
                t.Errorf("Invalid RetryTimer generated, got %v, want %v", *subReq.RetryTimer, *subReq2.RetryTimer)
        }
}

func TestPostSubReq(t *testing.T) {
	b := []byte(`{"ID":"deadbeef1234567890", "Version":0, "EventType":"all"}`)
	l, err := net.Listen("tcp", "127.0.0.1:3000")
        if err != nil {
                t.Error("Failed to create listener: " + err.Error())
        }
        ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                t.Log(r.Method)
                t.Log(r.URL)
                if r.Method == "POST" && r.URL.String() == "/ric/v1/subscriptions" {
                        t.Log("Sending reply")
                        w.Header().Add("Content-Type", "application/json")
                        w.WriteHeader(http.StatusOK)
                        w.Write(b)
                }
        }))
        ts.Listener.Close()
        ts.Listener = l

        ts.Start()
	defer ts.Close()
	//subReq := CreateSubReq("localhost","8888")
	err = PostSubReq("http://127.0.0.1:3000/","http://localhost:8888")
	if err != nil {
		t.Error("Error occured: "+err.Error())
	}
}

func TestPostSubReqWithInvalidUrls(t *testing.T) {
	// invalid Xapp Manager URL
        err := PostSubReq("http://127.0","http://localhost:8888")
        if err == nil {
                t.Error("Error occured: "+err.Error())
        }
	// invalid rest api url
	err = PostSubReq("http://127.0.0.1:3000/","localhost:8888")
	if err == nil {
		t.Error("Error occured: "+err.Error())
	}
}

