package claridate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DetermineDateFormat(t *testing.T) {
	testCases := map[string]struct {
		input, expectedOutput string
		expectedError         error
	}{
		"YYYY-MM-DD":                                     {"2006-09-28", "YYYY-MM-DD", nil},
		"YYYY-MM":                                        {"1978-12", "YYYY-MM", nil},
		"YYYY":                                           {"1492", "YYYY", nil},
		"less than 4 digits for the year":                {"22-01-01", "", ErrUnsupportedDateFormat},
		"more than 4 digits for the year":                {"201954-07-20", "", ErrUnsupportedDateFormat},
		"date with dots is unsupported":                  {"20.07.1983", "", ErrUnsupportedDateFormat},
		"date with dots is unsupported 2":                {"07.1983", "", ErrUnsupportedDateFormat},
		"more than 2 digits for the day lead to error":   {"1492-09-324", "", ErrUnsupportedDateFormat},
		"more than 2 digits for the month lead to error": {"1492-908-12", "", ErrUnsupportedDateFormat},
		"some gibberish":                                 {"sfv_24w4e", "", ErrUnsupportedDateFormat},
		"some gibberish 2":                               {"_!@§hahaha", "", ErrUnsupportedDateFormat},
		"some gibberish 3":                               {"            ", "", ErrUnsupportedDateFormat},
		"some gibberish 4":                               {"¯\\_(ツ)_/¯", "", ErrUnsupportedDateFormat},
	}

	for name, tc := range testCases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			parsedFormat, err := DetermineDateFormat(tc.input)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("got unexpected error: %v", err)
				}

				assert.Equal(t, tc.expectedError, err)
			}

			assert.Equal(t, tc.expectedOutput, parsedFormat)
		})
	}
}

func Test_TransformToDashedDate(t *testing.T) {
	testCases := map[string]struct {
		input, expectedOutput string
		expectedError         error
	}{
		"DD.MM.YYYY":                           {"22.04.1712", "1712-04-22", nil},
		"MM.YYYY":                              {"08.1492", "1492-08", nil},
		"YYYY":                                 {"2023", "2023", nil},
		"YYYY.MM.DD":                           {"1983.07.20", "", ErrUnsupportedDateFormat},
		"date string already in dashed format": {"2012-10-03", "2012-10-03", nil},
		"neither dashed nor dotted":            {"20 07 1983", "", ErrUnsupportedDateFormat},
		"some gibberish":                       {"sfv_24w4e", "", ErrUnsupportedDateFormat},
		"some gibberish 2":                     {"_!@§hahaha", "", ErrUnsupportedDateFormat},
		"some gibberish 3":                     {"            ", "", ErrUnsupportedDateFormat},
		"some gibberish 4":                     {"¯\\_(ツ)_/¯", "", ErrUnsupportedDateFormat},
	}

	for name, tc := range testCases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			parsedFormat, err := TransformToDashedDate(tc.input)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("got unexpected error: %v", err)
				}

				assert.Equal(t, tc.expectedError, err)
			}

			assert.Equal(t, tc.expectedOutput, parsedFormat)
		})
	}

}
