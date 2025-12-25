package controllers

import (
	"fmt"
	"gindemo2/config"
	"gindemo2/models"
	"gindemo2/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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
	if err := config.DB.Create(&user).Error; err != nil {
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
	err := config.DB.Where(`username = ?`, input.Username).First(&user).Error
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
	var user models.User
	uid := c.GetUint(`user_id`)
	fmt.Println(`uid=`, uid)
	err := config.DB.Where(`id=?`, uid).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			`error`: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
