package handler

import (
	"hash/fnv"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hmdp-backend/internal/dto"
)

// UploadHandler mirrors UploadController.java.
type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/upload")
	group.POST("/blog", h.uploadImage)
	group.GET("/blog/delete", h.deleteBlogImage)
}

func (h *UploadHandler) uploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail("missing file"))
		return
	}
	fileName := h.createNewFileName(file.Filename)
	target := filepath.Join(h.uploadDir, strings.TrimPrefix(fileName, "/"))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail("failed to create dir"))
		return
	}
	if err := ctx.SaveUploadedFile(file, target); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail("文件上传失败"))
		return
	}
	ctx.JSON(http.StatusOK, dto.OkWithData(fileName))
}

func (h *UploadHandler) deleteBlogImage(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, dto.Fail("invalid filename"))
		return
	}
	target := filepath.Join(h.uploadDir, strings.TrimPrefix(name, "/"))
	info, err := os.Stat(target)
	if err != nil {
		ctx.JSON(http.StatusOK, dto.Ok())
		return
	}
	if info.IsDir() {
		ctx.JSON(http.StatusBadRequest, dto.Fail("错误的文件名称"))
		return
	}
	if err := os.Remove(target); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Fail("删除失败"))
		return
	}
	ctx.JSON(http.StatusOK, dto.Ok())
}

func (h *UploadHandler) createNewFileName(original string) string {
	suffix := ""
	if idx := strings.LastIndex(original, "."); idx >= 0 {
		suffix = original[idx+1:]
	}
	name := uuid.NewString()
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(name))
	hash := hasher.Sum32()
	d1 := int(hash & 0xF)
	d2 := int((hash >> 4) & 0xF)
	rel := filepath.ToSlash(filepath.Join("blogs", strconv.Itoa(d1), strconv.Itoa(d2), name))
	if suffix != "" {
		rel = rel + "." + suffix
	}
	return "/" + rel
}
