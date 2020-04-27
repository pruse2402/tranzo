package utils

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
)

// Simple struct to store the Message & Key of a validation error
type ValidationError struct {
	Message, Key string
}

// String returns the Message field of the ValidationError struct.
func (e *ValidationError) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// A Validation context manages data validation and error messages.
type Validation struct {
	Errors []*ValidationError
	keep   bool
}

// Keep tells revel to set a flash cookie on the client to make the validation
// errors available for the next request.
// This is helpful  when redirecting the client after the validation failed.
// It is good practice to always redirect upon a HTTP POST request. Thus
// one should use this method when HTTP POST validation failed and redirect
// the user back to the form.
func (v *Validation) Keep() {
	v.keep = true
}

// Clear *all* ValidationErrors
func (v *Validation) Clear() {
	v.Errors = []*ValidationError{}
}

// HasErrors returns true if there are any (ie > 0) errors. False otherwise.
func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

// ErrorMap returns the errors mapped by key.
// If there are multiple validation errors associated with a single key, the
// first one "wins".  (Typically the first validation will be the more basic).
func (v *Validation) ErrorMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, e := range v.Errors {
		if _, ok := m[e.Key]; !ok {
			m[e.Key] = e.Message
		}
	}
	return m
}

// Error adds an error to the validation context.
func (v *Validation) Error(message string, args ...interface{}) *ValidationResult {
	result := (&ValidationResult{
		Ok:    false,
		Error: &ValidationError{},
	}).Message(message, args...)
	v.Errors = append(v.Errors, result.Error)
	return result
}

// A ValidationResult is returned from every validation method.
// It provides an indication of success, and a pointer to the Error (if any).
type ValidationResult struct {
	Error *ValidationError
	Ok    bool
}

// Key sets the ValidationResult's Error "key" and returns itself for chaining
func (r *ValidationResult) Key(key string) *ValidationResult {
	if r.Error != nil {
		r.Error.Key = key
	}
	return r
}

// Message sets the error message for a ValidationResult. Returns itself to
// allow chaining.  Allows Sprintf() type calling with multiple parameters
func (r *ValidationResult) Message(message string, args ...interface{}) *ValidationResult {
	if r.Error != nil {
		if len(args) == 0 {
			r.Error.Message = message
		} else {
			r.Error.Message = fmt.Sprintf(message, args...)
		}
	}
	return r
}

// Required tests that the argument is non-nil and non-empty (if string or list)
func (v *Validation) Required(obj interface{}) *ValidationResult {
	return v.apply(Required{}, obj)
}

func (v *Validation) Min(n int, min int) *ValidationResult {
	return v.apply(Min{min}, n)
}

func (v *Validation) MinFloat(n float64, min float64) *ValidationResult {
	return v.apply(MinFloat{min}, n)
}

func (v *Validation) Max(n int, max int) *ValidationResult {
	return v.apply(Max{max}, n)
}

func (v *Validation) MaxFloat(n float64, max float64) *ValidationResult {
	return v.apply(MaxFloat{max}, n)
}

func (v *Validation) Range(n, min, max int) *ValidationResult {
	return v.apply(Range{Min{min}, Max{max}}, n)
}

func (v *Validation) RangeFloat(n, min, max float64) *ValidationResult {
	return v.apply(RangeFloat{MinFloat{min}, MaxFloat{max}}, n)
}

func (v *Validation) MinSize(obj interface{}, min int) *ValidationResult {
	return v.apply(MinSize{min}, obj)
}

func (v *Validation) MaxSize(obj interface{}, max int) *ValidationResult {
	return v.apply(MaxSize{max}, obj)
}

func (v *Validation) Length(obj interface{}, n int) *ValidationResult {
	return v.apply(Length{n}, obj)
}

func (v *Validation) Match(str string, regex *regexp.Regexp) *ValidationResult {
	return v.apply(Match{regex}, str)
}

func (v *Validation) Email(str string) *ValidationResult {
	return v.apply(Email{Match{emailPattern}}, str)
}

func (v *Validation) apply(chk Validator, obj interface{}) *ValidationResult {
	if chk.IsSatisfied(obj) {
		return &ValidationResult{Ok: true}
	}

	// Get the default key.
	var key string
	if pc, _, line, ok := runtime.Caller(2); ok {
		f := runtime.FuncForPC(pc)
		if defaultKeys, ok := DefaultValidationKeys[f.Name()]; ok {
			key = defaultKeys[line]
		}
	} else {
		log.Println("Failed to get Caller information to look up Validation key")
	}

	// Add the error to the validation context.
	err := &ValidationError{
		Message: chk.DefaultMessage(),
		Key:     key,
	}
	v.Errors = append(v.Errors, err)

	// Also return it in the result.
	return &ValidationResult{
		Ok:    false,
		Error: err,
	}
}

// Apply a group of validators to a field, in order, and return the
// ValidationResult from the first one that fails, or the last one that
// succeeds.
func (v *Validation) Check(obj interface{}, checks ...Validator) *ValidationResult {
	var result *ValidationResult
	for _, check := range checks {
		result = v.apply(check, obj)
		if !result.Ok {
			return result
		}
	}
	return result
}

// Register default validation keys for all calls to Controller.Validation.Func().
// Map from (package).func => (line => name of first arg to Validation func)
// E.g. "myapp/controllers.helper" or "myapp/controllers.(*Application).Action"
// This is set on initialization in the generated main.go file.
var DefaultValidationKeys map[string]map[int]string
