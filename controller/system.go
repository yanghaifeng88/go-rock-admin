package controller

import (
	"github.com/daodao97/egin/consts"
	"github.com/daodao97/egin/utils/upload/filesystem"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// @Controller 系统接口
type System struct {
}

// @GetApi /system/config
func (s System) Config(c *gin.Context) (interface{}, consts.ErrCode, error) {
	config := map[string]interface{}{
		"logo":  "https://f.mengtuiapp.com/upload/oss/1/201912/7589f22c5cef9e7c-158x158.png",
		"title": "萌推小二",
	}

	return config, 0, nil
}

// @PostApi /upload
func (s System) Upload(c *gin.Context) (interface{}, consts.ErrCode, error) {
	file, err := c.FormFile("file")
	if file == nil {
		return nil, 500, errors.New("请选择上传文件")
	}
	localFile := "/tmp/" + file.Filename // 本地文件路径
	// gin 简单做了封装,拷贝了文件流
	if err := c.SaveUploadedFile(file, localFile); err != nil {
		// ignore
		return nil, 0, errors.New("本地保存失败")
	}

	oss, err := filesystem.NewDefault()
	if err != nil {
		return nil, 0, err
	}
	url, err := oss.Save("qupinapptest", localFile, "9527/"+file.Filename)

	return map[string]string{"url": url}, 0, err
}
