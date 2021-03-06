// Copyright 2016 Michael Schenk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package command

import (
	"flag"
	"strings"
	"time"

	"github.com/git-time-metric/gtm/metric"
	"github.com/git-time-metric/gtm/note"
	"github.com/git-time-metric/gtm/project"
	"github.com/git-time-metric/gtm/report"
	"github.com/git-time-metric/gtm/util"
	"github.com/mitchellh/cli"
)

// StatusCmd containt methods for status command
type StatusCmd struct {
	Ui cli.Ui
}

// NewStatus returns new StatusCmd struct
func NewStatus() (cli.Command, error) {
	return StatusCmd{}, nil
}

// Help returns help for status command
func (c StatusCmd) Help() string {
	helpText := `
Usage: gtm status [options]

  Show pending time for git repositories.

Options:

  -terminal-off=false        Exclude time spent in terminal (Terminal plug-in is required)

  -color=false               Always output color even if no terminal is detected, i.e 'gtm status -color | less -R'

  -total-only=false          Only display total pending time

  -tags=""                   Project tags to report status for, i.e --tags tag1,tag2

  -all=false                 Show status for all projects
`
	return strings.TrimSpace(helpText)
}

// Run executes status command with args
func (c StatusCmd) Run(args []string) int {
	var color, terminalOff, totalOnly, all, profile bool
	var tags string
	cmdFlags := flag.NewFlagSet("status", flag.ContinueOnError)
	cmdFlags.BoolVar(&color, "color", false, "Always output color even if no terminal is detected. Use this with pagers i.e 'less -R' or 'more -R'")
	cmdFlags.BoolVar(&terminalOff, "terminal-off", false, "Exclude time spent in terminal (Terminal plugin is required)")
	cmdFlags.BoolVar(&totalOnly, "total-only", false, "Only display total time")
	cmdFlags.StringVar(&tags, "tags", "", "Project tags to show status on")
	cmdFlags.BoolVar(&all, "all", false, "Show status for all projects")
	cmdFlags.BoolVar(&profile, "profile", false, "Enable profiling")
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	util.TimeTrackEnable = profile
	defer util.TimeTrack(time.Now(), "status.Run")

	if totalOnly && (all || tags != "") {
		c.Ui.Error("\n-tags and -all options not allowed with -total-only\n")
		return 1
	}

	var (
		err        error
		commitNote note.CommitNote
		out        string
	)

	index, err := project.NewIndex()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	tagList := []string{}
	if tags != "" {
		tagList = util.Map(strings.Split(tags, ","), strings.TrimSpace)
	}

	projects, err := index.Get(tagList, all)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	options := report.OutputOptions{
		TotalOnly:   totalOnly,
		TerminalOff: terminalOff,
		Color:       color}

	for _, projPath := range projects {
		if commitNote, err = metric.Process(true, projPath); err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
		o, err := report.Status(commitNote, options, projPath)
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
		out += o
	}

	c.Ui.Output(out)
	return 0
}

// Synopsis returns help for status command
func (c StatusCmd) Synopsis() string {
	return "Show pending time"
}
