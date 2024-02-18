package gitcmd

import (
	"os"

	"example.com/uploader/pkg/model"
)

type GitService struct {
	repoUrl  string
	cloneDir string
	username string
	password string
	sshKey   string
	status   string
}

func (g *GitService) GetStatus() string {
	return g.status
}
func CreateGitService(req model.Request) GitService {
	var g GitService
	g.cloneDir = os.Getenv("LOCAL_FOLDER_PATH") + "/" + req.RequestId + "/"
	g.repoUrl = req.RepoUrl
	g.username = req.Username
	g.password = req.Password
	g.sshKey = req.Ssh_key
	return g
}
