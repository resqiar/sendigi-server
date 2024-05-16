package utils

import (
	"context"
	"log"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/media"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func UploadImages(image string, packageName string) (string, error) {
	ik, err := imagekit.New()
	if err != nil {
		log.Println(err)
		return "", err
	}

	imageExist, err := ik.Media.Files(context.Background(), media.FilesParam{
		Tags: packageName,
	})
	if len(imageExist.Data) > 0 {
		return imageExist.Data[0].Url, nil
	}

	// upload image base64 to ImageKit
	resp, err := ik.Uploader.Upload(context.Background(), image, uploader.UploadParam{
		FileName: packageName,
		Folder:   "sendigi",
		Tags:     packageName,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	return resp.Data.Url, nil
}
