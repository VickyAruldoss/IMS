package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vickyaruldoss/ims/controller"
	"github.com/vickyaruldoss/ims/mocks"
	"github.com/vickyaruldoss/ims/model"
	"github.com/vickyaruldoss/ims/repository"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter(ctrl *controller.MemberController) *gin.Engine {
	r := gin.New()
	r.POST("/members", ctrl.CreateMember)
	r.GET("/members", ctrl.GetAllMembers)
	r.GET("/members/:id", ctrl.GetMember)
	r.PUT("/members/:id", ctrl.UpdateMember)
	r.DELETE("/members/:id", ctrl.DeleteMember)
	return r
}

func TestCreateMember_Created(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	expected := &model.Member{ID: "abc", Name: "Alice", Email: "alice@example.com", Role: "admin"}
	mockSvc.On("CreateMember", mock.AnythingOfType("*model.CreateMemberRequest")).Return(expected, nil)

	body, _ := json.Marshal(model.CreateMemberRequest{Name: "Alice", Email: "alice@example.com", Role: "admin"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var got model.Member
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Alice", got.Name)
	assert.Equal(t, "abc", got.ID)
	mockSvc.AssertExpectations(t)
}

func TestCreateMember_BadRequest(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	// Missing required fields
	body, _ := json.Marshal(map[string]string{"name": "Alice"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "CreateMember")
}

func TestGetMember_OK(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	expected := &model.Member{ID: "1", Name: "Alice", Email: "alice@example.com", Role: "admin"}
	mockSvc.On("GetMember", "1").Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/members/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got model.Member
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Alice", got.Name)
	mockSvc.AssertExpectations(t)
}

func TestGetMember_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	mockSvc.On("GetMember", "999").Return(nil, repository.ErrNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/members/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetAllMembers_OK(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	expected := []*model.Member{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
	}
	mockSvc.On("GetAllMembers").Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/members", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got []*model.Member
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Len(t, got, 2)
	mockSvc.AssertExpectations(t)
}

func TestUpdateMember_OK(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	updated := &model.Member{ID: "1", Name: "Alice Updated", Email: "alice@example.com", Role: "admin"}
	mockSvc.On("UpdateMember", "1", mock.AnythingOfType("*model.UpdateMemberRequest")).Return(updated, nil)

	body, _ := json.Marshal(model.UpdateMemberRequest{Name: "Alice Updated", Role: "admin"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/members/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got model.Member
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Alice Updated", got.Name)
	mockSvc.AssertExpectations(t)
}

func TestUpdateMember_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	mockSvc.On("UpdateMember", "999", mock.AnythingOfType("*model.UpdateMemberRequest")).Return(nil, repository.ErrNotFound)

	body, _ := json.Marshal(model.UpdateMemberRequest{Name: "Ghost"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/members/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteMember_NoContent(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	mockSvc.On("DeleteMember", "1").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/members/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteMember_NotFound(t *testing.T) {
	mockSvc := new(mocks.MockMemberService)
	ctrl := controller.NewMemberController(mockSvc)
	r := newTestRouter(ctrl)

	mockSvc.On("DeleteMember", "999").Return(repository.ErrNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/members/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}
