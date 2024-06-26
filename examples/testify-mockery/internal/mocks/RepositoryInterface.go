// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/ciazhar/go-zhar/examples/testify-mockery/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// RepositoryInterface is an autogenerated mock type for the RepositoryInterface type
type RepositoryInterface struct {
	mock.Mock
}

// GetAccidentReport provides a mock function with given fields: ctx, id
func (_m *RepositoryInterface) GetAccidentReport(ctx context.Context, id string) (*model.AccidentReport, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetAccidentReport")
	}

	var r0 *model.AccidentReport
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.AccidentReport, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.AccidentReport); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccidentReport)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepositoryInterface creates a new instance of RepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *RepositoryInterface {
	mock := &RepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
