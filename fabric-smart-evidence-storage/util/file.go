package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)


var ipfsClient *IPFSClient

func InitIPFSClient(ipfsNodeURL string) {
	ipfsClient = NewIPFSClient(ipfsNodeURL)
}

func CalculateFileHash(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hashAlgorithm := sha256.New()
	// 将文件内容写入哈希计算器
	if _, err := io.Copy(hashAlgorithm, file); err != nil {
		return "", err
	}

	// 计算哈希值并转换为十六进制字符串
	hashInBytes := hashAlgorithm.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

// IPFSClient 封装了 IPFS 的上传和下载功能
type IPFSClient struct {
	sh *shell.Shell
}

func NewIPFSClient(ipfsNodeURL string) *IPFSClient {
	return &IPFSClient{
		sh: shell.NewShell(ipfsNodeURL),
	}
}

// UploadFile 上传文件到 IPFS，返回文件的 CID
func UploadFile(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 上传文件到 IPFS
	cid, err := ipfsClient.sh.Add(file)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %v", err)
	}

	return cid, nil
}

// DownloadFile 从 IPFS 下载文件
func DownloadFile(cid string) ([]byte, error) {
	// 从 IPFS 获取文件内容
	reader, err := ipfsClient.sh.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %v", err)
	}
	defer reader.Close()

	// 读取文件内容
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	return content, nil
}
