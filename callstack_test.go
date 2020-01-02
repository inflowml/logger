// Package logger offers utility functions related to logging.
// This file contains unit tests for the callstack.go module.
package logger

import (
	"os"
	"reflect"
	"testing"
)

// Define several functions to simulate a callstack.
func f1(task func()) { f2(task) }
func f2(task func()) { f3(task) }
func f3(task func()) { f4(task) }
func f4(task func()) { f5(task) }
func f5(task func()) { task() }

// TestGetCallerName tests the GetCallerName() functions.
func TestGetCallerName(t *testing.T) {
	tests := []struct {
		caller   func(func())
		wantName string
	}{
		{f1, "github.com/inflowml/inflow-micro/common/logger.f5"},
		{f2, "github.com/inflowml/inflow-micro/common/logger.f5"},
		{f3, "github.com/inflowml/inflow-micro/common/logger.f5"},
		{f4, "github.com/inflowml/inflow-micro/common/logger.f5"},
		{f5, "github.com/inflowml/inflow-micro/common/logger.f5"},
	}

	var err error
	var name string
	task := func() { name, err = GetCallerName() }

	for i, test := range tests {
		test.caller(task)
		if err != nil {
			t.Errorf("TestGetCallerName()[%d] - failed to get caller name: %v.", i, err)
		} else if name != test.wantName {
			t.Errorf("TestGetCallerName()[%d] = %q, want name %q.", i, name, test.wantName)
		}
	}
}

// TestGetCallStack tests the GetCallStack() functions.
func TestGetCallStack(t *testing.T) {
	tests := []struct {
		caller    func(func())
		size      int
		wantStack []string
	}{
		{
			f5,
			5,
			[]string{
				"callstack_test.go:16",
			},
		}, {
			f4,
			5,
			[]string{
				"callstack_test.go:16",
			},
		}, {
			f4,
			6,
			[]string{
				"callstack_test.go:15",
				"callstack_test.go:16",
			},
		}, {
			f1,
			9,
			[]string{
				"callstack_test.go:12",
				"callstack_test.go:13",
				"callstack_test.go:14",
				"callstack_test.go:15",
				"callstack_test.go:16",
			},
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("TestGetCallStack() - failed to get working directory: %v.", err)
	}

	for i, test := range tests {
		var err error
		var stack []string
		task := func() { stack, err = GetCallStack(test.size) }

		test.caller(task)
		if err != nil {
			t.Errorf("TestGetCallStack()[%d] - failed to get caller name: %v.", i, err)
			continue
		}

		// Prepend the path to the current directory to each stack entry.
		for i, line := range test.wantStack {
			test.wantStack[i] = cwd + "/" + line
		}

		// Ignore the stack prefix to test.caller() and the stack suffix
		// from the anonymous task function.
		haveStack := stack[3 : len(stack)-1]
		if !reflect.DeepEqual(haveStack, test.wantStack) {
			t.Errorf("TestGetCallStack()[%d] = %v, want stack %v.", i, haveStack, test.wantStack)
		}
	}
}
