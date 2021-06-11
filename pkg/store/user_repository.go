package store

import (
	"GoStore/pkg/models"
	"context"
	"database/sql"
	"errors"
)

type UserRepository struct {
	store *Store
}

func (userRepository *UserRepository) UserByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	if err := userRepository.store.DB.QueryRowxContext(ctx, `SELECT uuid, 
       name, 
       email, 
       password, 
       create_at, 
       updated_at
FROM "public"."admins" us 
WHERE us.email = $1`, user.Email).StructScan(user); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("bad login")
		}
		return nil, err
	}

	return user, nil
}

func (userRepository *UserRepository) UserByUUID(ctx context.Context, user *models.User) (*models.User, error) {
	if err := userRepository.store.DB.QueryRowxContext(ctx, `SELECT uuid, 
       name, 
       email, 
       password, 
       create_at, 
       updated_at
FROM "public"."admins" us 
WHERE us.uuid = $1`, user.UUID).StructScan(user); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("bad login")
		}
		return nil, err
	}

	return user, nil
}

func (userRepository *UserRepository) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := userRepository.store.DB.QueryRowxContext(ctx, `insert into "public"."admins"(name, email, password) 
VALUES($1, $2, $3) returning uuid`, user.Name, user.Email, user.Password).Scan(&user.UUID); err != nil {
		return nil, err
	}

	return user, nil
}

