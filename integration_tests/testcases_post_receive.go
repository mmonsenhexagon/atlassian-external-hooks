package main

import (
	"github.com/kovetskiy/stash"
	"github.com/reconquest/atlassian-external-hooks/integration_tests/internal/external_hooks"
	"github.com/reconquest/atlassian-external-hooks/integration_tests/internal/runner"
	"github.com/stretchr/testify/assert"
)

func Testcase_PostReceive_OutputMessage(
	run *runner.Runner,
	assert *assert.Assertions,
	project *stash.Project,
	repository *stash.Repository,
) {
	err := run.Bitbucket().WriteFile(
		"shared/external-hooks/fail.sh",
		text(
			`#!/bin/bash`,
			`echo XXX`,
			`exit 1`,
		),
		0777,
	)
	assert.NoError(err, "should be able to write hook script to container")

	addon := run.ExternalHooks()

	postReceive := addon.OnProject(project.Key).PostReceive(
		external_hooks.NewSettings().
			UseSafePath(true).
			WithExecutable("fail.sh"),
	)

	err = postReceive.Configure()
	assert.NoError(err, "should be able to configure post-receive hook")

	err = postReceive.Enable()
	assert.NoError(err, "should be able to enable post-receive hook")

	git := run.GitClone(repository)

	run.GitCommitRandomFile(git)

	stdout, err := git.Push()
	assert.NoError(err, "git push should succeed")
	assert.Contains(
		string(stdout),
		"remote: XXX",
		"git push message should contain output from post-receive hook",
	)

	err = postReceive.Disable()
	assert.NoError(err, "should be able to disable post-receive hook")

	run.GitCommitRandomFile(git)

	stdout, err = git.Push()
	assert.NoError(err, "git push should succeed")
	assert.NotContains(
		string(stdout),
		"remote: external-post-receive-hook declined",
		"post-receive-hook should not decline push",
	)
	assert.NotContains(
		string(stdout),
		"remote: XXX",
		"post-receive-hook should not contain output from post-receive hook",
	)
}
