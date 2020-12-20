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

func RegCommonConfigRouter(r *gin.Engine) {
	ctrl := controller.CommonConfig{}

	r.Handle("GET", "/common_config/list", func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			var params controller.CommonConfigFilter
			errs := utils.Validated(ctx, &params)
			if errs != nil {
				egin.Fail(ctx, consts.ErrorParam, strings.Join(errs, "\n"))
				return
			}
			result, code, err := ctrl.Get(ctx, params)
			egin.Response(ctx, result, code, err)
		}
	}(), middleware.IpLimiter())

	r.Handle("POST", "/common_config/update", func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			var params controller.CommonConfigForm
			errs := utils.Validated(ctx, &params)
			if errs != nil {
				egin.Fail(ctx, consts.ErrorParam, strings.Join(errs, "\n"))
				return
			}
			result, code, err := ctrl.Post(ctx, params)
			egin.Response(ctx, result, code, err)
		}
	}())

	r.Handle("PUT", "/common_config/update/:id", func() func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			var params controller.CommonConfigForm
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

	r.Handle("DELETE", "/common_config/delete/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		result, code, err := ctrl.Delete(ctx, id)
		egin.Response(ctx, result, code, err)
	})

}
