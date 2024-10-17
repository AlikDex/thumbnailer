package image

import (
	"app/pkg/filesystem"
	"errors"
	"os/exec"
	"path/filepath"
	"strconv"
)

func ResizeToWidth(srcFilepath string, width int, maxHeight int, outFilepath string) error {
	if !filesystem.IsFile(srcFilepath) {
		return errors.New("Source path is not a file or not exists: " + srcFilepath)
	}

	outDir := filepath.Dir(outFilepath)

	if !filesystem.IsDir(outDir) {
		return errors.New("Out directory not exists: " + srcFilepath)
	}

	w := strconv.Itoa(width)
	mh := strconv.Itoa(maxHeight)

	// -resize 300x -gravity center -crop 300x400+0+0 +repage
	// -thumbnail 320x -gravity center -crop 320x400+0+0 -unsharp 0x.5 -strip +repage
	cmd := exec.Command("convert", srcFilepath, "-thumbnail", w+"x", "-gravity", "center", "-crop", w+"x"+mh+"+0+0", "-unsharp", "0x.5", "-strip", "+repage", outFilepath)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func CropToFit(srcFilepath string, width int, height int, outFilepath string) error {
	if !filesystem.IsFile(srcFilepath) {
		return errors.New("Source path is not a file or not exists: " + srcFilepath)
	}

	outDir := filepath.Dir(outFilepath)

	if !filesystem.IsDir(outDir) {
		return errors.New("Out directory not exists: " + srcFilepath)
	}

	w := strconv.Itoa(width)
	h := strconv.Itoa(height)

	// -resize 300x300^ -gravity center -crop 300x300+0+0 +repage
	// -thumbnail 320x180^ -gravity center -crop 320x180+0+0 -unsharp 0x.5 -strip +repage
	cmd := exec.Command("convert", srcFilepath, "-thumbnail", w+"x"+h+"^", "-gravity", "center", "-crop", w+"x"+h+"+0+0", "-unsharp", "-0x.5", "-strip", "+repage", outFilepath)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
