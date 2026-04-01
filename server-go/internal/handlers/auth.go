package handlers

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"bycigar-server/internal/config"
	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginAttempt struct {
	count       int
	lastAttempt time.Time
}

var (
	loginFailures = make(map[string]*loginAttempt)
	failureMutex  sync.Mutex
)

const maxLoginAttempts = 3

func getRequireCaptcha(key string) bool {
	failureMutex.Lock()
	defer failureMutex.Unlock()
	a, ok := loginFailures[key]
	if !ok {
		return false
	}
	if time.Since(a.lastAttempt) > 15*time.Minute {
		delete(loginFailures, key)
		return false
	}
	return a.count >= maxLoginAttempts
}

func recordFailure(key string) {
	failureMutex.Lock()
	defer failureMutex.Unlock()
	a, ok := loginFailures[key]
	if !ok {
		a = &loginAttempt{}
		loginFailures[key] = a
	}
	a.count++
	a.lastAttempt = time.Now()
}

func clearFailures(key string) {
	failureMutex.Lock()
	defer failureMutex.Unlock()
	delete(loginFailures, key)
}

// Register godoc
// @Summary 用户注册
// @Description 注册新用户账号
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RegisterInput true "注册信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请填写完整信息")
		return
	}

	if !captchaInstance.Verify(input.CaptchaId, input.CaptchaCode, true) {
		utils.ErrorResponse(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	if len(input.Password) < 6 {
		utils.ErrorResponse(c, http.StatusBadRequest, "密码至少需要6个字符")
		return
	}

	var existing models.User
	if database.DB.Where("email = ?", input.Email).First(&existing).Error == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "该邮箱已被注册")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	name := input.Name
	if name == "" {
		idx := strings.Index(input.Email, "@")
		if idx > 0 {
			name = input.Email[:idx]
		}
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Name:     name,
		Role:     "customer",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录获取token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "登录信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "邮箱和密码不能为空")
		return
	}

	key := strings.ToLower(input.Email)
	needCaptcha := getRequireCaptcha(key)

	if needCaptcha {
		if input.CaptchaId == "" || input.CaptchaCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          "请输入验证码",
				"requireCaptcha": true,
			})
			return
		}
		if !captchaInstance.Verify(input.CaptchaId, input.CaptchaCode, true) {
			recordFailure(key)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          "验证码错误或已过期",
				"requireCaptcha": true,
			})
			return
		}
	}

	var user models.User
	loginEmail := strings.ToLower(input.Email)
	if err := database.DB.Where("email = ?", loginEmail).First(&user).Error; err != nil {
		recordFailure(key)
		requireCaptcha := getRequireCaptcha(key)
		resp := gin.H{"error": "邮箱或密码错误"}
		if requireCaptcha {
			resp["requireCaptcha"] = true
		}
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		recordFailure(key)
		requireCaptcha := getRequireCaptcha(key)
		resp := gin.H{"error": "邮箱或密码错误"}
		if requireCaptcha {
			resp["requireCaptcha"] = true
		}
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	clearFailures(key)

	if user.IsBanned {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被封禁，请联系客服"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// GetProfile godoc
// @Summary 获取用户信息
// @Description 获取当前登录用户的个人信息
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/me [get]
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"name":     user.Name,
			"role":     user.Role,
			"isBanned": user.IsBanned,
		},
	})
}

// UpdateProfile godoc
// @Summary 更新用户信息
// @Description 更新当前用户的个人信息
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.UpdateProfileInput true "用户信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [put]
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not found")
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input models.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请填写完整信息")
		return
	}

	if !captchaInstance.Verify(input.CaptchaId, input.CaptchaCode, true) {
		utils.ErrorResponse(c, http.StatusBadRequest, "验证码错误或已过期")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "原密码错误")
		return
	}

	if input.OldPassword == input.NewPassword {
		utils.ErrorResponse(c, http.StatusBadRequest, "新密码不能与原密码相同")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码加密失败")
		return
	}
	user.Password = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码更新失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
