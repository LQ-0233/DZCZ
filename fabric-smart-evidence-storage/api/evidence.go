package api

import (
	"errors"
	"fabric-smart-evidence-storage/fabric"
	"fabric-smart-evidence-storage/model"
	"fabric-smart-evidence-storage/util"
	"math/big"

	"github.com/gin-gonic/gin"
)

// EvidenceCreate 创建证据
func EvidenceCreate(c *gin.Context) {
	var evidence model.Evidence
	if err := c.ShouldBind(&evidence); err == nil {
		username, exists := c.Get("username")
		if !exists {
			serverError(c, errors.New("用户未登录"))
			return
		}
		role, exists := c.Get("role")
		if !exists {
			serverError(c, errors.New("用户未登录"))
			return
		}
		roleBigInt := new(big.Int)
		_, ok := roleBigInt.SetString(role.(string), 10)
		if !ok {
			serverError(c, errors.New("角色错误"))
			return
		}

		ipfsLink, err := util.UploadFile(evidence.FilePath)
		if err != nil {
			serverError(c, err)
			return
		}
		evidence.IpfsLink = ipfsLink
		evidence.FileHash, err = util.CalculateFileHash(evidence.FilePath)
		if err != nil {
			serverError(c, err)
			return
		}
		userCredential, proof, err := util.GenUserCredential(roleBigInt)
		if err != nil {
			serverError(c, err)
			return
		}
		// 调用 fabric 创建证据
		err = fabric.AddEvidence(
			username.(string),
			string(userCredential),
			string(proof),
			evidence.CaseNumber,
			evidence.CaseInfo,
			evidence.Manager1,
			evidence.Manager2,
			evidence.EvidenceType,
			evidence.IpfsLink,
			evidence.FileHash,
			evidence.FileName,
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

// GetUserEvidence 获取用户创建的所有证据
func GetUserEvidence(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		serverError(c, errors.New("用户未登录"))
		return
	}

	evidence, err := fabric.GetUserEvidence(username.(string))
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(200, evidence)
}

// 授权请求结构
type AuthorizeRequest struct {
	CaseNumber     string `json:"caseNumber" binding:"required"`
	AuthorizedUser string `json:"authorizedUser" binding:"required"`
}

// EvidenceAuthorize 授权其他用户查看证据
func EvidenceAuthorize(c *gin.Context) {
	var req AuthorizeRequest
	if err := c.ShouldBind(&req); err == nil {
		username, exists := c.Get("username")
		if !exists {
			serverError(c, errors.New("用户未登录"))
			return
		}

		err = fabric.AddAuthorized(
			username.(string),
			req.CaseNumber,
			req.AuthorizedUser,
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

// 取消授权请求结构
type CancelAuthorizeRequest struct {
	Id string `json:"id" binding:"required"`
}

// EvidenceCancelAuthorize 取消授权
func EvidenceCancelAuthorize(c *gin.Context) {
	var req CancelAuthorizeRequest
	if err := c.ShouldBind(&req); err == nil {
		username, exists := c.Get("username")
		if !exists {
			serverError(c, errors.New("用户未登录"))
			return
		}

		err := fabric.CancelAuthorized(username.(string), req.Id)
		if err != nil {
			serverError(c, err)
			return
		}
		c.Status(200)
	} else {
		serverError(c, err)
	}
}

// EvidenceAuthorizedList 获取用户创建的所有授权记录
func EvidenceAuthorizedList(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		serverError(c, errors.New("用户未登录"))
		return
	}

	authorized, err := fabric.GetUserAllAuthorized(username.(string))
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(200, authorized)
}

// EvidenceReceivedAuthorizedList 获取用户收到的所有授权记录
func EvidenceReceivedAuthorizedList(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		serverError(c, errors.New("用户未登录"))
		return
	}

	authorized, err := fabric.GetUserReceivedAuthorized(username.(string))
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(200, authorized)
}

// 查看证据请求结构
type ViewRequest struct {
	CaseNumber string `json:"caseNumber" binding:"required"`
}

// EvidenceView 记录证据查看行为
func EvidenceView(c *gin.Context) {
	var req ViewRequest
	if err := c.ShouldBind(&req); err == nil {
		username, exists := c.Get("username")
		if !exists {
			serverError(c, errors.New("用户未登录"))
			return
		}
		evidence, err := fabric.View(username.(string), req.CaseNumber)
		if err != nil {
			serverError(c, err)
			return
		}
		c.JSON(200, evidence)
	} else {
		serverError(c, err)
	}
}

// EvidenceViewRecordList 获取证据的查看记录
func EvidenceViewRecordList(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		serverError(c, errors.New("用户未登录"))
		return
	}

	viewRecords, err := fabric.GetUserViewRecord(username.(string))
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(200, viewRecords)
}
