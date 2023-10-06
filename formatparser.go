package claridate

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dashedDateYearFirstRegex = regexp.MustCompile(`^\d{4}(-\d{1,2}){0,2}$`)
var dashedDateYearLastRegex = regexp.MustCompile(`^(\d{1,2}-){0,2}\d{4}$`)
var dottedDateRegex = regexp.MustCompile(`^(\d{1,2}\.){0,2}(\d{2}|\d{4})$`)
var slashedDateYearLastRegex = regexp.MustCompile(`^(\d{1,2}/){0,2}\d{4}$`)
var slashedDateYearFirstRegex = regexp.MustCompile(`^\d{4}(/\d{1,2}){0,2}$`)
var shortMonthRegex = regexp.MustCompile(`\b(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\b`)

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
	// clean the date string of potential spaces
	date = strings.TrimSpace(date)

	if date == "" {
		return "", nil
	}

	if !dashedDateYearFirstRegex.MatchString(date) {
		return "", ErrUnsupportedDateFormat
	}

	split := strings.Split(date, "-")

	builder := strings.Builder{}
	builder.WriteString("YYYY")

	if len(split) > 1 {
		builder.WriteString("-MM")
	}

	if len(split) > 2 {
		builder.WriteString("-DD")
	}

	return builder.String(), nil
}

// TransformToDashedDate takes a date string that is in the dotted date format,
// for example "MM.YYYY" or "DD.MM.YYYY", and converts it to the dashed date format.
// If the given date string already is in the dashed date format, it is returned without error.
// It returns an empty string and no error if the input is an empty string.
func TransformToDashedDate(date string) (string, error) {
	// clean the date string of potential spaces
	date = strings.TrimSpace(date)

	if date == "" {
		return "", nil
	}

	if strings.Contains(date, ";") {
		date, _, _ = strings.Cut(date, ";")
	}

	if strings.Contains(date, "ca.") {
		_, date, _ = strings.Cut(date, "ca. ")
	}

	if dashedDateYearFirstRegex.MatchString(date) {
		return sanitizeDashedDate(date), nil
	}

	if dashedDateYearLastRegex.MatchString(date) {
		// DD-MM-YYYY: in this case, split the date, reverse the slice and put it back together.
		return strings.Join(reverse(strings.Split(sanitizeDashedDate(date), "-")), "-"), nil
	}

	if shortMonthRegex.MatchString(date) {
		date = strings.ReplaceAll(date, ".", "")
		return parseShortMonthDate(date)
	}

	isDotted := dottedDateRegex.MatchString(date)
	isSlashedYearLast := slashedDateYearLastRegex.MatchString(date)
	isSlashedYearFirst := slashedDateYearFirstRegex.MatchString(date)

	if !isDotted && !isSlashedYearLast && !isSlashedYearFirst {
		return "", ErrUnsupportedDateFormat
	}

	separator := "/"
	if isDotted {
		separator = "."
	}

	split := strings.Split(date, separator)

	var year, month, day string

	switch len(split) {
	case 1: // YYYY
		year = split[0]
	case 2: // MM.YYYY or MM/YYYY or YYYY/MM
		if isSlashedYearFirst { // YYYY/MM
			year = split[0]
			month = split[1]
		} else {
			// MM.YYYY or MM/YYYY
			year = split[1]
			month = split[0]

		}
	case 3: // DD.MM.YYYY or DD/MM/YYYY or YYYY/MM/DD
		if isSlashedYearFirst { // YYYY/MM/DD
			year = split[0]
			month = split[1]
			day = split[2]
		} else {
			// DD.MM.YYYY or DD/MM/YYYY
			year = split[2]
			month = split[1]
			day = split[0]
		}
	default:
		return "", ErrUnsupportedDateFormat
	}

	return buildDashedDateResponse(year, month, day), nil
}

func sanitizeDashedDate(date string) string {
	split := strings.Split(date, "-")
	for i := range split {
		// only happens in cases where day or month are single digits
		if len(split[i]) == 1 {
			split[i] = "0" + split[i]
		}
	}

	return strings.Join(split, "-")
}

func transformTwoDigitYearToFourDigitYear(year string) string {
	// we can safely ignore this error, as we already matched with the regex that this value is an integer
	yearInt, _ := strconv.Atoi(year)

	// this is exactly how go converts 2 digit years to 4 digit years
	// https://cs.opensource.google/go/go/+/master:src/time/format.go;l=1082;drc=3ad6393f8676b1b408673bf40b8a876f29561eef
	if yearInt >= 69 {
		yearInt += 1900
	} else {
		yearInt += 2000
	}

	return strconv.Itoa(yearInt)
}

func parseShortMonthDate(date string) (string, error) {
	wordsAmount := len(strings.Fields(date))
	switch wordsAmount {
	case 2: // example: Jul 1957
		parsedDate, err := time.Parse("Jan 2006", date)
		if err != nil {
			return "", ErrUnsupportedDateFormat
		}
		return parsedDate.Format("2006-01"), nil
	case 3: // example: 30 Jul 1957
		dayLength := len(strings.Fields(date)[0])
		if dayLength == 1 {
			parsedDate, err := time.Parse("2 Jan 2006", date)
			if err != nil {
				return "", ErrUnsupportedDateFormat
			}
			return parsedDate.Format("2006-01-02"), nil
		}

		parsedDate, err := time.Parse("02 Jan 2006", date)
		if err != nil {
			return "", ErrUnsupportedDateFormat
		}
		return parsedDate.Format("2006-01-02"), nil
	default:
		return "", ErrUnsupportedDateFormat
	}
}

func reverse(strSlice []string) []string {
	for i, j := 0, len(strSlice)-1; i < j; i, j = i+1, j-1 {
		strSlice[i], strSlice[j] = strSlice[j], strSlice[i]
	}

	return strSlice
}

func buildDashedDateResponse(year, month, day string) string {
	builder := strings.Builder{}
	if year != "" {
		if len(year) == 2 {
			year = transformTwoDigitYearToFourDigitYear(year)
		}

		builder.WriteString(year)

		if month != "" {
			builder.WriteString("-")
		}
	}

	if month != "" {
		if len(month) == 1 {
			month = "0" + month
		}

		builder.WriteString(month)
		if day != "" {
			builder.WriteString("-")
		}
	}

	if day != "" {
		if len(day) == 1 {
			day = "0" + day
		}

		builder.WriteString(day)
	}

	return builder.String()
}
