package pkgbuilder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type BasicBuild struct {
	repoPath   string
	outputPath string
	status     string
}
type BuildReact struct {
	BasicBuild
}

type Build interface {
	BuildStatic() (string, error)
	GetOutputPath() string
}

func CreateConfig(framework, repoPath string) Build {
	switch framework {
	case "react":
		var b BuildReact
		b.repoPath = repoPath
		return &b
	default:
		log.Panicln("invalid project type")
	}
	return nil
}
func (b *BuildReact) BuildStatic() (string, error) {
	var buildLog string = ""
	cmd := exec.Command("npm", "install")
	cmd.Dir = b.repoPath
	output, err := cmd.Output()
	if err != nil {
		b.status = "failed"
		log.Println(err)
		return string(output), err
	}
	buildLog += string(output)
	preBuildFolder := CaptureFolderList(b.repoPath)
	preBuildFolder["node_modules"] = true
	fmt.Println(string(output))
	cmd = exec.Command("npm", "run", "build")
	cmd.Dir = b.repoPath
	output, err = cmd.Output()
	if err != nil {
		b.status = "failed"
		log.Println(err)
		return string(output), err
	}
	buildLog += string(output)
	b.status = "success"
	fmt.Println(string(output))
	postBuildFolder := CaptureFolderList(b.repoPath)
	for key := range postBuildFolder {
		if _, ok := preBuildFolder[key]; !ok {
			b.outputPath = b.repoPath + "/" + key
			return buildLog, nil
		}
	}
	return buildLog, nil
}

func (b *BuildReact) GetOutputPath() string {
	return b.outputPath
}
func CaptureFolderList(path string) map[string]bool {
	allFiles, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	folders := make(map[string]bool)
	for _, file := range allFiles {
		if file.IsDir() {
			folders[file.Name()] = true
		}
	}
	return folders
}
