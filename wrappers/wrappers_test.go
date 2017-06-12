package wrappers

import (
	"github.com/mushorg/go-dpi"
	"github.com/pkg/errors"
	"testing"
)

type MockClassifier struct {
	initializeSuccessfully bool
	initializeCalled       bool
	destroyCalled          bool
	classifyCalled         bool
}

func (classifier *MockClassifier) InitializeWrapper() error {
	classifier.initializeCalled = true
	if classifier.initializeSuccessfully {
		return nil
	} else {
		return errors.New("Init fail")
	}
}

func (classifier *MockClassifier) DestroyWrapper() error {
	classifier.destroyCalled = true
	return nil
}

func (classifier *MockClassifier) ClassifyFlow(flow *godpi.Flow) (godpi.Protocol, error) {
	classifier.classifyCalled = true
	return godpi.Http, nil
}

func TestClassifyFlowUninitialized(t *testing.T) {
	flow := godpi.NewFlow()
	uninitialized := &MockClassifier{initializeSuccessfully: false}
	wrappersList = [...]Wrapper{
		uninitialized,
	}
	InitializeWrappers()
	if !uninitialized.initializeCalled {
		t.Error("Initialize not called on wrapper")
	}
	result := ClassifyFlow(flow)
	if uninitialized.classifyCalled {
		t.Error("Classify called on uninitialized wrapper")
	}
	if result != godpi.Unknown {
		t.Error("Empty classify did not return unknown")
	}
	DestroyWrappers()
	if uninitialized.destroyCalled {
		t.Error("Destroy called on uninitialized wrapper")
	}
}

func TestClassifyFlowInitialized(t *testing.T) {
	flow := godpi.NewFlow()
	initialized := &MockClassifier{initializeSuccessfully: true}
	wrappersList = [...]Wrapper{
		initialized,
	}
	InitializeWrappers()
	if !initialized.initializeCalled {
		t.Error("Initialize not called on wrapper")
	}
	result := ClassifyFlow(flow)
	if !initialized.classifyCalled {
		t.Error("Classify not called on active wrapper")
	}
	if result != godpi.Http {
		t.Error("Classify did not return correct result")
	}
	DestroyWrappers()
	if !initialized.destroyCalled {
		t.Error("Destroy not called on active wrapper")
	}
}
