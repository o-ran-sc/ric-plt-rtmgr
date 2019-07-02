// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	appmgr_model "routing-manager/pkg/appmgr_model"
)

// GetSubscriptionsReader is a Reader for the GetSubscriptions structure.
type GetSubscriptionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSubscriptionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetSubscriptionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetSubscriptionsOK creates a GetSubscriptionsOK with default headers values
func NewGetSubscriptionsOK() *GetSubscriptionsOK {
	return &GetSubscriptionsOK{}
}

/*GetSubscriptionsOK handles this case with default header values.

successful query of subscriptions
*/
type GetSubscriptionsOK struct {
	Payload appmgr_model.AllSubscriptions
}

func (o *GetSubscriptionsOK) Error() string {
	return fmt.Sprintf("[GET /subscriptions][%d] getSubscriptionsOK  %+v", 200, o.Payload)
}

func (o *GetSubscriptionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}