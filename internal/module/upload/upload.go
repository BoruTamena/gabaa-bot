package upload

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/platform"
)

type uploadModule struct {
	uploader platform.FileUploader
}

func NewUploadModule(uploader platform.FileUploader) module.UploadModule {
	return &uploadModule{uploader: uploader}
}

func (m *uploadModule) UploadImages(ctx context.Context, files []interface{}, fileNames []string) ([]string, error) {
	return m.uploader.UploadMultiple(ctx, files, fileNames)
}
