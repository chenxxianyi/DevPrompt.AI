package service

import (
	"errors"
	"time"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/utils"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtSecret  string
	jwtExpireH int
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpireH int) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		jwtExpireH: jwtExpireH,
	}
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResult struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register 注册新用户
func (s *AuthService) Register(input *RegisterInput) (*AuthResult, error) {
	// 检查邮箱是否已注册
	existing, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("该邮箱已注册")
	}

	// 检查用户名是否已占用
	existing, err = s.userRepo.FindByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("该用户名已被占用")
	}

	// 密码哈希
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         "user",
		MembershipLevel: "free",
		Status:       "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 生成 JWT
	token, err := utils.GenerateToken(s.jwtSecret, s.jwtExpireH, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: user}, nil
}

// Login 用户登录
func (s *AuthService) Login(input *LoginInput) (*AuthResult, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("邮箱或密码错误")
	}

	if user.Status == "disabled" {
		return nil, errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return nil, errors.New("邮箱或密码错误")
	}

	token, err := utils.GenerateToken(s.jwtSecret, s.jwtExpireH, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: user}, nil
}

// GetProfile 获取用户信息
func (s *AuthService) GetProfile(userID uint64) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// GetMembershipExpiry 获取会员过期时间（格式化）
func GetMembershipExpiry(user *model.User) *string {
	return user.MembershipExpiredAt
}

// CheckMembership 检查会员是否有效
func CheckMembership(user *model.User) bool {
	if user.MembershipLevel == "free" {
		return true
	}
	if user.MembershipExpiredAt == nil {
		return false
	}
	expiry, err := time.Parse("2006-01-02 15:04:05", *user.MembershipExpiredAt)
	if err != nil {
		return false
	}
	return time.Now().Before(expiry)
}
