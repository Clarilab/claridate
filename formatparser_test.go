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
		"YYYY-MM-DD":                      {"2006-09-28", "YYYY-MM-DD", nil},
		"YYYY-MM":                         {"1978-12", "YYYY-MM", nil},
		"YYYY":                            {"1492", "YYYY", nil},
		"YYYY-M-DD":                       {"1492-8-22", "YYYY-MM-DD", nil},
		"YYYY-MM-D":                       {"1492-11-5", "YYYY-MM-DD", nil},
		"less than 4 digits for the year": {"22-01-01", "", ErrUnsupportedDateFormat},
		"more than 4 digits for the year": {"201954-07-20", "", ErrUnsupportedDateFormat},
		"date with dots is unsupported":   {"20.07.1983", "", ErrUnsupportedDateFormat},
		"date with dots is unsupported 2": {"07.1983", "", ErrUnsupportedDateFormat},
		"more than 2 digits for the day lead to error":   {"1492-09-324", "", ErrUnsupportedDateFormat},
		"more than 2 digits for the month lead to error": {"1492-908-12", "", ErrUnsupportedDateFormat},
		"some gibberish":   {"sfv_24w4e", "", ErrUnsupportedDateFormat},
		"some gibberish 2": {"_!@§hahaha", "", ErrUnsupportedDateFormat},
		"some gibberish 3": {"            ", "", ErrUnsupportedDateFormat},
		"some gibberish 4": {"¯\\_(ツ)_/¯", "", ErrUnsupportedDateFormat},
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
		"DD.MM.YYYY":                            {"22.04.1712", "1712-04-22", nil},
		"D.M.YYYY":                              {"2.4.1712", "1712-4-2", nil},
		"D.MM.YYYY":                             {"2.04.1712", "1712-04-2", nil},
		"D/MM/YYYY":                             {"2/04/1712", "1712-04-2", nil},
		"D/M/YYYY":                              {"2/4/1712", "1712-4-2", nil},
		"YYYY-M-DD":                             {"1712-4-22", "1712-4-22", nil},
		"DD.M.YYYY":                             {"22.4.1712", "1712-4-22", nil},
		"MM.YYYY":                               {"08.1492", "1492-08", nil},
		"YYYY":                                  {"2023", "2023", nil},
		"YYYY.MM.DD":                            {"1983.07.20", "", ErrUnsupportedDateFormat},
		"date string already in dashed format":  {"2012-10-03", "2012-10-03", nil},
		"neither dashed nor dotted":             {"20 07 1983", "", ErrUnsupportedDateFormat},
		"some gibberish":                        {"sfv_24w4e", "", ErrUnsupportedDateFormat},
		"some gibberish 2":                      {"_!@§hahaha", "", ErrUnsupportedDateFormat},
		"some gibberish 3":                      {"            ", "", ErrUnsupportedDateFormat},
		"some gibberish 4":                      {"¯\\_(ツ)_/¯", "", ErrUnsupportedDateFormat},
		"can handle forward slashes":            {"20/07/1983", "1983-07-20", nil},
		"can handle forward slashes 2":          {"1983/07/20", "1983-07-20", nil},
		"can handle dashed date with year last": {"20-07-1983", "1983-07-20", nil},
		"can handle dashed date with year last 2": {"20-7-1983", "1983-7-20", nil},
		"leading/trailing space":                  {" 30.07.1957 ", "1957-07-30", nil},
		"ca. DD.MM.YYYY":                          {"ca. 30.07.1984", "1984-07-30", nil},
		"ca. MM.YYYY":                             {"ca. 07.1984", "1984-07", nil},
		"ca. YYYY":                                {"ca. 1984", "1984", nil},
		"DD Mon YYYY with ;":                      {"30 Jul 1957; 1958", "1957-07-30", nil},
		"DD Mon. YYYY":                            {"30 Jul. 1957", "1957-07-30", nil},
		"D Mon YYYY":                              {"3 Jul. 1957", "1957-07-3", nil},
		"Mon YYYY":                                {"Jul 1867", "1867-07", nil},
		"Mon. YYYY":                               {"Jul. 1867", "1867-07", nil},
		"DD.MM.YYYY with ;":                       {"30.07.1957; 1958", "1957-07-30", nil},
		"ca. gibberish":                           {"ca. foobar", "", ErrUnsupportedDateFormat},
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
