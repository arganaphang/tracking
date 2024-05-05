package handler

import (
	"application/internals/dto"
	"application/internals/model"
	"application/internals/repository"
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const IDENTITY_KEY = "id"

func JWTMiddleware(userRepo repository.IUserRepository) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "tracking system",
		Key:         []byte(IDENTITY_KEY),
		Timeout:     time.Hour * 10, // TODO: Update this
		MaxRefresh:  time.Hour * 24,
		IdentityKey: IDENTITY_KEY,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					IDENTITY_KEY: v.ID.String(),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id := uuid.MustParse(claims[IDENTITY_KEY].(string))
			return &model.User{
				ID: id,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var data dto.LoginRequest
			if err := c.ShouldBind(&data); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			user, err := userRepo.GetByEmail(c, data.Email)
			if err != nil {
				return nil, errors.New("incorrect Email or Password")
			}
			if !user.IsActive {
				return nil, errors.New("user is not active")
			}
			if !user.ComparePassword(data.Password) {
				return nil, errors.New("incorrect Email or Password")
			}
			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*model.User); ok {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:    "header: Authorization, cookie: jwt",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
		SendCookie:     true,
		CookieHTTPOnly: true,
	})
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}
