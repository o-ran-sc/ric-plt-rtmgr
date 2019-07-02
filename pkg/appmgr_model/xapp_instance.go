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

// XappInstance xapp instance
// swagger:model XappInstance
type XappInstance struct {

	// ip
	IP string `json:"ip,omitempty"`

	// name
	// Required: true
	Name *string `json:"name"`

	// port
	Port int64 `json:"port,omitempty"`

	// rx messages
	RxMessages []string `json:"rxMessages"`

	// xapp instance status
	// Enum: [pending running succeeded failed unknown completed crashLoopBackOff]
	Status string `json:"status,omitempty"`

	// tx messages
	TxMessages []string `json:"txMessages"`
}

// Validate validates this xapp instance
func (m *XappInstance) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *XappInstance) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

var xappInstanceTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["pending","running","succeeded","failed","unknown","completed","crashLoopBackOff"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		xappInstanceTypeStatusPropEnum = append(xappInstanceTypeStatusPropEnum, v)
	}
}

const (

	// XappInstanceStatusPending captures enum value "pending"
	XappInstanceStatusPending string = "pending"

	// XappInstanceStatusRunning captures enum value "running"
	XappInstanceStatusRunning string = "running"

	// XappInstanceStatusSucceeded captures enum value "succeeded"
	XappInstanceStatusSucceeded string = "succeeded"

	// XappInstanceStatusFailed captures enum value "failed"
	XappInstanceStatusFailed string = "failed"

	// XappInstanceStatusUnknown captures enum value "unknown"
	XappInstanceStatusUnknown string = "unknown"

	// XappInstanceStatusCompleted captures enum value "completed"
	XappInstanceStatusCompleted string = "completed"

	// XappInstanceStatusCrashLoopBackOff captures enum value "crashLoopBackOff"
	XappInstanceStatusCrashLoopBackOff string = "crashLoopBackOff"
)

// prop value enum
func (m *XappInstance) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, xappInstanceTypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *XappInstance) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *XappInstance) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *XappInstance) UnmarshalBinary(b []byte) error {
	var res XappInstance
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}