package image

import (
	"os/exec"
	"strconv"
)

func Optimize(filename string, quality *int) error {
	q := 90 // default quality

	if quality != nil && (*quality >= 60 && *quality <= 100) {
		q = *quality
	}

	qualityParameter := strconv.Itoa(q)

	//jpegoptim -m <quality> <filepath>
	cmd := exec.Command("jpegoptim", "-m", qualityParameter, filename)

	return cmd.Run()
}
