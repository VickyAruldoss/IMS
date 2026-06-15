package repository

import (
	"database/sql"
	"errors"
	"sync"

	_ "github.com/lib/pq"
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

// --- PostgreSQL implementation ---

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) MemberRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(member *model.Member) error {
	const q = `INSERT INTO members (id, name, email, role, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(q, member.ID, member.Name, member.Email, member.Role, member.CreatedAt, member.UpdatedAt)
	return err
}

func (r *postgresRepository) GetByID(id string) (*model.Member, error) {
	const q = `SELECT id, name, email, role, created_at, updated_at FROM members WHERE id = $1`
	m := &model.Member{}
	err := r.db.QueryRow(q, id).Scan(&m.ID, &m.Name, &m.Email, &m.Role, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *postgresRepository) GetAll() ([]*model.Member, error) {
	const q = `SELECT id, name, email, role, created_at, updated_at FROM members ORDER BY created_at DESC`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*model.Member
	for rows.Next() {
		m := &model.Member{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Role, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *postgresRepository) Update(member *model.Member) error {
	const q = `UPDATE members SET name=$1, email=$2, role=$3, updated_at=$4 WHERE id=$5`
	res, err := r.db.Exec(q, member.Name, member.Email, member.Role, member.UpdatedAt, member.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepository) Delete(id string) error {
	const q = `DELETE FROM members WHERE id=$1`
	res, err := r.db.Exec(q, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// --- In-memory implementation (used by unit tests) ---

type inMemoryMemberRepository struct {
	mu      sync.RWMutex
	members map[string]*model.Member
}

func NewInMemoryMemberRepository() MemberRepository {
	return &inMemoryMemberRepository{members: make(map[string]*model.Member)}
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
	m, ok := r.members[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
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
