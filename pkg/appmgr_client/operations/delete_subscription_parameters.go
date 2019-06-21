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
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteSubscriptionParams creates a new DeleteSubscriptionParams object
// with the default values initialized.
func NewDeleteSubscriptionParams() *DeleteSubscriptionParams {
	var ()
	return &DeleteSubscriptionParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteSubscriptionParamsWithTimeout creates a new DeleteSubscriptionParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteSubscriptionParamsWithTimeout(timeout time.Duration) *DeleteSubscriptionParams {
	var ()
	return &DeleteSubscriptionParams{

		timeout: timeout,
	}
}

// NewDeleteSubscriptionParamsWithContext creates a new DeleteSubscriptionParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteSubscriptionParamsWithContext(ctx context.Context) *DeleteSubscriptionParams {
	var ()
	return &DeleteSubscriptionParams{

		Context: ctx,
	}
}

// NewDeleteSubscriptionParamsWithHTTPClient creates a new DeleteSubscriptionParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteSubscriptionParamsWithHTTPClient(client *http.Client) *DeleteSubscriptionParams {
	var ()
	return &DeleteSubscriptionParams{
		HTTPClient: client,
	}
}

/*DeleteSubscriptionParams contains all the parameters to send to the API endpoint
for the delete subscription operation typically these are written to a http.Request
*/
type DeleteSubscriptionParams struct {

	/*SubscriptionID
	  ID of subscription

	*/
	SubscriptionID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete subscription params
func (o *DeleteSubscriptionParams) WithTimeout(timeout time.Duration) *DeleteSubscriptionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete subscription params
func (o *DeleteSubscriptionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete subscription params
func (o *DeleteSubscriptionParams) WithContext(ctx context.Context) *DeleteSubscriptionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete subscription params
func (o *DeleteSubscriptionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete subscription params
func (o *DeleteSubscriptionParams) WithHTTPClient(client *http.Client) *DeleteSubscriptionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete subscription params
func (o *DeleteSubscriptionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSubscriptionID adds the subscriptionID to the delete subscription params
func (o *DeleteSubscriptionParams) WithSubscriptionID(subscriptionID int64) *DeleteSubscriptionParams {
	o.SetSubscriptionID(subscriptionID)
	return o
}

// SetSubscriptionID adds the subscriptionId to the delete subscription params
func (o *DeleteSubscriptionParams) SetSubscriptionID(subscriptionID int64) {
	o.SubscriptionID = subscriptionID
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteSubscriptionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param subscriptionId
	if err := r.SetPathParam("subscriptionId", swag.FormatInt64(o.SubscriptionID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
