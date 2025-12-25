package controllers

import (
	"golang_per_day_30/internal/components"
	"golang_per_day_30/internal/errors"
	"golang_per_day_30/internal/models"
	"golang_per_day_30/internal/services"
	"golang_per_day_30/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type inputUser struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Passwd   string `json:"passwd" binding:"required,min=6"`
}

// 注册用户
func Register(c *gin.Context) {
	input := inputUser{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			`error`: err.Error(),
		})
		return
	}
	hasedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Passwd), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			`error`: "哈希密码失败",
		})
		return
	}
	user := models.User{
		Username: input.Username,
		PassWd:   string(hasedPwd),
	}
	if err := components.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{`error`: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		`message`: "用户创建成功",
		`user`: gin.H{
			`id`:       user.ID,
			`username`: user.Username,
		},
	})
}

// 登录
func Login(c *gin.Context) {
	input := inputUser{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			`error`: err.Error(),
		})
		return
	}
	var user models.User
	err := components.DB.Where(`username = ?`, input.Username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PassWd), []byte(input.Passwd))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// 查询当前登录者的信息
func Me(c *gin.Context) {
	user, err := services.User.GetById(c)
	if err != nil {
		// 业务错误
		if biz, ok := err.(*errors.BizError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": biz.Code,
				"msg":  biz.Msg,
			})
			return
		}
		fields := []zap.Field{zap.String("trace_id", c.GetString(`trace_id`))}
		zap.L().Error(err.Error(), fields...)
		c.JSON(http.StatusInternalServerError, gin.H{
			`code`: http.StatusInternalServerError,
			`msg`:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "查询成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
