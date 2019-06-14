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
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateRouteHandlerFunc turns a function with the right signature into a update route handler
type UpdateRouteHandlerFunc func(UpdateRouteParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateRouteHandlerFunc) Handle(params UpdateRouteParams) middleware.Responder {
	return fn(params)
}

// UpdateRouteHandler interface for that can handle valid update route params
type UpdateRouteHandler interface {
	Handle(UpdateRouteParams) middleware.Responder
}

// NewUpdateRoute creates a new http.Handler for the update route operation
func NewUpdateRoute(ctx *middleware.Context, handler UpdateRouteHandler) *UpdateRoute {
	return &UpdateRoute{Context: ctx, Handler: handler}
}

/*UpdateRoute swagger:route PUT /routes route updateRoute

Update an existing route

By performing a PUT method on the routes resource, the API caller is able to update an already existing route.

*/
type UpdateRoute struct {
	Context *middleware.Context
	Handler UpdateRouteHandler
}

func (o *UpdateRoute) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateRouteParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
