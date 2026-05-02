package upload

import (
	"net/http"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadModule module.UploadModule
}

func NewUploadHandler(uModule module.UploadModule) *UploadHandler {
	return &UploadHandler{uploadModule: uModule}
}

// UploadImages godoc
// @Summary Upload multiple images
// @Description Upload multiple images to Cloudinary and get public URLs
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Images to upload"
// @Success 200 {object} response.BaseResponse{data=[]string}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /upload/images [post]
func (h *UploadHandler) UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, "Failed to parse form", http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		appErr := errorx.New(errorx.ErrBadRequest, "No files uploaded", http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	var interfaces []interface{}
	var fileNames []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			appErr := errorx.New(errorx.ErrInternal, "Failed to open file", http.StatusInternalServerError)
			response.CustomError(c, appErr)
			return
		}
		defer file.Close()
		interfaces = append(interfaces, file)
		fileNames = append(fileNames, fileHeader.Filename)
	}

	urls, err := h.uploadModule.UploadImages(c.Request.Context(), interfaces, fileNames)
	if err != nil {
		appErr := errorx.New(errorx.ErrInternal, err.Error(), http.StatusInternalServerError)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, urls)
}
