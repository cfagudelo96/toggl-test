// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/cfagudelo96/toggle-test/deck/domain"
	mock "github.com/stretchr/testify/mock"
)

// DeckRepository is an autogenerated mock type for the DeckRepository type
type DeckRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, uuid
func (_m *DeckRepository) Get(ctx context.Context, uuid string) (*domain.Deck, error) {
	ret := _m.Called(ctx, uuid)

	var r0 *domain.Deck
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Deck); ok {
		r0 = rf(ctx, uuid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Deck)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, uuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, d
func (_m *DeckRepository) Save(ctx context.Context, d *domain.Deck) error {
	ret := _m.Called(ctx, d)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Deck) error); ok {
		r0 = rf(ctx, d)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
