package image

import (
	"app/pkg/filesystem"
	"fmt"
	goImg "image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

/**
 * ResizeToWidth изменяет размер изображения по ширине, сохраняя пропорции.
 * Если задан maxHeight, ограничивает высоту изображения, обрезая края сверху и снизу.
 */
func ResizeToWidth(srcFilepath string, width int, maxHeight *int, outFilepath string) error {
	if !filesystem.IsFile(srcFilepath) {
		return fmt.Errorf("source path is not a file or not exists: %s", srcFilepath)
	}

	outDir := filepath.Dir(outFilepath)

	if !filesystem.IsDir(outDir) {
		return fmt.Errorf("out directory not exists: %s", outDir)
	}

	w := strconv.Itoa(width)
	var cmd *exec.Cmd

	if maxHeight == nil {
		cmd = exec.Command("convert", srcFilepath, "-thumbnail", w+"x", "-strip", "-quality", "100", "+repage", outFilepath)
	} else {
		mh := strconv.Itoa(*maxHeight)

		cmd = exec.Command("convert", srcFilepath, "-thumbnail", w+"x", "-gravity", "center", "-crop", w+"x"+mh+"+0+0", "-strip", "-quality", "100", "+repage", outFilepath)
	}

	return cmd.Run()
}

func CropToFit(srcFilepath string, width *int, height *int, outFilepath string) error {
	if !filesystem.IsFile(srcFilepath) {
		return fmt.Errorf("source path is not a file or not exists: %s", srcFilepath)
	}

	outDir := filepath.Dir(outFilepath)

	if !filesystem.IsDir(outDir) {
		return fmt.Errorf("out directory not exists: %s", outDir)
	}

	if width == nil || height == nil {
		return fmt.Errorf("width and height should be specified for this operation")
	}

	w := strconv.Itoa(*width)
	h := strconv.Itoa(*height)

	// -resize 300x300^ -gravity center -crop 300x300+0+0 +repage
	// -thumbnail 320x180^ -gravity center -crop 320x180+0+0 -unsharp 0x.5 -strip +repage
	cmd := exec.Command("convert", srcFilepath, "-thumbnail", w+"x"+h+"^", "-gravity", "center", "-crop", w+"x"+h+"+0+0", "-strip", "-quality", "100", "+repage", outFilepath)

	return cmd.Run()
}

func Thumbnail16x9(srcFilepath string, width *int, height *int, outFilepath string) error {
	if !filesystem.IsFile(srcFilepath) {
		return fmt.Errorf("source path is not a file or not exists: %s", srcFilepath)
	}

	outDir := filepath.Dir(outFilepath)

	if !filesystem.IsDir(outDir) {
		return fmt.Errorf("out directory not exists: %s", outDir)
	}

	if width == nil && height == nil {
		return fmt.Errorf("please specify one of the sides")
	}

	w := ""
	h := ""
	var cmd *exec.Cmd
	imgSize, err := getImageSize(srcFilepath)

	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	if (width != nil && *width > 0) && (height != nil && *height > 0) {
		w = strconv.Itoa(*width)
		h = strconv.Itoa(int(float64(*width) * 9.0 / 16.0))
	} else if height == nil {
		w = strconv.Itoa(*width)
		h = strconv.Itoa(int(float64(*width) * 9.0 / 16.0))
	} else if width == nil {
		h = strconv.Itoa(*height)
		w = strconv.Itoa(int(float64(*height) * 16.0 / 9.0))
	}

	boxSize := w + "x" + h

	if isWiderThan16by9(imgSize.Width, imgSize.Height) {
		cmd = exec.Command("convert", srcFilepath, "-resize", boxSize+"^", "-background", "black", "-gravity", "center", "-extent", boxSize, "-strip", "-quality", "100", "-colorspace", "sRGB", outFilepath)
	} else {
		cmd = exec.Command("convert", srcFilepath, "-resize", boxSize, "-background", "black", "-gravity", "center", "-extent", boxSize, "-strip", "-quality", "100", "-colorspace", "sRGB", outFilepath)
	}

	return cmd.Run()
}

/**
 * Проверяет, является ли прямоугольник соотношением 16:9 или шире
 */
func isWiderThan16by9(width int, height int) bool {
	if height == 0 {
		return false
	}

	isHorizontal := width > height
	aspectRatio16x9 := 16.0 / 9.0
	currentAspectRatio := float64(width) / float64(height)

	return isHorizontal && currentAspectRatio >= aspectRatio16x9
}

type ImageDimension struct {
	Width  int
	Height int
}

func getImageSize(srcPath string) (ImageDimension, error) {
	dimension := ImageDimension{}

	file, err := os.Open(srcPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка при открытии файла: %v\n", err)

		return dimension, err
	}

	defer file.Close() // Закрываем файл после завершения работы

	// Получаем размеры изображения
	config, _, err := goImg.DecodeConfig(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка при декодировании изображения: %v\n", err)

		return dimension, err
	}

	dimension.Width = config.Width
	dimension.Height = config.Height

	return dimension, nil
}
