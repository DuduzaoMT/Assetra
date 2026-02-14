package auth

import (
	"assetra/authentication/models"
	"assetra/authentication/repository"
	"assetra/authentication/validators"
	"assetra/pb"
	"assetra/security"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type authService struct {
	pb.UnimplementedAuthServiceServer
	usersRepository         repository.UserRepository
	refreshTokenRepository repository.RefreshTokenRepository
}

func NewAuthService(usersRepository repository.UserRepository, refreshTokenRepository repository.RefreshTokenRepository) pb.AuthServiceServer {
	return &authService{
		usersRepository:         usersRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

// Implementation stubs
func (s *authService) SignUp(ctx context.Context, user *pb.User) (*pb.SignInResponse, error) {
	log.Printf("[SECURITY] SignUp attempt for email: %s, username: %s", user.Email, user.Name)

	if err := validators.ValidateSignUp(user); err != nil {
		log.Printf("[SECURITY] SignUp validation failed for %s: %v", user.Email, err)
		return nil, err
	}

	user.Name = validators.SanitizeName(strings.TrimSpace(user.Name))
	user.Email = validators.NormalizeEmail(user.Email)

	// As we are creating a new user we need to compare to all of them
	err := s.checkUserUniqueness("", user.Email, user.Name)
	if err != nil {
		log.Printf("[SECURITY] SignUp failed - user already exists: %s", user.Email)
		return nil, err
	}

	// encrypt the password before storing it
	user.Password, err = security.EncryptPassword(user.Password)
	if err != nil {
		log.Printf("[ERROR] Password encryption failed: %v", err)
		return nil, errors.New("failed to create user")
	}

	savedUser := new(models.User)
	savedUser.FromProtoBuffer(user)
	newuser, err := s.usersRepository.Save(savedUser)
	if err != nil {
		log.Printf("[ERROR] Failed to save user %s: %v", user.Email, err)
		return nil, errors.New("failed to create user")
	}

	// Generate access token (15 minutes)
	accessToken, err := security.NewAccessToken(newuser.ID.String(), newuser.Roles)
	if err != nil {
		log.Printf("[ERROR] Access token generation failed for user %s: %v", newuser.Email, err)
		return nil, errors.New("authentication failed")
	}

	// Generate refresh token
	refreshToken, err := security.NewRefreshToken()
	if err != nil {
		log.Printf("[ERROR] Refresh token generation failed for user %s: %v", newuser.Email, err)
		return nil, errors.New("authentication failed")
	}

	// Save refresh token hash to database
	refreshTokenModel := &models.RefreshToken{
		UserID:    newuser.ID,
		TokenHash: security.HashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	if err := s.refreshTokenRepository.Save(refreshTokenModel); err != nil {
		log.Printf("[ERROR] Failed to save refresh token for user %s: %v", newuser.Email, err)
		return nil, errors.New("authentication failed")
	}

	log.Printf("[SECURITY] Successful signup for user: %s (ID: %s)", newuser.Email, newuser.ID.String())
	
	return &pb.SignInResponse{
		User:         newuser.ToProtoBuffer(),
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	log.Printf("[SECURITY] SignIn attempt for email: %s", req.Email)
	req.Email = validators.NormalizeEmail(req.Email)

	// Check if account is locked
	isLocked, err := s.usersRepository.IsAccountLocked(req.Email)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("[ERROR] Failed to check account lock status: %v", err)
		return nil, errors.New("authentication failed")
	}
	
	if isLocked {
		log.Printf("[SECURITY] Login attempt for locked account: %s", req.Email)
		return nil, errors.New("account locked due to multiple failed login attempts")
	}

	user, err := s.usersRepository.FindByEmail(req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("[SECURITY] Failed login attempt - user not found: %s", req.Email)
			// Don't reveal whether user exists
			return nil, errors.New("invalid credentials")
		}
		log.Printf("[ERROR] Database error during signin for %s: %v", req.Email, err)
		return nil, errors.New("authentication failed")
	}

	// Compare the provided password with the stored hashed password
	if err := security.ComparePassword(user.Password, req.Password); err != nil {
		log.Printf("[SECURITY] Failed login attempt - invalid password for user: %s", req.Email)
		
		// Increment failed login attempts
		if err := s.usersRepository.IncrementFailedLoginAttempts(req.Email); err != nil {
			log.Printf("[ERROR] Failed to increment login attempts: %v", err)
		}
		
		// Check if we need to lock the account (after 3 failed attempts)
		updatedUser, _ := s.usersRepository.FindByEmail(req.Email)
		if updatedUser != nil && updatedUser.FailedLoginAttempts >= 3 {
			lockUntil := time.Now().Add(15 * time.Minute)
			if err := s.usersRepository.LockAccount(req.Email, lockUntil); err != nil {
				log.Printf("[ERROR] Failed to lock account: %v", err)
			} else {
				log.Printf("[SECURITY] Account locked for 15 minutes due to failed attempts: %s", req.Email)
			}
		}
		
		// Don't reveal whether password is wrong
		return nil, errors.New("invalid credentials")
	}

	// Reset failed login attempts on successful login
	if err := s.usersRepository.ResetFailedLoginAttempts(req.Email); err != nil {
		log.Printf("[ERROR] Failed to reset login attempts: %v", err)
	}

	// Generate access token (15 minutes)
	accessToken, err := security.NewAccessToken(user.ID.String(), user.Roles)
	if err != nil {
		log.Printf("[ERROR] Access token generation failed for user %s: %v", req.Email, err)
		return nil, errors.New("authentication failed")
	}

	// Generate refresh token
	refreshToken, err := security.NewRefreshToken()
	if err != nil {
		log.Printf("[ERROR] Refresh token generation failed for user %s: %v", req.Email, err)
		return nil, errors.New("authentication failed")
	}

	// Save refresh token hash to database
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: security.HashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	if err := s.refreshTokenRepository.Save(refreshTokenModel); err != nil {
		log.Printf("[ERROR] Failed to save refresh token for user %s: %v", req.Email, err)
		return nil, errors.New("authentication failed")
	}

	log.Printf("[SECURITY] Successful login for user: %s (ID: %s)", req.Email, user.ID.String())
	
	return &pb.SignInResponse{
		User:         user.ToProtoBuffer(),
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, validators.ErrInvalidUserId
	}

	user, err := s.usersRepository.FindByID(req.Id)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, fmt.Errorf("error finding user by ID: %w", err)
		}
		return nil, fmt.Errorf("no user found with id: %s", req.Id)
	}
	return user.ToProtoBuffer(), nil
}

func (s *authService) ListUsers(req *pb.ListUsersRequest, stream pb.AuthService_ListUsersServer) error {

	users, err := s.usersRepository.FindAll()
	if err != nil {
		if err != pgx.ErrNoRows {
			return fmt.Errorf("error listing users: %w", err)
		}
		return errors.New("no users found")
	}

	for _, user := range users {
		if err := stream.Send(user.ToProtoBuffer()); err != nil {
			return fmt.Errorf("error sending user to stream: %w", err)
		}
	}
	return nil
}

// the user only can update the password , username and email
// the roles must be updated from someone with acess to the DB
func (s *authService) UpdateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	if err := validators.ValidateUpdateUser(user); err != nil {
		log.Printf("[SECURITY] UpdateUser validation failed for ID %s: %v", user.Id, err)
		return nil, err
	}

	existingUser, err := s.usersRepository.FindByID(user.Id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no user found with id: %s", user.Id)
		}
		return nil, fmt.Errorf("error finding user by ID: %w", err)
	}

	// Only update fields that were provided
	if user.Password != "" {
		encryptedPassword, err := security.EncryptPassword(user.Password)
		if err != nil {
			log.Printf("[ERROR] Password encryption failed: %v", err)
			return nil, errors.New("failed to update user")
		}
		existingUser.Password = encryptedPassword
	}

	if user.Name != "" && existingUser.Username != user.Name {
		existingUser.Username = validators.SanitizeName(strings.TrimSpace(user.Name))
	}

	if user.Email != "" {
		existingUser.Email = validators.NormalizeEmail(user.Email)
	}

	existingUser.UpdatedAt = time.Now()

	// Check uniqueness only for changed fields
	err = s.checkUserUniqueness(existingUser.ID.String(), existingUser.Email, existingUser.Username)
	if err != nil {
		log.Printf("[SECURITY] Update failed - duplicate user: %v", err)
		return nil, err
	}

	if err := s.usersRepository.Update(existingUser); err != nil {
		log.Printf("[ERROR] Failed to update user %s: %v", user.Id, err)
		return nil, errors.New("failed to update user")
	}

	log.Printf("[SECURITY] User updated successfully: %s (ID: %s)", existingUser.Email, existingUser.ID.String())
	return existingUser.ToProtoBuffer(), nil
}

func (s *authService) DeleteUser(ctx context.Context, req *pb.GetUserRequest) (*pb.DeleteUserResponse, error) {
	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, validators.ErrInvalidUserId
	}

	err := s.usersRepository.Delete(req.Id)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, fmt.Errorf("error finding user by ID: %w", err)
		}
		return nil, fmt.Errorf("no user found with id: %s", req.Id)
	}
	return nil, nil
}

func (s *authService) checkUserUniqueness(id, email, username string) error {
	found, err := s.usersRepository.FindByEmail(email)
	if err != nil {
		// Error
		if err != pgx.ErrNoRows {
			return fmt.Errorf("error finding user by ID: %w", err)
		}
		// err == ErrnoRows
	}

	if found != nil && found.ID.String() != id {
		return errors.New("there is already one user with this email registered")
	}

	found, err = s.usersRepository.FindByUserName(username)
	if err != nil {
		// Error
		if err != pgx.ErrNoRows {
			return fmt.Errorf("error finding user by ID: %w", err)
		}
		// err == ErrnoRows 
	}

	if found != nil && found.ID.String() != id {
		return errors.New("there is already one user with this username registered")
	}

	return nil
}

// RefreshToken validates a refresh token and issues a new access token
func (s *authService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	log.Printf("[SECURITY] Refresh token request")

	// Hash the provided refresh token to search in database
	tokenHash := security.HashRefreshToken(req.RefreshToken)

	// Find the refresh token in the database
	storedToken, err := s.refreshTokenRepository.FindByTokenHash(tokenHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("[SECURITY] Refresh token not found or invalid")
			return nil, errors.New("invalid refresh token")
		}
		log.Printf("[ERROR] Error finding refresh token: %v", err)
		return nil, errors.New("authentication failed")
	}

	// Get the user
	user, err := s.usersRepository.FindByID(storedToken.UserID.String())
	if err != nil {
		log.Printf("[ERROR] User not found for refresh token: %v", err)
		return nil, errors.New("authentication failed")
	}

	// Generate new access token
	accessToken, err := security.NewAccessToken(user.ID.String(), user.Roles)
	if err != nil {
		log.Printf("[ERROR] Access token generation failed: %v", err)
		return nil, errors.New("authentication failed")
	}

	// Optional: Refresh Token Rotation
	// Revoke the old refresh token
	if err := s.refreshTokenRepository.RevokeToken(tokenHash); err != nil {
		log.Printf("[ERROR] Failed to revoke old refresh token: %v", err)
		// Continue anyway, not critical
	}

	// Generate new refresh token
	newRefreshToken, err := security.NewRefreshToken()
	if err != nil {
		log.Printf("[ERROR] New refresh token generation failed: %v", err)
		return nil, errors.New("authentication failed")
	}

	// Save new refresh token to database
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: security.HashRefreshToken(newRefreshToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	if err := s.refreshTokenRepository.Save(refreshTokenModel); err != nil {
		log.Printf("[ERROR] Failed to save new refresh token: %v", err)
		return nil, errors.New("authentication failed")
	}

	log.Printf("[SECURITY] Token refreshed successfully for user: %s (ID: %s)", user.Email, user.ID.String())

	return &pb.RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
