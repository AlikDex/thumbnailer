package controllers

import (
    "app/internal/config"
    "app/internal/image"
    "app/pkg/filesystem"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/gorilla/mux"
)

var cfg *config.Config

func init() {
    cfg = config.LoadConfig()
}

func ThumbController(w http.ResponseWriter, r *http.Request) {
    routeVars := mux.Vars(r)
    queryParams := r.URL.Query()

    path := routeVars["path"]

    sourceImagePath := cfg.Storage.Path + "/" + path

    if !filesystem.IsFile(sourceImagePath) || !isSupportedExtension(sourceImagePath) {
        render404(w)

        return
    }

    width := 320
    requestWidth := queryParams.Get("w")
    if requestWidth != "" {
        parsedWidth, err := strconv.Atoi(requestWidth)

        if err == nil && parsedWidth <= 2048 {
            width = parsedWidth
        }
    }

    height := 180
    requestHeight := queryParams.Get("h")
    if requestHeight != "" {
        parsedHeight, err := strconv.Atoi(requestHeight)

        if err == nil && parsedHeight <= 1024 {
            height = parsedHeight
        }
    }

    quality := 85
    requestQuality := queryParams.Get("q")
    if requestQuality != "" {
        parsedQuality, err := strconv.Atoi(requestQuality)

        if err == nil && parsedQuality >= 60 && parsedQuality <= 100 {
            quality = parsedQuality
        }
    }

    op := "none" // operation
    requestOp := queryParams.Get("op")

    if requestOp != "" {
        op = requestOp
    }

    renderImagePath := sourceImagePath

    if isSupportedOp(op) {
        // Создадим временный файл
        tmpFile, err := os.CreateTemp(cfg.TmpDir, "thumb-*.jpg")

        if err != nil {
            log.Fatal(err)
        }

        defer os.Remove(tmpFile.Name()) // Удаляем файл после использования
        defer tmpFile.Close()

        renderImagePath = tmpFile.Name()

        switch op {
            case "r2w":
                image.ResizeToWidth(sourceImagePath, width, height, renderImagePath)
            case "c2f":
                image.CropToFit(sourceImagePath, width, height, renderImagePath)
        }

        image.Optimize(renderImagePath, quality)
    }

    renderImage(w, renderImagePath)

    /*file, err := os.Open(tmpFile.Name())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer file.Close()*/

    /*buffer := make([]byte, 512)
    _, err = file.Read(buffer)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    contentType := http.DetectContentType(buffer)

    _, err = file.Seek(0, 0)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }*/

    /*w.Header().Set("Content-Type", "image/jpeg")
    _, err = io.Copy(w, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }*/
}

func render404(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNotFound)

    fmt.Fprint(w, "file not found")
}

func renderImage(w http.ResponseWriter, sourceFilepath string) {
    file, err := os.Open(sourceFilepath)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)

        return
    }

    defer file.Close()

    w.Header().Set("Content-Type", "image/jpeg")
    _, err = io.Copy(w, file)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)

        return
    }
}

func isSupportedOp(operation string) bool {
    supportedOp := []string{"r2w", "c2f"}

    for _, op := range supportedOp {
        if op == operation {
            return true
        }
    }

    return false
}

func isSupportedExtension(filename string) bool {
    ext := filepath.Ext(filename)
    supportedExtensions := []string{".jpeg", ".jpg", ".png", ".webp"}

    for _, supportedExt := range supportedExtensions {
        if ext == supportedExt {
            return true
        }
    }

    return false
}
