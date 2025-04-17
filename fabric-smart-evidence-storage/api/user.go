package api

import (
	"errors"
	"fabric-smart-evidence-storage/fabric"
	"fabric-smart-evidence-storage/model"
	"fabric-smart-evidence-storage/util"
	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册接口
func UserRegister(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err == nil {
		// 调用 fabric 注册接口
		encryptedPwd, err := util.PasswordEncrypt(user.Pwd)
		if err != nil {
			serverError(c, err)
			return
		}
		err = fabric.Register(
			user.Username,
			encryptedPwd,
			user.Nickname,
			user.Role,
		)
		if err != nil {
			serverError(c, err)
			return
		}
		c.Status(200)
	} else {
		serverError(c, err)
	}
}

type UserLoginVo struct {
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

// UserLogin 用户登录接口
func UserLogin(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err == nil {
		userChain, err := fabric.GetUser(user.Username)
		if err != nil {
			serverError(c, err)
			return
		}
		if !util.PasswordMatch(user.Pwd, userChain.Pwd) {
			serverError(c, errors.New("密码错误"))
			return
		}
		token, err := util.GenerateJWT(user.Username, userChain.Role)
		if err != nil {
			serverError(c, err)
			return
		}

		c.JSON(200, UserLoginVo{
			Token:    token,
			Nickname: userChain.Nickname,
			Role:     userChain.Role,
		})
	} else {
		serverError(c, err)
	}
}

// 用户列表
func UserList(c *gin.Context) {
	// 调用 fabric 获取所有用户
	users, err := fabric.GetAllUsers()
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(200, users)
}

// 获取可授权用户列表，除掉自己
func GetAuthorizeUserList(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		serverError(c, errors.New("用户未登录"))
		return
	}
	users, err := fabric.GetAllUsers()
	if err != nil {
		serverError(c, err)
		return
	}
	userList := make([]string, 0)
	for _, user := range users {
		if user.Username != username {
			userList = append(userList, user.Username)
		}
	}
	c.JSON(200, userList)
}

func UserUpdate(c *gin.Context) {
	ao := new(model.User)
	if err := c.ShouldBind(ao); err == nil {
		// 调用 fabric 更新用户状态和角色
		err = fabric.UpdateRoleAndStatus(ao.Username, ao.Role, ao.Status)
		if err != nil {
			serverError(c, err)
			return
		}
		c.Status(200)
	} else {
		serverError(c, err)
	}
}

// 自己修改呢称或者密码
type UserUpdateMeAo struct {
	OldPassword string `json:"oldPwd"`
	Password    string `json:"pwd"`
}

func UserChangePwd(c *gin.Context) {
	ao := new(UserUpdateMeAo)
	if err := c.ShouldBind(ao); err == nil {
		username, exists := c.Get("username")
		if !exists {
			serverError(c, errors.New("用户未登录"))
			return
		}

		if ao.OldPassword == "" && ao.Password == "" {
			serverError(c, errors.New("参数错误"))
			return
		}
		userChain, err := fabric.GetUser(username.(string))
		if err != nil {
			serverError(c, err)
			return
		}
		if !util.PasswordMatch(ao.OldPassword, userChain.Pwd) {
			serverError(c, errors.New("旧密码错误"))
			return
		}
		encryptedPwd, err := util.PasswordEncrypt(ao.Password)
		if err != nil {
			serverError(c, err)
			return
		}
		err = fabric.UpdatePwd(username.(string), encryptedPwd)
		if err != nil {
			serverError(c, err)
			return
		}
		c.Status(200)
	} else {
		serverError(c, err)
	}
}
