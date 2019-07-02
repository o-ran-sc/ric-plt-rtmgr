// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeployXappParams creates a new DeployXappParams object
// with the default values initialized.
func NewDeployXappParams() *DeployXappParams {
	var ()
	return &DeployXappParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeployXappParamsWithTimeout creates a new DeployXappParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeployXappParamsWithTimeout(timeout time.Duration) *DeployXappParams {
	var ()
	return &DeployXappParams{

		timeout: timeout,
	}
}

// NewDeployXappParamsWithContext creates a new DeployXappParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeployXappParamsWithContext(ctx context.Context) *DeployXappParams {
	var ()
	return &DeployXappParams{

		Context: ctx,
	}
}

// NewDeployXappParamsWithHTTPClient creates a new DeployXappParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeployXappParamsWithHTTPClient(client *http.Client) *DeployXappParams {
	var ()
	return &DeployXappParams{
		HTTPClient: client,
	}
}

/*DeployXappParams contains all the parameters to send to the API endpoint
for the deploy xapp operation typically these are written to a http.Request
*/
type DeployXappParams struct {

	/*XAppInfo
	  xApp information

	*/
	XAppInfo DeployXappBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the deploy xapp params
func (o *DeployXappParams) WithTimeout(timeout time.Duration) *DeployXappParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the deploy xapp params
func (o *DeployXappParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the deploy xapp params
func (o *DeployXappParams) WithContext(ctx context.Context) *DeployXappParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the deploy xapp params
func (o *DeployXappParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the deploy xapp params
func (o *DeployXappParams) WithHTTPClient(client *http.Client) *DeployXappParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the deploy xapp params
func (o *DeployXappParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXAppInfo adds the xAppInfo to the deploy xapp params
func (o *DeployXappParams) WithXAppInfo(xAppInfo DeployXappBody) *DeployXappParams {
	o.SetXAppInfo(xAppInfo)
	return o
}

// SetXAppInfo adds the xAppInfo to the deploy xapp params
func (o *DeployXappParams) SetXAppInfo(xAppInfo DeployXappBody) {
	o.XAppInfo = xAppInfo
}

// WriteToRequest writes these params to a swagger request
func (o *DeployXappParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.XAppInfo); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}