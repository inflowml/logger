// Package logger offers utility functions related to logging.
// This file contains unit tests for the logger.go module.
package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
	"time"
)

// TestSeverity tests the Network(), Info(), Warning(), Error(), and Fatal() functions.
func TestSeverity(t *testing.T) {
	os.Setenv("DEBUG", "1")

	tests := []struct {
		time    time.Time
		msg     string
		args    []interface{}
		wantFmt string
	}{
		{
			time.Date(1234, 5, 6, 7, 8, 9, 0, time.UTC),
			"Message",
			[]interface{}{},
			"[1234/05/06 07:08:09.000].*%s.*: Message\n",
		}, {
			time.Date(2016, 11, 8, 22, 0, 0, 0, time.UTC),
			"Trump received %d electoral votes compared to %d for Clinton.",
			[]interface{}{306, 232},
			"[2016/11/08 22:00:00.000].*%s.*: Trump received 306 electoral votes compared to 232 for Clinton.\n",
		}, {
			time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
			"The %s hackers leaked %d emails.",
			[]interface{}{"Russian", 150000},
			"[0000/01/01 00:00:00.000].*%s.*: The Russian hackers leaked 150000 emails.\n",
		},
	}

	halt = func(_ interface{}) {}

	sevMap := []struct {
		sev string
		f   func(msg string, args ...interface{})
	}{
		{"NET", Network},
		{"INFO", Info},
		{"WARN", Warning},
		{"ERROR", Error},
		{"ERROR", Fatal},
	}

	for _, sevFunc := range sevMap {
		sev := sevFunc.sev
		f := sevFunc.f

		for i, test := range tests {
			exe := func() {
				getTime = func() time.Time { return test.time }
				f(test.msg, test.args...)
			}
			out := captureStdout(t, exe)
			wantOut := fmt.Sprintf(test.wantFmt, sev)
			if match, _ := regexp.MatchString(wantOut, out); !match {
				t.Errorf("TestSeverity()[%s][%d] = %q, want output %q.", sev, i, out, wantOut)
			}
		}
	}
}

// TestPrint tests the print() function.
func TestPrint(t *testing.T) {
	tests := []struct {
		time    time.Time
		sev     string
		msg     string
		args    []interface{}
		wantOut string
	}{
		{
			time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC),
			"Severity",
			"Message",
			[]interface{}{},
			"[2000/01/02 03:04:05].*Severity.*: Message\n",
		}, {
			time.Date(2018, 1, 13, 8, 7, 0, 0, time.UTC),
			"WARNING",
			"Missile incoming in %d minutes.",
			[]interface{}{38},
			"[2018/01/13 08:07:00].*WARNING.*: Missile incoming in 38 minutes.\n",
		}, {
			time.Date(0, 12, 25, 12, 34, 56, 789, time.UTC),
			"INFO",
			"The first %d digits of %s are %.2f.",
			[]interface{}{3, "Pi", 3.14159},
			"[0000/12/25 12:34:56].*INFO.*: The first 3 digits of Pi are 3.14.\n",
		},
	}

	for i, test := range tests {
		exe := func() {
			getTime = func() time.Time { return test.time }
			print(test.sev, test.msg, "", test.args...)
		}
		out := captureStdout(t, exe)
		if match, _ := regexp.MatchString(test.wantOut, out); !match {
			t.Errorf("TestPrint()[%d] = %q, want output %q.", i, out, test.wantOut)
		}
	}
}

// captureStdout returns a string representation of stdout when the
// given function is called using the specified arguments.
func captureStdout(t *testing.T, f func()) string {
	defer func(out *os.File) { os.Stdout = out }(os.Stdout)

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Failed to create OS pipe.")
	}

	os.Stdout = w
	f()
	w.Close()

	haveFile, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal("Failed to read OS pipe.")
	}
	return string(haveFile)
}
