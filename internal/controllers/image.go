package controllers

import (
	"app/internal/config"
	"app/internal/image"
	"app/internal/utils"
	"app/pkg/filesystem"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AllowedExtensions struct {
	Extensions map[string]string
}

type QueryParams struct {
	Width  *int    `json:"width,omitempty"`
	Height *int    `json:"height,omitempty"`
	Op     *string `json:"op,omitempty"`
}

var (
	cfg      *config.Config
	fileExts AllowedExtensions
)

func init() {
	cfg = config.GetConfig()

	fileExts.Extensions = make(map[string]string)
	for _, ext := range cfg.AllovedImageExtensions {
		fileExts.Extensions[ext] = "1"
	}
}

func ImageController(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	fileDirectory := filepath.Dir(r.URL.Path)
	relativePath := filepath.Join(fileDirectory, filename)

	if !isAllowedExtension(filename) {
		render404(w)

		return
	}

	sourceUrl := cfg.Upstream + relativePath
	cachePath := prepareCachePath(r.URL)
	absolutePath := filepath.Join(cfg.Storage.Path, cachePath)

	if filesystem.IsFile(absolutePath) {
		renderImage(w, absolutePath)

		return
	}

	qParams := populateQueryParams(r)

	err := utils.DownloadFile(sourceUrl, absolutePath)

	if err != nil {
		return
	}

	fileInfo, _ := os.Stat(absolutePath)

	modifiedTime := fileInfo.ModTime() // запоминаем время создания оригинала, чтобы потом сверять и инвалидировать этот кеш

	op := ""

	if qParams.Op != nil {
		op = *qParams.Op
	}

	err = nil
	// Удаляем изображение из кеша, если по какой-то причине обработка завершилась неудачно.
	switch op {
	case "t16x9":
		err = image.Thumbnail16x9(absolutePath, qParams.Width, qParams.Height, absolutePath)
	case "c2f":
		err = image.CropToFit(absolutePath, qParams.Width, qParams.Height, absolutePath)
	default:
		if qParams.Width != nil {
			err = image.ResizeToWidth(absolutePath, *qParams.Width, qParams.Height, absolutePath)
		}
	}

	if err != nil {
		os.Remove(absolutePath)
	}

	if filesystem.IsFile(absolutePath) {
		image.Optimize(absolutePath, nil)
		os.Chtimes(absolutePath, modifiedTime, modifiedTime)

		renderImage(w, absolutePath)

		return
	}

	render502(w)
}

func render404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)

	fmt.Fprint(w, "file not found")
}

func render502(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)

	fmt.Fprint(w, "internal server error")
}

/**
 * Рендер изображения
 */
func renderImage(w http.ResponseWriter, sourceFilepath string) {
	file, err := os.Open(sourceFilepath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	defer file.Close()

	fileInfo, err := os.Stat(sourceFilepath)

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	lastModified := fileInfo.ModTime()
	w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))

	fileExtension := filepath.Ext(sourceFilepath)

	contentType := mime.TypeByExtension(fileExtension)

	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

/**
 * Проверяет является ли расширение запрашиваемого файла поддерживаемым.
 */
func isAllowedExtension(filename string) bool {

	fileExtension := filepath.Ext(filename)[1:]

	return fileExts.Extensions[fileExtension] != ""
}

/**
 * Нормализует путь для кеш файла.
 */
func prepareCachePath(requestUrl *url.URL) string {
	path := strings.NewReplacer(
		"./", "",
		"../", "",
		"..", "",
		"?", ".",
		"&", "_",
		"=", "_",
		"#", "_",
		":", "_",
		"*", "_",
		"\"", "_",
		"#", "_",
		"<", "_",
		">", "_",
		"|", "_",
	).Replace(requestUrl.Path)

	querySting := strings.NewReplacer(
		".", "_",
		"/", "_",
		"?", "_",
		"&", "_",
		"=", "_",
		"#", "_",
		":", "_",
		"*", "_",
		"\"", "_",
		"#", "_",
		"<", "_",
		">", "_",
		"|", "_",
	).Replace(requestUrl.RawQuery)

	if len(querySting) > 200 {
		querySting = querySting[:200]
	}

	filename := filepath.Base(path)

	if querySting != "" {
		filename = querySting + "." + filename
	}

	fileDirectory := filepath.Dir(path)

	return filepath.Join(fileDirectory, filename)
}

func populateQueryParams(r *http.Request) QueryParams {
	// Создаем экземпляр структуры Dimensions
	qParams := QueryParams{}

	// Получаем параметры из строки запроса
	widthStr := r.URL.Query().Get("width")
	heightStr := r.URL.Query().Get("height")
	opStr := r.URL.Query().Get("op")

	// Обрабатываем параметр width
	if widthStr != "" {
		width, err := strconv.Atoi(widthStr)

		if err == nil {
			qParams.Width = &width
		}
	}

	// Обрабатываем параметр height
	if heightStr != "" {
		height, err := strconv.Atoi(heightStr)

		if err == nil {
			qParams.Height = &height
		}
	}

	// Обрабатываем параметр op
	if opStr != "" {
		op := strings.Trim(opStr, " ")
		qParams.Op = &op
	}

	return qParams
}
