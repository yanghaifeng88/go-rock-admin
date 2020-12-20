package controller

import (
	"encoding/json"

	"github.com/daodao97/egin/consts"
	"github.com/daodao97/egin/db"
	"github.com/gin-gonic/gin"

	"oms/model"
)

// @Controller
type CommonConfig struct {
}

type CommonConfigFilter struct {
	Id int
}

// @GetApi /common_config/list
// @Summary 列表接口
// @Desc 列表接口 维护者: 刀刀
// @Params CommonConfigFilter
// @Response
// @Middleware IpLimiter
func (u CommonConfig) Get(c *gin.Context, params CommonConfigFilter) (interface{}, consts.ErrCode, error) {
	filter := db.Filter{}
	if params.Id != 0 {
		filter["id"] = params.Id
	}
	var result []model.CommonConfigEntity
	err := model.CommonConfig.Select(filter, db.Attr{
		Select:  []string{"id"},
		OrderBy: "id desc",
	}, &result)
	return result, 0, err
}

type CommonConfigForm struct {
}

// @PostApi /common_config/update
// @Summary 创建
// @Desc 创建接口 维护者: 刀刀
// @Params CommonConfigForm 接口参数所对应的结构体
// @Response
func (u CommonConfig) Post(c *gin.Context, params CommonConfigForm) (interface{}, consts.ErrCode, error) {
	var paramsMap db.Record
	tmp, _ := json.Marshal(params)
	err := json.Unmarshal(tmp, &paramsMap)
	lastId, _, err := model.CommonConfig.Insert(paramsMap)
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

// @PutApi /common_config/update/:id
// @Summary 更新接口
// @Desc 维护者: 刀刀
// @Params CommonConfigForm 接口参数所对应的结构体
// @Response
func (u CommonConfig) Put(c *gin.Context, id int, params CommonConfigForm) (interface{}, consts.ErrCode, error) {
	var paramsMap db.Record
	tmp, _ := json.Marshal(params)
	err := json.Unmarshal(tmp, &paramsMap)
	_, affected, err := model.CommonConfig.Update(db.Filter{"id": id}, paramsMap)
	var code consts.ErrCode
	if err != nil {
		code = consts.ErrorSystem
	}
	return affected, code, err
}

// @DeleteApi /common_config/delete/:id
// @Summary 删除
// @Desc 维护者: 刀刀
// @Response
func (u CommonConfig) Delete(c *gin.Context, id int) (interface{}, consts.ErrCode, error) {
	_, affected, err := model.CommonConfig.Delete(db.Filter{
		"id": id,
	})
	return affected, 0, err
}
