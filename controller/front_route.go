package controller

import (
	"encoding/json"

	"github.com/daodao97/egin/consts"
	"github.com/daodao97/egin/db"
	userInstance "github.com/daodao97/egin/service/user"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"oms/model"
	"oms/service/menu"
)

// @Controller
type FrontRoute struct {
}

type FrontRouteFilter struct {
	Pid    int    `form:"pid"`
	Status int    `form:"status"`
	Name   string `form:"name"`
	Filter
}

type RoutesList struct {
	model.FrontRouteEntity
	Children    []interface{} `json:"children"`
	HasChildren bool          `json:"hasChildren"`
}

// @GetApi /front_route/list
func (f FrontRoute) List(c *gin.Context, params FrontRouteFilter) (interface{}, consts.ErrCode, error) {
	filter := db.Filter{}
	filter["pid"] = params.Pid
	_, exists := c.Request.Form["status"]
	if exists {
		filter["status"] = params.Status
	}
	var result []RoutesList
	err := model.FrontRoute.Select(filter, db.Attr{
		Select:  []string{"id", "name", "path", "icon", "sort", "status", "type", "page_type", "view"},
		OrderBy: "sort desc",
	}, &result)
	for i, v := range result {
		result[i].Children = make([]interface{}, 0)
		result[i].HasChildren = v.Type != 3
	}
	response := ListResponse{
		List: result,
		Page: Page{Count: 1, Page: 1},
	}

	return response, 0, err
}

// @GetApi /front_route/get/:id
// @Summary 列表接口
// @Desc 列表接口 维护者: 刀刀
// @Params FrontRouteFilter
// @Response
// @Middleware IpLimiter
func (f FrontRoute) Get(c *gin.Context, id int) (interface{}, consts.ErrCode, error) {
	filter := db.Filter{"id": id}
	var result []model.FrontRouteEntity
	err := model.FrontRoute.Select(filter, db.Attr{
		Select:  []string{"id", "name", "path", "icon", "sort", "status", "type", "page_type", "page_schema", "pid"},
		OrderBy: "id desc",
	}, &result)
	if len(result) != 1 {
		return nil, 500, errors.New("select error")
	}
	return result[0], 0, err
}

type FrontRouteForm struct {
	Id         int    `form:"id" json:"id"`
	Pid        int    `form:"pid" json:"pid" comment:"父id"`
	ModuleId   int    `form:"module_id" json:"module_id" comment:"模块id"`
	Name       string `form:"name" json:"name" comment:"路由名"`
	Type       int    `form:"type" json:"type" comment:"类型  1 目录, 2 菜单, 3 页面"`
	Path       string `form:"path" json:"path" comment:"前端路由"`
	Icon       string `form:"icon" json:"icon" comment:"图标"`
	PageType   int    `form:"page_type" json:"page_type" comment:"页面类型 0 自定义, 1 表单页, 2 列表页, 3复杂schema"`
	PageSchema string `form:"page_schema" json:"page_schema" comment:"页面定义"`
	View       string `form:"view" json:"view" comment:"自定义组价路径"`
	Sort       int    `form:"sort" json:"sort" comment:"倒序排序"`
	Status     int    `form:"status" json:"status" comment:"状态 0 禁用, 1 启用"`
}

// @PostApi /front_route/create
// @Summary 创建
// @Desc 创建接口 维护者: 刀刀
// @Params FrontRouteForm 接口参数所对应的结构体
// @Response
func (f FrontRoute) Post(c *gin.Context, params FrontRouteForm) (interface{}, consts.ErrCode, error) {
	var paramsMap db.Record
	tmp, _ := json.Marshal(params)
	err := json.Unmarshal(tmp, &paramsMap)
	lastId, _, err := model.FrontRoute.Insert(paramsMap)
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	var result struct {
		LastId int64
	}
	result.LastId = lastId
	return result, code, err
}

// @PostApi /front_route/update/:id
// @Summary 更新接口
// @Desc 维护者: 刀刀
// @Params FrontRouteForm 接口参数所对应的结构体
// @Response
func (f FrontRoute) Put(c *gin.Context, id int, params FrontRouteForm) (interface{}, consts.ErrCode, error) {
	var paramsMap db.Record
	tmp, _ := json.Marshal(params)
	err := json.Unmarshal(tmp, &paramsMap)
	_, affected, err := model.FrontRoute.Update(db.Filter{"id": id}, paramsMap)
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	return map[string]interface{}{"affected": affected}, code, err
}

// @DeleteApi /front_route/delete:id
// @Summary 删除
// @Desc 维护者: 刀刀
// @Response
func (f FrontRoute) Delete(c *gin.Context, id int) (interface{}, consts.ErrCode, error) {
	_, affected, err := model.FrontRoute.Delete(db.Filter{
		"id": id,
	})
	return affected, 0, err
}

// @GetApi /front_route/tree
func (f FrontRoute) Tree(c *gin.Context) (interface{}, consts.ErrCode, error) {
	info, exists := c.Get("user")
	if !exists {
		return nil, 500, nil
	}
	me := info.(userInstance.Info)

	tree := menu.GetMenuSelect(me.Id)

	return tree, 0, nil
}
