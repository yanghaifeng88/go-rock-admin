package controller

import (
	"encoding/json"
	"fmt"

	"github.com/daodao97/egin/consts"
	"github.com/daodao97/egin/db"
	egin_user "github.com/daodao97/egin/service/user"
	"github.com/daodao97/egin/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"oms/model"
	"oms/service/menu"
	"oms/service/user"
)

// @Controller 用户管理 这里是简介
type User struct {
}

type UserFilter struct {
	Id     int    `form:"id"`
	Name   string `form:"name"`
	RoleId int    `form:"role_id"`
	Status int    `form:"status"`
}

type UserForm struct {
	Id       int    `form:"id" json:"id"`
	Name     string `form:"name" json:"name" comment:"名称"`
	Avatar   string `form:"avatar" json:"avatar" comment:"头像"`
	Password string `form:"password" json:"password" comment:"密码"`
	Status   int    `form:"status" json:"status" comment:"状态 0禁用 1启用"`
	Email    string `form:"email" json:"email" comment:"邮箱"`
	Mobile   string `form:"mobile" json:"mobile" comment:"手机号"`
	RoleIds  [][]int `form:"role_ids" json:"role_ids" comment:"用户角色"`
}

// @PostApi /user/create 若无默认为全小写的方法名
// @Summary 创建用户
// @Desc 这个接口支持你创建一个用户 维护者: 刀刀
// @Params UserForm 接口参数所对应的结构体
// @Response
func (u User) Post(c *gin.Context, newUser UserForm) (interface{}, consts.ErrCode, error) {
	var userMap db.Record
	tmp, _ := json.Marshal(newUser)
	err := json.Unmarshal(tmp, &userMap)
	lastId, _, err := model.User.Insert(userMap)
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	var result struct {
		LastId int64 `json:"last_id"`
	}
	result.LastId = lastId
	return result, code, err
}

type UserLogin struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Ticket   string
}

// @PostApi /user/login
// @Summary 删除用户
// @Desc 维护者: 刀刀
// @Params UserLogin 接口参数所对应的结构体
// @Response
func (u User) Login(c *gin.Context, userPost UserLogin) (interface{}, consts.ErrCode, error) {
	info, err := user.New().Login(userPost.Name, userPost.Password)
	if err != nil {
		return nil, 500, err
	}
	token, err := utils.GenerateToken(info.Id, info.Name)
	if err != nil {
		return nil, 500, errors.New("gen token error")
	}
	return map[string]interface{}{
		"name":   info.Name,
		"avatar": "",
		"token":  fmt.Sprintf("Bearer %v", token),
	}, 0, nil
}

// @GetApi /user/info
func (u User) Info(c *gin.Context) (interface{}, consts.ErrCode, error) {
	info, exists := c.Get("user")
	if !exists {
		return nil, 500, errors.New("token error")
	}
	me := info.(egin_user.Info)
	full, err := user.FullAttr(me.Id)
	if err != nil {
		return nil, 500, err
	}
	return full, 0, nil
}

// @PostApi /user/logout
func (u User) Logout(c *gin.Context) (interface{}, consts.ErrCode, error) {
	return "success", 0, nil
}

// @GetApi /user/routes
func (u User) Routes(c *gin.Context) (interface{}, consts.ErrCode, error) {
	info, exists := c.Get("user")
	if !exists {
		return nil, 500, nil
	}
	me := info.(egin_user.Info)

	userRoute := menu.GetPages(me.Id)

	return userRoute, 0, nil
}

// @GetApi /user/list 列表接口
func (u User) List(c *gin.Context, params UserFilter) (interface{}, consts.ErrCode, error) {
	filter := db.Filter{}
	if params.Id != 0 {
		filter["id"] = params.Id
	}
	if params.Name != "" {
		filter["name"] = map[string]string{
			"like": params.Name + "%",
		}
	}
	_, exists := c.Request.Form["status"]
	if exists {
		filter["status"] = params.Status
	}
	var result []model.UserEntity
	err := model.User.Select(filter, db.Attr{
		Select:  []string{"name", "id", "avatar", "email", "mobile", "status"},
		OrderBy: "id desc",
	}, &result)
	response := ListResponse{
		List: result,
		Page: Page{Count: 1, Page: 1},
	}

	return response, 0, err
}

// @GetApi /user/get/:id 若无默认为全小写的方法名
// @Summary 用户列表接口
// @Desc 接口简介, 若无则为空 维护者: 刀刀
// @Params UserFilter 接口参数所对应的结构体
// @Response
// @Middleware IpLimiter
func (u User) Get(c *gin.Context, id int) (interface{}, consts.ErrCode, error) {
	var result model.UserEntity
	err := model.User.FindById(id, []string{"id", "name", "avatar", "name", "email", "mobile", "status", "role_ids"}, &result)
	return result, 0, err
}

// @PostApi /user/update/:id
// @Summary 更新用户信息
// @Desc 维护者: 刀刀
// @Params UserForm 接口参数所对应的结构体
// @Response
func (u User) Put(c *gin.Context, id int, upUser UserForm) (interface{}, consts.ErrCode, error) {
	var userMap db.Record
	tmp, _ := json.Marshal(upUser)
	err := json.Unmarshal(tmp, &userMap)
	delete(userMap, "password")
	_, affected, err := model.User.Update(db.Filter{"id": id}, userMap)
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	return affected, code, err
}

// @DeleteApi /user/delete/:id
// @Summary 删除用户
// @Desc 维护者: 刀刀
// @Response
func (u User) Delete(c *gin.Context, id int) (interface{}, consts.ErrCode, error) {
	_, affected, err := model.User.Delete(db.Filter{
		"id": id,
	})
	return affected, 0, err
}
