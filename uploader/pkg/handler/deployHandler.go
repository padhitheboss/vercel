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
		g := gitcmd.CreateGitService(req)
		outputPath, err := g.Clone()
		if err != nil {
			log.Panicf("unable to clone git repo %v", err)
			return
		}
		fmt.Println(outputPath)
		outputPath, err = ziphandler.CreateZipFile(fmt.Sprintf("/tmp/srccode/%v.zip", req.RequestId), outputPath)
		if err != nil {
			log.Panicf("unable to zip the files %v", err)
			return
		}
		b := azureBlob.CreateConfig("")
		uploadPath, err := b.UploadFile(outputPath, path.Base(outputPath))
		if err != nil {
			log.Panicf("unable to upload the files %v", err)
			return
		}
		message := queuemodel.CreateResponse(req, uploadPath, "successful")
		queuemodel.SendResponse(message)
	}()
	// message := model.CreateResponse(req)
	// fmt.Println(message)
	// queuemodel.SendResponse(message)
	json.NewEncoder(w).Encode(req)
}

func DeployStatus(w http.ResponseWriter, r *http.Request) {
	// TO BE IMPLEMENTED
}
