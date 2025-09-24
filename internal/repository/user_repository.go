package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserEntity struct {
	bun.BaseModel `bun:"users"`

	Id        uuid.UUID `bun:"id,pk,type:uuid"`
	Name      string    `bun:"name,notnull"`
	Email     string    `bun:"email,notnull"`
	Password  string    `bun:"password,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

type UserRepository interface {
	GetUsers(offset, limit int) ([]UserEntity, error)
	GetUsersCount() (int64, error)
	InsertUser(user UserEntity) (out UserEntity, err error)
	GetUserById(id uuid.UUID) (out UserEntity, err error)
	GetUserByEmail(email string) (out UserEntity, err error)
}

type DefaultUserRepository struct {
	db  *bun.DB
	ctx context.Context
}

func NewUserRepository(db *bun.DB, ctx context.Context) UserRepository {
	return &DefaultUserRepository{db: db, ctx: ctx}
}

func (r *DefaultUserRepository) GetUsers(offset, limit int) ([]UserEntity, error) {
	var users []UserEntity
	err := r.db.NewSelect().
		Model(&users).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Scan(r.ctx)
	if err != nil {
		return []UserEntity{}, err
	}
	return users, nil
}

func (r *DefaultUserRepository) GetUsersCount() (int64, error) {
	count, err := r.db.NewSelect().
		Model((*UserEntity)(nil)).
		Count(r.ctx)
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

func (r *DefaultUserRepository) InsertUser(user UserEntity) (out UserEntity, err error) {
	_, err = r.db.NewInsert().Model(&user).Exec(r.ctx)
	if err != nil {
		return out, err
	}
	return user, nil
}

func (r *DefaultUserRepository) GetUserById(id uuid.UUID) (out UserEntity, err error) {
	err = r.db.NewSelect().Model(&out).Where("id = ?", id).Scan(r.ctx)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (r *DefaultUserRepository) GetUserByEmail(email string) (out UserEntity, err error) {
	err = r.db.NewSelect().Model(&out).Where("email = ?", email).Scan(r.ctx)
	if err != nil {
		return out, err
	}
	return out, nil
}
