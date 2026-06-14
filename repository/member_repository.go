package repository

import (
	"errors"
	"sync"

	"github.com/vickyaruldoss/ims/model"
)

var ErrNotFound = errors.New("member not found")

type MemberRepository interface {
	Create(member *model.Member) error
	GetByID(id string) (*model.Member, error)
	GetAll() ([]*model.Member, error)
	Update(member *model.Member) error
	Delete(id string) error
}

type inMemoryMemberRepository struct {
	mu      sync.RWMutex
	members map[string]*model.Member
}

func NewInMemoryMemberRepository() MemberRepository {
	return &inMemoryMemberRepository{
		members: make(map[string]*model.Member),
	}
}

func (r *inMemoryMemberRepository) Create(member *model.Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.members[member.ID] = member
	return nil
}

func (r *inMemoryMemberRepository) GetByID(id string) (*model.Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	member, ok := r.members[id]
	if !ok {
		return nil, ErrNotFound
	}
	return member, nil
}

func (r *inMemoryMemberRepository) GetAll() ([]*model.Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	members := make([]*model.Member, 0, len(r.members))
	for _, m := range r.members {
		members = append(members, m)
	}
	return members, nil
}

func (r *inMemoryMemberRepository) Update(member *model.Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.members[member.ID]; !ok {
		return ErrNotFound
	}
	r.members[member.ID] = member
	return nil
}

func (r *inMemoryMemberRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.members[id]; !ok {
		return ErrNotFound
	}
	delete(r.members, id)
	return nil
}
