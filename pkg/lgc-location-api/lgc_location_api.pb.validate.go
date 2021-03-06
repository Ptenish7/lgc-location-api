// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: ozonmp/lgc_location_api/v1/lgc_location_api.proto

package lgc_location_api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// Validate checks the field values on Location with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Location) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for Latitude

	// no validation rules for Longitude

	// no validation rules for Title

	return nil
}

// LocationValidationError is the validation error returned by
// Location.Validate if the designated constraints aren't met.
type LocationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LocationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LocationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LocationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LocationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LocationValidationError) ErrorName() string { return "LocationValidationError" }

// Error satisfies the builtin error interface
func (e LocationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLocation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LocationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LocationValidationError{}

// Validate checks the field values on CreateLocationV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateLocationV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if val := m.GetLatitude(); val < -90 || val > 90 {
		return CreateLocationV1RequestValidationError{
			field:  "Latitude",
			reason: "value must be inside range [-90, 90]",
		}
	}

	if val := m.GetLongitude(); val < -180 || val > 180 {
		return CreateLocationV1RequestValidationError{
			field:  "Longitude",
			reason: "value must be inside range [-180, 180]",
		}
	}

	if utf8.RuneCountInString(m.GetTitle()) < 1 {
		return CreateLocationV1RequestValidationError{
			field:  "Title",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// CreateLocationV1RequestValidationError is the validation error returned by
// CreateLocationV1Request.Validate if the designated constraints aren't met.
type CreateLocationV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateLocationV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateLocationV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateLocationV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateLocationV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateLocationV1RequestValidationError) ErrorName() string {
	return "CreateLocationV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateLocationV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateLocationV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateLocationV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateLocationV1RequestValidationError{}

// Validate checks the field values on CreateLocationV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateLocationV1Response) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for LocationId

	return nil
}

// CreateLocationV1ResponseValidationError is the validation error returned by
// CreateLocationV1Response.Validate if the designated constraints aren't met.
type CreateLocationV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateLocationV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateLocationV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateLocationV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateLocationV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateLocationV1ResponseValidationError) ErrorName() string {
	return "CreateLocationV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateLocationV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateLocationV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateLocationV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateLocationV1ResponseValidationError{}

// Validate checks the field values on DescribeLocationV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribeLocationV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetLocationId() <= 0 {
		return DescribeLocationV1RequestValidationError{
			field:  "LocationId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// DescribeLocationV1RequestValidationError is the validation error returned by
// DescribeLocationV1Request.Validate if the designated constraints aren't met.
type DescribeLocationV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribeLocationV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribeLocationV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribeLocationV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribeLocationV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribeLocationV1RequestValidationError) ErrorName() string {
	return "DescribeLocationV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e DescribeLocationV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribeLocationV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribeLocationV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribeLocationV1RequestValidationError{}

// Validate checks the field values on DescribeLocationV1Response with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribeLocationV1Response) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetLocation()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DescribeLocationV1ResponseValidationError{
				field:  "Location",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// DescribeLocationV1ResponseValidationError is the validation error returned
// by DescribeLocationV1Response.Validate if the designated constraints aren't met.
type DescribeLocationV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribeLocationV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribeLocationV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribeLocationV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribeLocationV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribeLocationV1ResponseValidationError) ErrorName() string {
	return "DescribeLocationV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DescribeLocationV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribeLocationV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribeLocationV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribeLocationV1ResponseValidationError{}

// Validate checks the field values on ListLocationsV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListLocationsV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetLimit() <= 0 {
		return ListLocationsV1RequestValidationError{
			field:  "Limit",
			reason: "value must be greater than 0",
		}
	}

	if m.GetOffset() < 0 {
		return ListLocationsV1RequestValidationError{
			field:  "Offset",
			reason: "value must be greater than or equal to 0",
		}
	}

	return nil
}

// ListLocationsV1RequestValidationError is the validation error returned by
// ListLocationsV1Request.Validate if the designated constraints aren't met.
type ListLocationsV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListLocationsV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListLocationsV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListLocationsV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListLocationsV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListLocationsV1RequestValidationError) ErrorName() string {
	return "ListLocationsV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListLocationsV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListLocationsV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListLocationsV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListLocationsV1RequestValidationError{}

// Validate checks the field values on ListLocationsV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListLocationsV1Response) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetLocations() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListLocationsV1ResponseValidationError{
					field:  fmt.Sprintf("Locations[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ListLocationsV1ResponseValidationError is the validation error returned by
// ListLocationsV1Response.Validate if the designated constraints aren't met.
type ListLocationsV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListLocationsV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListLocationsV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListLocationsV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListLocationsV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListLocationsV1ResponseValidationError) ErrorName() string {
	return "ListLocationsV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListLocationsV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListLocationsV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListLocationsV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListLocationsV1ResponseValidationError{}

// Validate checks the field values on UpdateLocationV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UpdateLocationV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetLocationId() <= 0 {
		return UpdateLocationV1RequestValidationError{
			field:  "LocationId",
			reason: "value must be greater than 0",
		}
	}

	if val := m.GetLatitude(); val < -90 || val > 90 {
		return UpdateLocationV1RequestValidationError{
			field:  "Latitude",
			reason: "value must be inside range [-90, 90]",
		}
	}

	if val := m.GetLongitude(); val < -180 || val > 180 {
		return UpdateLocationV1RequestValidationError{
			field:  "Longitude",
			reason: "value must be inside range [-180, 180]",
		}
	}

	if utf8.RuneCountInString(m.GetTitle()) < 1 {
		return UpdateLocationV1RequestValidationError{
			field:  "Title",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// UpdateLocationV1RequestValidationError is the validation error returned by
// UpdateLocationV1Request.Validate if the designated constraints aren't met.
type UpdateLocationV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateLocationV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateLocationV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateLocationV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateLocationV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateLocationV1RequestValidationError) ErrorName() string {
	return "UpdateLocationV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateLocationV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateLocationV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateLocationV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateLocationV1RequestValidationError{}

// Validate checks the field values on UpdateLocationV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UpdateLocationV1Response) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// UpdateLocationV1ResponseValidationError is the validation error returned by
// UpdateLocationV1Response.Validate if the designated constraints aren't met.
type UpdateLocationV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateLocationV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateLocationV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateLocationV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateLocationV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateLocationV1ResponseValidationError) ErrorName() string {
	return "UpdateLocationV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateLocationV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateLocationV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateLocationV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateLocationV1ResponseValidationError{}

// Validate checks the field values on RemoveLocationV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemoveLocationV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetLocationId() <= 0 {
		return RemoveLocationV1RequestValidationError{
			field:  "LocationId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// RemoveLocationV1RequestValidationError is the validation error returned by
// RemoveLocationV1Request.Validate if the designated constraints aren't met.
type RemoveLocationV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveLocationV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveLocationV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveLocationV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveLocationV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveLocationV1RequestValidationError) ErrorName() string {
	return "RemoveLocationV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e RemoveLocationV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveLocationV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveLocationV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveLocationV1RequestValidationError{}

// Validate checks the field values on RemoveLocationV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemoveLocationV1Response) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// RemoveLocationV1ResponseValidationError is the validation error returned by
// RemoveLocationV1Response.Validate if the designated constraints aren't met.
type RemoveLocationV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveLocationV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveLocationV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveLocationV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveLocationV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveLocationV1ResponseValidationError) ErrorName() string {
	return "RemoveLocationV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RemoveLocationV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveLocationV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveLocationV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveLocationV1ResponseValidationError{}

// Validate checks the field values on LocationEvent with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *LocationEvent) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for LocationId

	// no validation rules for Type

	// no validation rules for ExtraType

	// no validation rules for Status

	if v, ok := interface{}(m.GetEntity()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LocationEventValidationError{
				field:  "Entity",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return LocationEventValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// LocationEventValidationError is the validation error returned by
// LocationEvent.Validate if the designated constraints aren't met.
type LocationEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LocationEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LocationEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LocationEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LocationEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LocationEventValidationError) ErrorName() string { return "LocationEventValidationError" }

// Error satisfies the builtin error interface
func (e LocationEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLocationEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LocationEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LocationEventValidationError{}
