package main

import (
	"fmt"
	"log"
	"os"

	"github.com/padhitheboss/code-builder/pkg/logger"
	queuemodel "github.com/padhitheboss/code-builder/pkg/queueHelper/model"
	"github.com/padhitheboss/code-builder/pkg/service/azureBlob"
	pkgbuilder "github.com/padhitheboss/code-builder/pkg/service/buildStatic"
	ziphandler "github.com/padhitheboss/code-builder/pkg/service/zipHandler"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println(err)
	// }
	queuemodel.StartQueue()
	sessionLog, err := queuemodel.GetFromDB(os.Getenv("ENDPOINT_NAME"))
	if err != nil {
		fmt.Println("internal logger error")
	}
	layer := logger.CreateLayerLog()
	sessionLog.LayerLogs = append(sessionLog.LayerLogs, layer)
	layerLog := &sessionLog.LayerLogs[len(sessionLog.LayerLogs)-1]
	b := azureBlob.CreateConfig("", os.Getenv("REQUEST_ID"))
	downloadPath, err := b.DownloadFile(os.Getenv("BLOB_INPUT_FOLDER_PATH"), os.Getenv("DOWNLOAD_FOLDER_PATH"))
	if err != nil {
		log.Printf("unable to download repo from blob %v", err)
		layerLog.Status = "failed"
		layerLog.Rawlogs += fmt.Sprintf("unable to download repo from blob %v", err)
		queuemodel.UpdateDB(sessionLog)
		return
	}
	layerLog.Rawlogs += fmt.Sprintf("downloaded repo from blob %v", err)
	fmt.Println(downloadPath)
	unzipFilePath := ziphandler.UnzipFile(downloadPath, os.Getenv("BUILD_FOLDER_PATH"))
	layerLog.Rawlogs += fmt.Sprintf("\nunzipped repo files %v and starting build", unzipFilePath)
	pkg := pkgbuilder.CreateConfig(os.Getenv("PROJECT_TYPE"), unzipFilePath)
	consoleOutput, err := pkg.BuildStatic()
	layerLog.Rawlogs += consoleOutput
	if err != nil {
		log.Printf("unable to build the project %v", err)
		layerLog.Status = "failed"
		layerLog.Rawlogs += fmt.Sprintf("unable to build the project %v", err)
		queuemodel.UpdateDB(sessionLog)
		return
	}
	fmt.Println(pkg.GetOutputPath())
	b.UploadFolder(pkg.GetOutputPath())
	layerLog.Rawlogs += fmt.Sprintf("\nuploaded files to folder %v", pkg.GetOutputPath())
	queuemodel.UpdateDB(sessionLog)
	if err != nil {
		log.Println(err)
	}
}
