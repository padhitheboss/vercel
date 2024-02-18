package git

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (g *GitService) Clone() error {
	// If no authentication required to clone clone it
	isAuthRequired := g.CheckIfAuthRequired()
	fmt.Println(isAuthRequired)
	if !isAuthRequired {
		return g.NoAuthClone()
	}
	// If authentication required then give first priority to SSH Key then to basic username and password
	if g.sshKey != "" {
		return g.SSHAuthClone()
	} else if (g.username != "") && (g.password != "") {
		return g.BasicAuthClone()
	} else {
		return fmt.Errorf("git authentication failure")
	}
}
func (g *GitService) NoAuthClone() error {
	_, err := git.PlainClone(g.cloneDir, false, &git.CloneOptions{
		URL:      g.repoUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		g.status = "error"
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}

func (g *GitService) CheckIfAuthRequired() bool {
	req, err := http.NewRequest("HEAD", g.repoUrl, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return false
	} else if resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}

func (g *GitService) BasicAuthClone() error {
	auth := &githttp.BasicAuth{
		Username: g.username,
		Password: g.password,
	}

	// Clone the Git repository with authentication
	_, err := git.PlainClone(g.cloneDir, false, &git.CloneOptions{
		URL:      g.repoUrl,
		Auth:     auth,
		Progress: os.Stdout,
	})
	if err != nil {
		g.status = "error"
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}

func (g *GitService) SSHAuthClone() error {
	auth, err := ssh.NewPublicKeys("git", []byte(g.sshKey), "")
	if err != nil {
		return fmt.Errorf("failed to create SSH public keys: %v", err)
	}
	// Clone the Git repository with authentication
	_, err = git.PlainClone(g.cloneDir, false, &git.CloneOptions{
		URL:      g.repoUrl,
		Auth:     auth,
		Progress: os.Stdout,
	})
	if err != nil {
		g.status = "error"
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}
