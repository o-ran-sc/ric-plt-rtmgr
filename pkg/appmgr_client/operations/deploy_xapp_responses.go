// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"

	appmgr_model "routing-manager/pkg/appmgr_model"
)

// DeployXappReader is a Reader for the DeployXapp structure.
type DeployXappReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeployXappReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewDeployXappCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewDeployXappBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewDeployXappInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDeployXappCreated creates a DeployXappCreated with default headers values
func NewDeployXappCreated() *DeployXappCreated {
	return &DeployXappCreated{}
}

/*DeployXappCreated handles this case with default header values.

xApp successfully created
*/
type DeployXappCreated struct {
	Payload *appmgr_model.Xapp
}

func (o *DeployXappCreated) Error() string {
	return fmt.Sprintf("[POST /xapps][%d] deployXappCreated  %+v", 201, o.Payload)
}

func (o *DeployXappCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(appmgr_model.Xapp)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeployXappBadRequest creates a DeployXappBadRequest with default headers values
func NewDeployXappBadRequest() *DeployXappBadRequest {
	return &DeployXappBadRequest{}
}

/*DeployXappBadRequest handles this case with default header values.

Invalid input
*/
type DeployXappBadRequest struct {
}

func (o *DeployXappBadRequest) Error() string {
	return fmt.Sprintf("[POST /xapps][%d] deployXappBadRequest ", 400)
}

func (o *DeployXappBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeployXappInternalServerError creates a DeployXappInternalServerError with default headers values
func NewDeployXappInternalServerError() *DeployXappInternalServerError {
	return &DeployXappInternalServerError{}
}

/*DeployXappInternalServerError handles this case with default header values.

Internal error
*/
type DeployXappInternalServerError struct {
}

func (o *DeployXappInternalServerError) Error() string {
	return fmt.Sprintf("[POST /xapps][%d] deployXappInternalServerError ", 500)
}

func (o *DeployXappInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*DeployXappBody deploy xapp body
swagger:model DeployXappBody
*/
type DeployXappBody struct {

	// Name of the xApp
	// Required: true
	XAppName *string `json:"xAppName"`
}

// Validate validates this deploy xapp body
func (o *DeployXappBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateXAppName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *DeployXappBody) validateXAppName(formats strfmt.Registry) error {

	if err := validate.Required("xAppInfo"+"."+"xAppName", "body", o.XAppName); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *DeployXappBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *DeployXappBody) UnmarshalBinary(b []byte) error {
	var res DeployXappBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}