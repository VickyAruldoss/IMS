package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vickyaruldoss/ims/model"
)

type MockMemberRepository struct {
	mock.Mock
}

func (m *MockMemberRepository) Create(member *model.Member) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockMemberRepository) GetByID(id string) (*model.Member, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *MockMemberRepository) GetAll() ([]*model.Member, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Member), args.Error(1)
}

func (m *MockMemberRepository) Update(member *model.Member) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockMemberRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
