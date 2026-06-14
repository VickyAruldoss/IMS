package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vickyaruldoss/ims/mocks"
	"github.com/vickyaruldoss/ims/model"
	"github.com/vickyaruldoss/ims/repository"
	"github.com/vickyaruldoss/ims/service"
)

func TestCreateMember_Success(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	req := &model.CreateMemberRequest{Name: "Alice", Email: "alice@example.com", Role: "admin"}
	mockRepo.On("Create", mock.AnythingOfType("*model.Member")).Return(nil)

	member, err := svc.CreateMember(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, member.ID)
	assert.Equal(t, "Alice", member.Name)
	assert.Equal(t, "alice@example.com", member.Email)
	assert.Equal(t, "admin", member.Role)
	mockRepo.AssertExpectations(t)
}

func TestCreateMember_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	req := &model.CreateMemberRequest{Name: "Alice", Email: "alice@example.com", Role: "admin"}
	mockRepo.On("Create", mock.AnythingOfType("*model.Member")).Return(assert.AnError)

	member, err := svc.CreateMember(req)

	assert.Error(t, err)
	assert.Nil(t, member)
	mockRepo.AssertExpectations(t)
}

func TestGetMember_Success(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	expected := &model.Member{ID: "1", Name: "Alice", Email: "alice@example.com", Role: "admin"}
	mockRepo.On("GetByID", "1").Return(expected, nil)

	member, err := svc.GetMember("1")

	assert.NoError(t, err)
	assert.Equal(t, expected, member)
	mockRepo.AssertExpectations(t)
}

func TestGetMember_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	mockRepo.On("GetByID", "999").Return(nil, repository.ErrNotFound)

	member, err := svc.GetMember("999")

	assert.Nil(t, member)
	assert.ErrorIs(t, err, repository.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestGetAllMembers_Success(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	expected := []*model.Member{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
	}
	mockRepo.On("GetAll").Return(expected, nil)

	members, err := svc.GetAllMembers()

	assert.NoError(t, err)
	assert.Len(t, members, 2)
	mockRepo.AssertExpectations(t)
}

func TestUpdateMember_Success(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	existing := &model.Member{ID: "1", Name: "Alice", Email: "alice@example.com", Role: "member"}
	mockRepo.On("GetByID", "1").Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Member")).Return(nil)

	req := &model.UpdateMemberRequest{Name: "Alice Updated", Role: "admin"}
	member, err := svc.UpdateMember("1", req)

	assert.NoError(t, err)
	assert.Equal(t, "Alice Updated", member.Name)
	assert.Equal(t, "admin", member.Role)
	assert.Equal(t, "alice@example.com", member.Email)
	mockRepo.AssertExpectations(t)
}

func TestUpdateMember_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	mockRepo.On("GetByID", "999").Return(nil, repository.ErrNotFound)

	req := &model.UpdateMemberRequest{Name: "Ghost"}
	member, err := svc.UpdateMember("999", req)

	assert.Nil(t, member)
	assert.ErrorIs(t, err, repository.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestDeleteMember_Success(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	mockRepo.On("Delete", "1").Return(nil)

	err := svc.DeleteMember("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteMember_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMemberRepository)
	svc := service.NewMemberService(mockRepo)

	mockRepo.On("Delete", "999").Return(repository.ErrNotFound)

	err := svc.DeleteMember("999")

	assert.ErrorIs(t, err, repository.ErrNotFound)
	mockRepo.AssertExpectations(t)
}
