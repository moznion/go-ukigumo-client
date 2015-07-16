package ukigumo

import (
	"fmt"
	"log"
)

type VCGit struct {
	vcInfo   *VCInfo
	LogLimit uint
}

func NewVCGit(repository string, branch string) *VCGit {
	vg := new(VCGit)
	vg.vcInfo = &VCInfo{
		Repository:       repository,
		Branch:           branch,
		Description:      "",
		SkipIfUnmodified: false,
		revision:         "",
	}
	vg.LogLimit = 50

	return vg
}

func (vg *VCGit) GetLog(revFrom string, revTo string) (string, error) {
	cmd := "git log --pretty=format:\"%h %an: %s\" --abbrev-commit --source "
	if revFrom == revTo {
		cmd += "-1 'HEAD'"
	} else {
		cmd += fmt.Sprintf("-%d '%s..%s'", vg.LogLimit, revFrom, revTo)
	}

	out, err := executeCommand(cmd)
	if err != nil {
		log.Fatal(err)
	}
	return out, err
}

func (vg *VCGit) Update() error {
	var err error
	_, err = executeCommandWithOutput("git pull -f origin " + vg.vcInfo.Branch)
	if err != nil {
		return err
	}
	_, err = executeCommandWithOutput("git submodule init")
	if err != nil {
		return err
	}
	_, err = executeCommandWithOutput("git submodule update")
	if err != nil {
		return err
	}
	_, err = executeCommandWithOutput("git clean -dxf")
	if err != nil {
		return err
	}
	_, err = executeCommandWithOutput("git status")
	if err != nil {
		return err
	}

	rev, err := executeCommand("git rev-parse HEAD")
	if err != nil {
		rev = "Unknown"
	}
	vg.vcInfo.revision = rev

	return nil
}

func (vg *VCGit) GetVCInfo() *VCInfo {
	return vg.vcInfo
}

func (vg *VCGit) GetRevision() string {
	rev := vg.vcInfo.revision
	if rev != "" {
		return rev
	}

	rev, err := executeCommand("git rev-parse HEAD")
	if err != nil {
		return "Unknown"
	}
	return rev
}
