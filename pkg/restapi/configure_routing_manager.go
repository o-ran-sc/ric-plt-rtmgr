// This file is safe to edit. Once it exists it will not be overwritten

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

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"routing-manager/pkg/restapi/operations"
	"routing-manager/pkg/restapi/operations/handle"
	"routing-manager/pkg/restapi/operations/health"
)

//go:generate swagger generate server --target ../../pkg --name RoutingManager --spec ../../api/routing_manager.yaml --exclude-main

func configureFlags(api *operations.RoutingManagerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.RoutingManagerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.HandleGetHandlesHandler = handle.GetHandlesHandlerFunc(func(params handle.GetHandlesParams) middleware.Responder {
		return middleware.NotImplemented("operation handle.GetHandles has not yet been implemented")
	})
	api.HealthGetHealthHandler = health.GetHealthHandlerFunc(func(params health.GetHealthParams) middleware.Responder {
		return middleware.NotImplemented("operation health.GetHealth has not yet been implemented")
	})
	api.HandleProvideXappHandleHandler = handle.ProvideXappHandleHandlerFunc(func(params handle.ProvideXappHandleParams) middleware.Responder {
		return middleware.NotImplemented("operation handle.ProvideXappHandle has not yet been implemented")
	})
	api.HandleProvideXappSubscriptionHandleHandler = handle.ProvideXappSubscriptionHandleHandlerFunc(func(params handle.ProvideXappSubscriptionHandleParams) middleware.Responder {
		return middleware.NotImplemented("operation handle.ProvideXappSubscriptionHandle has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
