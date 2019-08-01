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
  Mnemonic:	nbi.go
  Abstract:	Contains NBI (NorthBound Interface) module definitions and generic NBI components
  Date:		12 March 2019
*/

package nbi

import (
	"errors"
	"routing-manager/pkg/rtmgr"
        "net/url"
        apiclient "routing-manager/pkg/appmgr_client"
        "routing-manager/pkg/appmgr_client/operations"
        "routing-manager/pkg/appmgr_model"
        httptransport "github.com/go-openapi/runtime/client"
        "github.com/go-openapi/strfmt"
        "github.com/go-openapi/swag"
        "time"

)

var (
	SupportedNbis = []*NbiEngineConfig{
		&NbiEngineConfig{
			Name:     "httpGetter",
			Version:  "v1",
			Protocol: "http",
			Instance: NewHttpGetter(),
			IsAvailable: true,
		},
		&NbiEngineConfig{
			Name:     "httpRESTful",
			Version:  "v1",
			Protocol: "http",
			Instance: NewHttpRestful(),
			IsAvailable: true,
		},
	}
)

type Nbi struct {

}

func GetNbi(nbiName string) (NbiEngine, error) {
	for _, nbi := range SupportedNbis {
		if nbi.Name == nbiName && nbi.IsAvailable {
			return nbi.Instance, nil
		}
	}
	return nil, errors.New("NBI:" + nbiName + " is not supported or still not a available")
}

func CreateSubReq(restUrl string, restPort string) *appmgr_model.SubscriptionRequest {
	// TODO: parametize function
        subReq := appmgr_model.SubscriptionRequest{
                TargetURL:  swag.String(restUrl + ":" + restPort + "/ric/v1/handles/xapp-handle/"),
                EventType:  swag.String("all"),
                MaxRetries: swag.Int64(5),
                RetryTimer: swag.Int64(10),
        }

        return &subReq
}

func PostSubReq(xmUrl string, nbiif string) error {
        // setting up POST request to Xapp Manager
        appmgrUrl, err := url.Parse(xmUrl)
        if err != nil {
                rtmgr.Logger.Error("Invalid XApp manager url/hostname: " + err.Error())
                return err
        }
	nbiifUrl, err := url.Parse(nbiif)
	if err != nil {
		rtmgr.Logger.Error("Invalid NBI address/port: " + err.Error())
		return err
	}
        transport := httptransport.New(appmgrUrl.Hostname()+":"+appmgrUrl.Port(), "/ric/v1", []string{"http"})
        client := apiclient.New(transport, strfmt.Default)
        addSubParams := operations.NewAddSubscriptionParamsWithTimeout(10 * time.Second)
        // create sub req with rest url and port
        subReq := CreateSubReq(string(nbiifUrl.Scheme+"://"+nbiifUrl.Hostname()), nbiifUrl.Port())
        resp, postErr := client.Operations.AddSubscription(addSubParams.WithSubscriptionRequest(subReq))
        if postErr != nil {
                rtmgr.Logger.Error("POST unsuccessful:"+postErr.Error())
                return postErr
        } else {
                // TODO: use the received ID
                rtmgr.Logger.Info("POST received: "+string(resp.Payload.ID))
                return nil
        }
}

