package queuemodel

import (
	"encoding/json"

	queueconfig "github.com/padhitheboss/code-builder/pkg/queueHelper/config"
)

var q queueconfig.QueueFunction

func StartQueue() {
	q = queueconfig.InitializeQueue("redis")
}

type GitAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Ssh_key  string `json:"ssh_key,omitempty"`
}
type Request struct {
	RequestId    string `json:"request_id"`
	RepoUrl      string `json:"repo_url"`
	EndPointName string `json:"endpoint_name"`
	GitAuth
}

type Response struct {
	RequestId   string `json:"request_id"`
	Status      string `json:"status"`
	BlobUrl     string `json:"blob_url"`
	ResourceId  string `json:"resource_id"`
	ProjectType string `json:"project_type"`
}

func (auth *GitAuth) GetGitUserPass() (string, string) {
	return auth.Username, auth.Password
}

func (auth *GitAuth) GetGitSSHKey() (string, string) {
	return auth.Username, auth.Ssh_key
}
func (req *Request) GetRepoUrl() string {
	return req.RepoUrl
}

func CreateResponse(id string, status string, blobResponse interface{}) string {
	res := Response{
		RequestId:  id,
		Status:     status,
		BlobUrl:    "",
		ResourceId: "",
	}
	responseString, _ := json.Marshal(res)
	return string(responseString)
}

func SendResponse(message string) error {
	return q.WriteToQueue(message)
}

func CollectRequest() (Response, error) {
	message, err := q.ReadFromQueue()
	var req Response
	json.Unmarshal([]byte(message), &req)
	return req, err
}