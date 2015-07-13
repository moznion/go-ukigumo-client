package main

type VCInfo struct {
	Repository       string
	Branch           string
	Description      string
	SkipIfUnmodified bool

	revision string
}

type VC interface {
	GetVCInfo() *VCInfo
	GetLog(revFrom string, revTo string) (string, error)
	Update() error
	GetRevision() string
}
