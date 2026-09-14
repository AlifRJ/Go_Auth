package service

import (
	"context"
	"errors"

	"github.com/AlifRJ/Go_Auth/app/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password must be 6 characters or longer")
	ErrEmailTaken       = errors.New("email is already registered")
	ErrUsernameTaken    = errors.New("username is already taken")
)

type UserService struct {
	repo model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll(ctx context.Context, limit, offset int) ([]*model.User, error) {
	if limit <= 0 {
		limit = 10 
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) RegisterUser(ctx context.Context, name, username, email, password string) (*model.User, error) {
	// validate password
	if len(password) < 6 {
		return nil, errors.New("Password must be 6 character or longer!")
	}

	if existingUser, _ := s.repo.GetByEmail(ctx, email); existingUser != nil {
		return nil, ErrEmailTaken
	}
	if existingUser, _ := s.repo.GetByUsername(ctx, username); existingUser != nil {
		return nil, ErrUsernameTaken
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	user := &model.User{Name: name, Username: username, Email: email, Password: string(hashedPassword)}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, name, username, email, password string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if username != "" && username != user.Username {
		if existing, _ := s.repo.GetByUsername(ctx, username); existing != nil {
			return nil, ErrUsernameTaken
		}
		user.Username = username
	}
	if email != "" && email != user.Email {
		if existing, _ := s.repo.GetByEmail(ctx, email); existing != nil {
			return nil, ErrEmailTaken
		}
		user.Email = email
	}
    if password != "" {
		if len(password) < 6 {
			return nil, ErrPasswordTooShort
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
        user.Password = string(hashedPassword) 
    }

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) DeleteUserPermanently(ctx context.Context, id uint) error {
	return s.repo.PermanentDelete(ctx, id)
}
