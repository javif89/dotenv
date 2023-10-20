package dotenv

import "testing"

type standardizeKeyTest struct {
	arg, expected string
}

var standardizeKeyTests = []standardizeKeyTest{
	{"test", "TEST"},
	{"test test", "TEST_TEST"},
	{"test-test", "TEST_TEST"},
	{"test+test", "TEST_TEST"},
}

func TestStandardizeKey(t *testing.T) {
	for _, test := range standardizeKeyTests {
		actual := standardizeKey(test.arg)
		if actual != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, actual)
		}
	}
}

type formatForPrintTest struct {
    arg1, arg2, expected string
}

var formatForPrintTests = []formatForPrintTest{
    {"test", "", "test"},
	{"test", "comment", "test # comment"},
	{"test test", "", "\"test test\""},
	{"test test", "comment", "\"test test\" # comment"},
	{"test-test", "", "\"test-test\""},
	{"test-test", "comment", "\"test-test\" # comment"},
	{"test+test", "", "\"test+test\""},
	{"test+test", "comment", "\"test+test\" # comment"},
}

func TestFormatValueForPrint(t *testing.T) {
	for _, test := range formatForPrintTests {
		actual := formatValueForPrint(test.arg1, test.arg2)
		if actual != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, actual)
		}
	}
}