package cloudinary

import (
	"context"
	"fmt"

	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryService struct {
	cl *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudName, apiKey, apiSecret string) (platform.FileUploader, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &cloudinaryService{cl: cld}, nil
}

func (s *cloudinaryService) UploadFile(ctx context.Context, file interface{}, fileName string) (string, error) {
	return s.uploadFileToFolder(ctx, file, fileName, "gabaa_products")
}

func (s *cloudinaryService) uploadFileToFolder(ctx context.Context, file interface{}, fileName, folder string) (string, error) {
	uploadResult, err := s.cl.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       fileName,
		Folder:         folder,
		UniqueFilename: api.Bool(true),
		Overwrite:      api.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func (s *cloudinaryService) UploadMultiple(ctx context.Context, files []interface{}, fileNames []string) ([]string, error) {
	return s.UploadMultipleToFolder(ctx, files, fileNames, "gabaa_products")
}

func (s *cloudinaryService) UploadMultipleToFolder(ctx context.Context, files []interface{}, fileNames []string, folder string) ([]string, error) {
	if len(files) != len(fileNames) {
		return nil, fmt.Errorf("files and fileNames length mismatch")
	}

	urls := make([]string, len(files))
	for i, file := range files {
		url, err := s.uploadFileToFolder(ctx, file, fileNames[i], folder)
		if err != nil {
			return nil, err
		}
		urls[i] = url
	}

	return urls, nil
}
