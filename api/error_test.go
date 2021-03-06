package api

import (
	"bytes"
	"testing"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/kylelemons/godebug/pretty"
)

func TestError(t *testing.T) {
	testcases := map[string]struct {
		err             Error
		expectedMessage string
	}{
		"EmptyValues": {
			err:             Error{},
			expectedMessage: "Code: , Message: ",
		},
		"EmptyCode": {
			err: Error{
				Message: "This is a message",
			},
			expectedMessage: "Code: , Message: This is a message",
		},
		"EmptyMessage": {
			err: Error{
				Code: "CODE",
			},
			expectedMessage: "Code: CODE, Message: ",
		},
		"NormalCase": {
			err: Error{
				Code:    UNKNOWN_API_ERROR,
				Message: "Unknown error",
			},
			expectedMessage: "Code: " + UNKNOWN_API_ERROR + ", Message: Unknown error",
		},
	}

	for x, testcase := range testcases {
		if testcase.err.Error() != testcase.expectedMessage {
			t.Errorf("Test %v failed. Received different messages (wanted:%v / received:%v)",
				x, testcase.expectedMessage, testcase.err.Error())
			continue
		}
	}
}

func TestLogErrorMessage(t *testing.T) {
	// Logger
	testOut := bytes.NewBuffer([]byte{})
	logger := logrus.Logger{
		Out:       testOut,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.ErrorLevel,
	}
	testcases := map[string]struct {
		err             *Error
		expectedMessage string
	}{
		"OkCase": {
			err: &Error{
				Code:    "Code Error",
				Message: "Error Message",
			},
			expectedMessage: "{\"code\":\"Code Error\",\"level\":\"error\",\"msg\":\"Error Message\",\"requestID\":\"RequestID\",\"time\"",
		},
	}

	for x, testcase := range testcases {
		LogErrorMessage(&logger, RequestInfo{Admin: true, RequestID: "RequestID", Identifier: "123"}, testcase.err)
		logMessage := testOut.String()
		diff := pretty.Compare(logMessage, testcase.expectedMessage)
		if !strings.Contains(logMessage, testcase.expectedMessage) {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
	}
}
