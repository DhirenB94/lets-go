package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form struct will hold any form data and form errors
type Form struct {
	FormData   url.Values
	FormErrors formErrors
}

// Define a new function that will initialise a new Form struct
func NewForm(formData url.Values) *Form {
	return &Form{
		FormData:   formData,
		FormErrors: map[string][]string{},
	}
}

// Required checks that specific fields in the form data are not blank, if the check fails add this to the formErrrors
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.FormData.Get(field)
		if strings.TrimSpace(value) == "" {
			f.FormErrors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength checks that a specific field in the form data contains a maximum number of characters, if the check fails add this to the formErrrors
func (f *Form) MaxLength(field string, d int) {
	value := f.FormData.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.FormErrors.Add(field, fmt.Sprintf("This field is too long (maximum is %d)", d))
	}
}

// PermittedValues checks that a specific field in the form data is is a permitetd value, if the check fails add this to the formErrrors
func (f *Form) PermittedValues(field string, permValues ...string) {
	value := f.FormData.Get(field)
	for _, pv := range permValues {
		if value == pv {
			return
		}
		f.FormErrors.Add(field, "This field is invalid")
	}

}

// Valid returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.FormErrors) == 0
}
