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
  Mnemonic:	nbi.go
  Abstract:	Contains NBI (NorthBound Interface) module definitions and generic NBI components
  Date:		12 March 2019
*/

package nbi

import (
	"errors"
	"net/url"
	apiclient "routing-manager/pkg/appmgr_client"
	"routing-manager/pkg/appmgr_client/operations"
	"routing-manager/pkg/appmgr_model"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

var (
	SupportedNbis = []*EngineConfig{
		{
			Name:        "httpGetter",
			Version:     "v1",
			Protocol:    "http",
			Instance:    NewHttpGetter(),
			IsAvailable: true,
		},
		{
			Name:        "httpRESTful",
			Version:     "v1",
			Protocol:    "http",
			Instance:    NewHttpRestful(),
			IsAvailable: true,
		},
	}
)

type Nbi struct {
}

func GetNbi(nbiName string) (Engine, error) {
	for _, nbi := range SupportedNbis {
		if nbi.Name == nbiName && nbi.IsAvailable {
			return nbi.Instance, nil
		}
	}
	return nil, errors.New("NBI:" + nbiName + " is not supported or still not a available")
}

func CreateSubReq(restUrl string, restPort string) *appmgr_model.SubscriptionRequest {
	// TODO: parameterize function
	subData := appmgr_model.SubscriptionData{
		TargetURL:  swag.String(restUrl + ":" + restPort + "/ric/v1/handles/xapp-handle/"),
		EventType:  appmgr_model.EventTypeAll,
		MaxRetries: swag.Int64(5),
		RetryTimer: swag.Int64(10),
	}

	subReq := appmgr_model.SubscriptionRequest{
		Data: &subData,
	}

	return &subReq
}

func PostSubReq(xmUrl string, nbiif string) error {
	// setting up POST request to Xapp Manager
	appmgrUrl, err := url.Parse(xmUrl)
	if err != nil {
		xapp.Logger.Error("Invalid XApp manager url/hostname: " + err.Error())
		return err
	}
	nbiifUrl, err := url.Parse(nbiif)
	if err != nil {
		xapp.Logger.Error("Invalid NBI address/port: " + err.Error())
		return err
	}
	transport := httptransport.New(appmgrUrl.Hostname()+":"+appmgrUrl.Port(), "/ric/v1", []string{"http"})
	client := apiclient.New(transport, strfmt.Default)
	addSubParams := operations.NewAddSubscriptionParamsWithTimeout(10 * time.Second)
	// create sub req with rest url and port
	subReq := CreateSubReq(nbiifUrl.Scheme+"://"+nbiifUrl.Hostname(), nbiifUrl.Port())
	resp, postErr := client.Operations.AddSubscription(addSubParams.WithSubscriptionRequest(subReq))
	if postErr != nil {
		xapp.Logger.Error("POST unsuccessful:" + postErr.Error())
		return postErr
	} else {
		// TODO: use the received ID
		xapp.Logger.Info("POST received: " + string(resp.Payload.ID))
		return nil
	}
}
