// Code generated by go-swagger; DO NOT EDIT.

package appmgr_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SubscriptionResponse subscription response
// swagger:model subscriptionResponse
type SubscriptionResponse struct {

	// Event which is subscribed
	// Enum: [created deleted updated all]
	EventType string `json:"eventType,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// version
	Version int64 `json:"version,omitempty"`
}

// Validate validates this subscription response
func (m *SubscriptionResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEventType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var subscriptionResponseTypeEventTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["created","deleted","updated","all"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		subscriptionResponseTypeEventTypePropEnum = append(subscriptionResponseTypeEventTypePropEnum, v)
	}
}

const (

	// SubscriptionResponseEventTypeCreated captures enum value "created"
	SubscriptionResponseEventTypeCreated string = "created"

	// SubscriptionResponseEventTypeDeleted captures enum value "deleted"
	SubscriptionResponseEventTypeDeleted string = "deleted"

	// SubscriptionResponseEventTypeUpdated captures enum value "updated"
	SubscriptionResponseEventTypeUpdated string = "updated"

	// SubscriptionResponseEventTypeAll captures enum value "all"
	SubscriptionResponseEventTypeAll string = "all"
)

// prop value enum
func (m *SubscriptionResponse) validateEventTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, subscriptionResponseTypeEventTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *SubscriptionResponse) validateEventType(formats strfmt.Registry) error {

	if swag.IsZero(m.EventType) { // not required
		return nil
	}

	// value enum
	if err := m.validateEventTypeEnum("eventType", "body", m.EventType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SubscriptionResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriptionResponse) UnmarshalBinary(b []byte) error {
	var res SubscriptionResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
