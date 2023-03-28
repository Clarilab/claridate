package claridate

import (
	"regexp"
	"strings"
)

var dashedDateRegex = regexp.MustCompile(`^\d{4}(-\d{1,2}){0,2}$`)
var dottedDateRegex = regexp.MustCompile(`^(\d{1,2}\.){0,2}\d{4}$`)

// DetermineDateFormat receives a date string and returns the format the date string is in.
// It returns an empty string and no error if the input is an empty string.
// Examples:
//
//	   ParseDateFormat("1983-07-20") -> "YYYY-MM-DD"
//		 ParseDateFormat("2006") -> "YYYY"
//
// ParseDateFormat will return an error, if
//  1. the year is represented with more or less than 4 digits
//  2. the month or day are represented with more than 2 digits
//  3. the date string is separated into more than 3 parts
//  4. the date string is separated with something other than hyphens
func DetermineDateFormat(date string) (string, error) {
	if date == "" {
		return "", nil
	}

	if !dashedDateRegex.MatchString(date) {
		return "", ErrUnsupportedDateFormat
	}

	split := strings.Split(date, "-")

	result := "YYYY"

	if len(split) > 1 {
		result = result + "-" + strings.Repeat("M", len(split[1]))
	}

	if len(split) > 2 {
		result = result + "-" + strings.Repeat("D", len(split[2]))
	}

	return result, nil
}

// TransformToDashedDate takes a date string that is in the dotted date format,
// for example "MM.YYYY" or "DD.MM.YYYY", and converts it to the dashed date format.
// If the given date string already is in the dashed date format, it is returned without error.
// It returns an empty string and no error if the input is an empty string.
func TransformToDashedDate(date string) (string, error) {
	if date == "" {
		return "", nil
	}

	if dashedDateRegex.MatchString(date) {
		return date, nil
	}

	if !dottedDateRegex.MatchString(date) {
		return "", ErrUnsupportedDateFormat
	}

	split := strings.Split(date, ".")

	switch len(split) {
	// case 1 (YYYY) can never happen, because that also matches the dashedDateRegex and is returned "as is" in the first if-block of the function
	case 2:
		// MM.YYYY
		return split[1] + "-" + split[0], nil
	case 3:
		// DD.MM.YYYY
		return split[2] + "-" + split[1] + "-" + split[0], nil
	default:
		return "", ErrUnsupportedDateFormat
	}
}
