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

// Subscription subscription
// swagger:model subscription
type Subscription struct {

	// Event which is subscribed
	// Enum: [created deleted updated all]
	EventType string `json:"eventType,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// Maximum number of retries
	MaxRetries int64 `json:"maxRetries,omitempty"`

	// Time in seconds to wait before next retry
	RetryTimer int64 `json:"retryTimer,omitempty"`

	// target Url
	TargetURL string `json:"targetUrl,omitempty"`
}

// Validate validates this subscription
func (m *Subscription) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEventType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var subscriptionTypeEventTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["created","deleted","updated","all"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		subscriptionTypeEventTypePropEnum = append(subscriptionTypeEventTypePropEnum, v)
	}
}

const (

	// SubscriptionEventTypeCreated captures enum value "created"
	SubscriptionEventTypeCreated string = "created"

	// SubscriptionEventTypeDeleted captures enum value "deleted"
	SubscriptionEventTypeDeleted string = "deleted"

	// SubscriptionEventTypeUpdated captures enum value "updated"
	SubscriptionEventTypeUpdated string = "updated"

	// SubscriptionEventTypeAll captures enum value "all"
	SubscriptionEventTypeAll string = "all"
)

// prop value enum
func (m *Subscription) validateEventTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, subscriptionTypeEventTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Subscription) validateEventType(formats strfmt.Registry) error {

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
func (m *Subscription) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Subscription) UnmarshalBinary(b []byte) error {
	var res Subscription
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}