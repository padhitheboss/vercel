package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"example.com/uploader/pkg/model"
	queuemodel "example.com/uploader/pkg/queueHelper/model"
	"example.com/uploader/pkg/service/azureBlob"
	"example.com/uploader/pkg/service/gitcmd"
	ziphandler "example.com/uploader/pkg/service/zipHandler"
)

func RequestDeploy(w http.ResponseWriter, r *http.Request) {
	var req model.Request
	json.NewDecoder(r.Body).Decode(&req)
	req.GenerateId()
	fmt.Println(req)
	go func() {
		logger := model.StartLog(req.RequestId, req.EndPointName)
		layer := model.CreateLayerLog()
		logger.LayerLogs = append(logger.LayerLogs, layer)
		layerLog := &logger.LayerLogs[len(logger.LayerLogs)-1]
		g := gitcmd.CreateGitService(req)
		outputPath, err := g.Clone()
		if err != nil {
			log.Printf("unable to clone git repo %v", err)
			layerLog.Rawlogs += fmt.Sprintf("\nunable to clone git repo %v", err)
			layerLog.Status = "failed"
			queuemodel.UpdateDB(logger)
			return
		}
		layerLog.Rawlogs += fmt.Sprintf("\ncloned git repo %v", req.RepoUrl)
		fmt.Println(outputPath)
		outputPath, err = ziphandler.CreateZipFile(fmt.Sprintf("/tmp/srccode/%v.zip", req.RequestId), outputPath)
		if err != nil {
			log.Printf("unable to zip the files %v", err)
			layerLog.Rawlogs += fmt.Sprintf("/nunable to zip the files %v", err)
			layerLog.Status = "failed"
			queuemodel.UpdateDB(logger)
			return
		}
		layerLog.Rawlogs += fmt.Sprintf("/ncreated zip file of the repository %v", req.RepoUrl)
		b := azureBlob.CreateConfig("")
		uploadPath, err := b.UploadFile(outputPath, path.Base(outputPath))
		if err != nil {
			log.Printf("unable to upload the files %v", err)
			layerLog.Rawlogs += fmt.Sprintf("/nunable to zip the files %v", err)
			layerLog.Status = "failed"
			queuemodel.UpdateDB(logger)
			return
		}
		layerLog.Rawlogs += fmt.Sprintf("/nuploaded zip file to blob %v", err)
		layerLog.Status = "success"
		message := queuemodel.CreateResponse(req, uploadPath, "successful",)
		queuemodel.SendResponse(message)
		queuemodel.UpdateDB(logger)
	}()
	json.NewEncoder(w).Encode(req)
}

func DeployStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	logger, err := queuemodel.GetFromDB(id)
	if err != nil {
		w.Write([]byte("unable to get details at this moment"))
	}
	json.NewEncoder(w).Encode(logger)
}
