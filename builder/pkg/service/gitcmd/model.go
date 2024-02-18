package git

import queuemodel "github.com/padhitheboss/code-builder/pkg/queueHelper/model"

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
func CreateGitService(req queuemodel.Request) GitService {
	var g GitService
	g.cloneDir = "/tmp/cloneRepo"
	g.repoUrl = req.RepoUrl
	g.username = req.Username
	g.password = req.Password
	// g.sshKey = req.Ssh_key
	return g
}
