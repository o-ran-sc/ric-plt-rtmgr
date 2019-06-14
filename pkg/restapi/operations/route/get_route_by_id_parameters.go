// Code generated by go-swagger; DO NOT EDIT.

// ==================================================================================
// Unless otherwise specified, all software contained herein is licensed
// under the Apache License, Version 2.0 (the "Software License");
// you may not use this software except in compliance with the Software
// License. You may obtain a copy of the Software License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the Software License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Software License for the specific language governing permissions
// and limitations under the Software License.
//
// ==================================================================================
//
// Unless otherwise specified, all documentation contained herein is licensed
// under the Creative Commons License, Attribution 4.0 Intl. (the
// "Documentation License"); you may not use this documentation except in
// compliance with the Documentation License. You may obtain a copy of the
// Documentation License at
//
// https://creativecommons.org/licenses/by/4.0/
//
// Unless required by applicable law or agreed to in writing, documentation
// distributed under the Documentation License is distributed on an "AS IS"
// BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the Documentation License for the specific language governing
// permissions and limitations under the Documentation License.
// ==================================================================================
//
//

package route

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetRouteByIDParams creates a new GetRouteByIDParams object
// no default values defined in spec.
func NewGetRouteByIDParams() GetRouteByIDParams {

	return GetRouteByIDParams{}
}

// GetRouteByIDParams contains all the bound params for the get route by id operation
// typically these are obtained from a http.Request
//
// swagger:parameters get_route_by_id
type GetRouteByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*ID of route to return
	  Required: true
	  In: path
	*/
	RouteID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetRouteByIDParams() beforehand.
func (o *GetRouteByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rRouteID, rhkRouteID, _ := route.Params.GetOK("route-id")
	if err := o.bindRouteID(rRouteID, rhkRouteID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindRouteID binds and validates parameter RouteID from path.
func (o *GetRouteByIDParams) bindRouteID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("route-id", "path", "int64", raw)
	}
	o.RouteID = value

	return nil
}
