package main

import (
	"errors"
	"log"
	"os"
	"reflect"
	"time"
)

type Runner struct {
	conf           *Config
	elapsedTimeSec *int64
}

func NewRunner(conf *Config, elapsedTimeSec *int64) *Runner {
	r := new(Runner)
	r.conf = conf
	r.elapsedTimeSec = elapsedTimeSec

	return r
}

func (r *Runner) Run(phase string) error {
	var commands []string = nil
	if phase == "Install" {
		commands = r.decideInstallCommand()
	} else if phase == "Script" {
		return errors.New("Must not run as 'Script' phase")
	} else {
		commands = reflect.ValueOf(r.conf).Elem().FieldByName(phase).Interface().([]string)
	}

	if commands != nil {
		for _, cmd := range commands {
			log.Printf("[%s] %s", phase, cmd)

			start := int64(time.Now().Unix())

			_, err := executeCommandWithOutput(cmd)
			if err != nil {
				log.Fatalf("Failure in %s: %s", phase, cmd)
				return err
			}

			elapsed := *r.elapsedTimeSec + int64(time.Now().Unix()) - start
			r.elapsedTimeSec = &elapsed
		}
	}

	return nil
}

func (r *Runner) decideInstallCommand() []string {
	installCommands := r.conf.Install

	if installCommands != nil {
		return installCommands
	}

	installCommandsForPerl := []string{"cpanm --notest --installdeps ."}
	for _, file := range []string{"Makefile.PL", "cpanfile", "Build.PL"} {
		_, err := os.Stat(file)
		if err == nil {
			return installCommandsForPerl
		}
	}

	return nil
}
