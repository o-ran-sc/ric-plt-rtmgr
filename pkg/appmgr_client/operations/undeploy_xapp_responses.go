// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// UndeployXappReader is a Reader for the UndeployXapp structure.
type UndeployXappReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UndeployXappReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 204:
		result := NewUndeployXappNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewUndeployXappBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewUndeployXappInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUndeployXappNoContent creates a UndeployXappNoContent with default headers values
func NewUndeployXappNoContent() *UndeployXappNoContent {
	return &UndeployXappNoContent{}
}

/*UndeployXappNoContent handles this case with default header values.

Successful deletion of xApp
*/
type UndeployXappNoContent struct {
}

func (o *UndeployXappNoContent) Error() string {
	return fmt.Sprintf("[DELETE /xapps/{xAppName}][%d] undeployXappNoContent ", 204)
}

func (o *UndeployXappNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUndeployXappBadRequest creates a UndeployXappBadRequest with default headers values
func NewUndeployXappBadRequest() *UndeployXappBadRequest {
	return &UndeployXappBadRequest{}
}

/*UndeployXappBadRequest handles this case with default header values.

Invalid xApp name supplied
*/
type UndeployXappBadRequest struct {
}

func (o *UndeployXappBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /xapps/{xAppName}][%d] undeployXappBadRequest ", 400)
}

func (o *UndeployXappBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUndeployXappInternalServerError creates a UndeployXappInternalServerError with default headers values
func NewUndeployXappInternalServerError() *UndeployXappInternalServerError {
	return &UndeployXappInternalServerError{}
}

/*UndeployXappInternalServerError handles this case with default header values.

Internal error
*/
type UndeployXappInternalServerError struct {
}

func (o *UndeployXappInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /xapps/{xAppName}][%d] undeployXappInternalServerError ", 500)
}

func (o *UndeployXappInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
