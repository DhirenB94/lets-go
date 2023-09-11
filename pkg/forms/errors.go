package forms

type formErrors map[string][]string

// Add will append error messages for a given field
func (fve formErrors) Add(field, message string) {
	fve[field] = append(fve[field], message)
}

// Get will retrieve the 1st error message for a given field
func (fve formErrors) Get(field string) string {
	errorMessages := fve[field]
	if len(errorMessages) == 0 {
		return ""
	}
	return errorMessages[0]
}
