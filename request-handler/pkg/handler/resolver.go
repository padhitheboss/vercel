package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	queuemodel "github.com/padhitheboss/request-handler/pkg/queueHelper/model"
	"github.com/padhitheboss/request-handler/pkg/utilities/azblob"
)

var b azblob.BlobConfig

func ConnectBlob() {
	b = azblob.CreateConfig()
}
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	id := strings.Split(host, ".")[0]
	id, err := queuemodel.GetFromDB(id)
	log.Println(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "index.html"
	}
	fmt.Println(id, r.URL.Path)
	downloadResponse, err := b.DownloadStream(path.Join(id, filePath))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	contentType := getContentType(filePath)
	w.Header().Set("Content-Type", contentType)
	// Serve Content
	io.Copy(w, downloadResponse)
}

func getContentType(filePath string) string {
	switch {
	case strings.HasSuffix(filePath, ".html"):
		return "text/html"
	case strings.HasSuffix(filePath, ".css"):
		return "text/css"
	case strings.HasSuffix(filePath, ".js"):
		return "application/javascript"
	default:
		return "application/octet-stream"
	}
}
