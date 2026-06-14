package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/vickyaruldoss/ims/model"
	"github.com/vickyaruldoss/ims/repository"
)

type MemberService interface {
	CreateMember(req *model.CreateMemberRequest) (*model.Member, error)
	GetMember(id string) (*model.Member, error)
	GetAllMembers() ([]*model.Member, error)
	UpdateMember(id string, req *model.UpdateMemberRequest) (*model.Member, error)
	DeleteMember(id string) error
}

type memberService struct {
	repo repository.MemberRepository
}

func NewMemberService(repo repository.MemberRepository) MemberService {
	return &memberService{repo: repo}
}

func (s *memberService) CreateMember(req *model.CreateMemberRequest) (*model.Member, error) {
	member := &model.Member{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *memberService) GetMember(id string) (*model.Member, error) {
	return s.repo.GetByID(id)
}

func (s *memberService) GetAllMembers() ([]*model.Member, error) {
	return s.repo.GetAll()
}

func (s *memberService) UpdateMember(id string, req *model.UpdateMemberRequest) (*model.Member, error) {
	member, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		member.Name = req.Name
	}
	if req.Email != "" {
		member.Email = req.Email
	}
	if req.Role != "" {
		member.Role = req.Role
	}
	member.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *memberService) DeleteMember(id string) error {
	return s.repo.Delete(id)
}
