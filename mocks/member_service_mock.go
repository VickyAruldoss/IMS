package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/vickyaruldoss/ims/model"
)

type MockMemberService struct {
	mock.Mock
}

func (m *MockMemberService) CreateMember(req *model.CreateMemberRequest) (*model.Member, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *MockMemberService) GetMember(id string) (*model.Member, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *MockMemberService) GetAllMembers() ([]*model.Member, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Member), args.Error(1)
}

func (m *MockMemberService) UpdateMember(id string, req *model.UpdateMemberRequest) (*model.Member, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *MockMemberService) DeleteMember(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
