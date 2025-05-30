package staff

import (
	"errors"
	"opsalert/internal/jwt"
	staffModel "opsalert/internal/model/staff"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInactiveAccount    = errors.New("account is inactive")
	ErrUserNotFound       = errors.New("user not found")
)

type Service struct {
	repo       Repository
	jwtService *jwt.Service
}

type Repository interface {
	Create(staff *staffModel.Staff) error
	GetByUsername(username string) (*staffModel.Staff, error)
	GetByID(id uint) (*staffModel.Staff, error)
	GetAll() ([]staffModel.Staff, error)
	Update(id uint, staff *staffModel.Staff) error
}

func NewService(repo Repository, jwtService *jwt.Service) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *Service) Register(req *staffModel.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newStaff := &staffModel.Staff{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         req.Role,
		IsActive:     true,
	}

	return s.repo.Create(newStaff)
}

func (s *Service) Login(req *staffModel.LoginRequest) (string, *staffModel.Staff, error) {
	staff, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !staff.IsActive {
		return "", nil, ErrInactiveAccount
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(staff.ID, staff.Username, staff.Role)
	if err != nil {
		return "", nil, err
	}

	return token, staff, nil
}

func (s *Service) GetProfile(userID uint) (*staffModel.Staff, error) {
	staff, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return staff, nil
}

func (s *Service) GetAccounts() ([]staffModel.Staff, error) {
	return s.repo.GetAll()
}

func (s *Service) UpdateStaff(id uint, req *staffModel.UpdateStaffRequest) error {
	staff, err := s.repo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	staff.FullName = req.FullName
	staff.Role = req.Role
	staff.IsActive = req.IsActive

	return s.repo.Update(id, staff)
}
