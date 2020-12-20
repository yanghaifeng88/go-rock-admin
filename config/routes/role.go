// ****************************
// 该文件为系统生成, 请勿更改
// ****************************
package routes

import (
	"strconv"
	"strings"

	"github.com/daodao97/egin"
	"github.com/daodao97/egin/consts"
	"github.com/daodao97/egin/middleware"
	"github.com/daodao97/egin/utils"
	"github.com/gin-gonic/gin"

	"oms/controller"
)

func RegRoleRouter(r *gin.Engine) {
	ctrl := controller.Role{}

	r.Handle("GET", "/role/list", func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			var params controller.RoleFilter
			errs := utils.Validated(ctx, &params)
			if errs != nil {
				egin.Fail(ctx, consts.ErrorParam, strings.Join(errs, "\n"))
				return
			}
			result, code, err := ctrl.List(ctx, params)
			egin.Response(ctx, result, code, err)
		}
	}())

	r.Handle("GET", "/role/get/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		result, code, err := ctrl.Get(ctx, id)
		egin.Response(ctx, result, code, err)
	}, middleware.IpLimiter())

	r.Handle("POST", "/role/create", func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			var params controller.RoleForm
			errs := utils.Validated(ctx, &params)
			if errs != nil {
				egin.Fail(ctx, consts.ErrorParam, strings.Join(errs, "\n"))
				return
			}
			result, code, err := ctrl.Post(ctx, params)
			egin.Response(ctx, result, code, err)
		}
	}())

	r.Handle("POST", "/role/update/:id", func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			var params controller.RoleUpdate
			errs := utils.Validated(ctx, &params)
			if errs != nil {
				egin.Fail(ctx, consts.ErrorParam, strings.Join(errs, "\n"))
				return
			}

			id, _ := strconv.Atoi(ctx.Param("id"))

			result, code, err := ctrl.Put(ctx, id, params)
			egin.Response(ctx, result, code, err)
		}
	}())

	r.Handle("DELETE", "/role/delete/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		result, code, err := ctrl.Delete(ctx, id)
		egin.Response(ctx, result, code, err)
	})

	r.Handle("GET", "/role/tree", func(ctx *gin.Context) {
		result, code, err := ctrl.Tree(ctx)
		egin.Response(ctx, result, code, err)
	})

	r.Handle("GET", "/role/resource", func(ctx *gin.Context) {
		result, code, err := ctrl.Resource(ctx)
		egin.Response(ctx, result, code, err)
	})

}
