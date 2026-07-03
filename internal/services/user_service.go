package services

import (
	"errors"

	"github.com/rodatboat/crong/internal/middleware"
	"github.com/rodatboat/crong/internal/utils"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/database"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
	"github.com/rodatboat/crong/internal/resp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(email string, firstName string, lastName string, password string) (*models.User, error) {
	log.Infof("Registering new user with email %v", email)

	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	userEntity, err := s.userRepo.Create(email, firstName, lastName, hashedPassword)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return nil, resp.ErrUserAlreadyExists
		}
		return nil, err
	}

	// TODO: Send verification email

	return utils.MapUserEntityToUserModel(userEntity), nil
}

func (s *UserService) LoginUser(email string, password string) (*models.User, error) {
	log.Infof("Logging in user with email %v", email)

	// Find user by email
	userEntity, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrNotFound
		}
		return nil, err
	}

	// TODO: Validate user email is verified

	// Verify password
	if VerifyPassword(password, userEntity.Password) == false {
		return nil, resp.ErrBadRequest
	}

	// Generate JWT
	token, err := middleware.GenerateJWT(userEntity.ID, userEntity.Email)
	if err != nil {
		log.Errorf("Failed to generate JWT: %v", err)
		return nil, err
	}

	user := utils.MapUserEntityToUserModel(userEntity)
	user.AuthToken = token

	return user, nil
}

func (s *UserService) UpdateUser(userID uint, updateReq models.UserUpdate) (*models.User, error) {
	log.Infof("Updating user with id %v", userID)

	// Find user by email
	userEntity, err := s.userRepo.Update(userID, updateReq.FirstName, updateReq.LastName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrNotFound
		}
		return nil, err
	}

	user := utils.MapUserEntityToUserModel(userEntity)

	return user, nil
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	// pepper := os.Getenv("PASSWORD_PEPPER")
	pepper := ""

	pepperedPassword := password + pepper

	bytes, err := bcrypt.GenerateFromPassword([]byte(pepperedPassword), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password string, hashedPassword string) bool {
	// pepper := os.Getenv("PASSWORD_PEPPER")
	pepper := ""

	pepperedPassword := password + pepper

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pepperedPassword))
	return err == nil
}
