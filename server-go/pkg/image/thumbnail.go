package image

import (
	"bytes"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

const (
	MaxOriginalWidth = 1920
	ThumbSize        = 800
	JPEGQuality      = 95
	ThumbQuality     = 90
)

type ProcessResult struct {
	Original    []byte
	Thumbnail   []byte
	OrigExt     string
	ContentType string
}

func Process(fileBytes []byte, ext string) (*ProcessResult, error) {
	isGIF := ext == ".gif"

	if isGIF {
		return processGIF(fileBytes, ext)
	}

	isJPEG := ext == ".jpg" || ext == ".jpeg"

	cfg, _, err := image.DecodeConfig(bytes.NewReader(fileBytes))
	if err != nil {
		return fallback(fileBytes, ext), nil
	}

	if isJPEG && cfg.Width <= MaxOriginalWidth {
		thumb, err := generateThumbnail(fileBytes)
		if err != nil {
			return &ProcessResult{
				Original:    fileBytes,
				Thumbnail:   nil,
				OrigExt:     ".jpg",
				ContentType: "image/jpeg",
			}, nil
		}
		return &ProcessResult{
			Original:    fileBytes,
			Thumbnail:   thumb,
			OrigExt:     ".jpg",
			ContentType: "image/jpeg",
		}, nil
	}

	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return fallback(fileBytes, ext), nil
	}

	if img.Bounds().Dx() > MaxOriginalWidth {
		img = imaging.Resize(img, MaxOriginalWidth, 0, imaging.Lanczos)
	}

	img = onWhiteBg(img)

	var origBuf bytes.Buffer
	if err := imaging.Encode(&origBuf, img, imaging.JPEG, imaging.JPEGQuality(JPEGQuality)); err != nil {
		return fallback(fileBytes, ext), nil
	}

	thumb := imaging.Fill(img, ThumbSize, ThumbSize, imaging.Center, imaging.Lanczos)
	var thumbBuf bytes.Buffer
	if err := imaging.Encode(&thumbBuf, thumb, imaging.JPEG, imaging.JPEGQuality(ThumbQuality)); err != nil {
		return &ProcessResult{
			Original:    origBuf.Bytes(),
			Thumbnail:   nil,
			OrigExt:     ".jpg",
			ContentType: "image/jpeg",
		}, nil
	}

	return &ProcessResult{
		Original:    origBuf.Bytes(),
		Thumbnail:   thumbBuf.Bytes(),
		OrigExt:     ".jpg",
		ContentType: "image/jpeg",
	}, nil
}

func generateThumbnail(fileBytes []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	img = onWhiteBg(img)
	thumb := imaging.Fill(img, ThumbSize, ThumbSize, imaging.Center, imaging.Lanczos)
	var thumbBuf bytes.Buffer
	if err := imaging.Encode(&thumbBuf, thumb, imaging.JPEG, imaging.JPEGQuality(ThumbQuality)); err != nil {
		return nil, err
	}
	return thumbBuf.Bytes(), nil
}

func processGIF(fileBytes []byte, ext string) (*ProcessResult, error) {
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return &ProcessResult{
			Original:    fileBytes,
			Thumbnail:   nil,
			OrigExt:     ext,
			ContentType: "image/gif",
		}, nil
	}

	img = onWhiteBg(img)
	thumb := imaging.Fill(img, ThumbSize, ThumbSize, imaging.Center, imaging.Lanczos)
	var thumbBuf bytes.Buffer
	if err := imaging.Encode(&thumbBuf, thumb, imaging.JPEG, imaging.JPEGQuality(ThumbQuality)); err != nil {
		return &ProcessResult{
			Original:    fileBytes,
			Thumbnail:   nil,
			OrigExt:     ext,
			ContentType: "image/gif",
		}, nil
	}

	return &ProcessResult{
		Original:    fileBytes,
		Thumbnail:   thumbBuf.Bytes(),
		OrigExt:     ext,
		ContentType: "image/gif",
	}, nil
}

func fallback(fileBytes []byte, ext string) *ProcessResult {
	ct := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		ct = "image/jpeg"
	case ".png":
		ct = "image/png"
	case ".webp":
		ct = "image/webp"
	case ".gif":
		ct = "image/gif"
	}
	return &ProcessResult{
		Original:    fileBytes,
		Thumbnail:   nil,
		OrigExt:     ext,
		ContentType: ct,
	}
}

func onWhiteBg(img image.Image) *image.NRGBA {
	b := img.Bounds()
	bg := imaging.New(b.Dx(), b.Dy(), color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	return imaging.Paste(bg, img, image.Pt(0, 0))
}
