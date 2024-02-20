package model

import (
	"github.com/gofrs/uuid"
)

type GitAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Ssh_key  string `json:"ssh_key,omitempty"`
}
type Request struct {
	RequestId    string `json:"request_id"`
	RepoUrl      string `json:"repo_url"`
	EndPointName string `json:"endpoint_name"`
	ProjectType  string `json:"project_type"`
	GitAuth
}

func (r *Request) GenerateId() {
	newUUID, _ := uuid.NewV4()
	r.RequestId = newUUID.String()
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
