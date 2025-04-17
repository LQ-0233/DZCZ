package api

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"fabric-smart-evidence-storage/util"

	"github.com/gin-gonic/gin"
)

type ServerError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func serverError(c *gin.Context, err error) {
	c.JSON(500, ServerError{
		Code: 500,
		Msg:  err.Error(),
	})
}

func UploadFile(c *gin.Context) {
	// 从请求中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		serverError(c, err)
		return
	}
	timestamp := time.Now().Unix()
	// 保存文件到指定路径
	filePath := "./uploads/" + file.Filename + strconv.FormatInt(timestamp, 10)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		serverError(c, err)
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{
		"fileName": file.Filename,
		"filePath": filePath,
	})
}

func DownloadFile(c *gin.Context) {
	// 从请求中获取文件路径
	ipfsLink := c.Query("ipfsLink")
	fileName := c.Query("fileName")
	if ipfsLink == "" {
		serverError(c, errors.New("ipfsLink不能为空"))
		return
	}

	// 检查文件名参数
	if fileName == "" {
		serverError(c, errors.New("fileName不能为空"))
		return
	}

	content, err := util.DownloadFile(ipfsLink)
	if err != nil {
		serverError(c, err)
		return
	}

	// 设置下载文件的 header utf-8
	encodedFileName := url.QueryEscape(fileName)
	c.Header("Content-Disposition", "attachment; filename="+encodedFileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(200, "application/octet-stream", content)
}
