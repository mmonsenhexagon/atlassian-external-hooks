package git

import (
	"strings"

	"github.com/reconquest/atlassian-external-hooks/integration_tests/internal/exec"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/lexec-go"
)

type Git struct {
	dir string
}

func Clone(href string, path string) (*Git, error) {
	err := exec.New("git", "clone", href, path).Run()
	if err != nil {
		return nil, karma.Format(
			err,
			"run git clone",
		)
	}

	git := &Git{
		dir: path,
	}

	return git, nil
}

func (git *Git) GetWorkDir() string {
	return git.dir
}

func (git *Git) Add(paths ...string) error {
	return git.command("add", paths...).Run()
}

func (git *Git) Commit(message string) error {
	return git.command("commit", "-m", message).Run()
}

func (git *Git) Push(args ...string) (string, error) {
	_, stderr, err := git.command("push", args...).Output()
	return string(stderr), err
}

func (git *Git) Branch(name string) error {
	err := git.command("checkout", "-b", name).Run()
	if err != nil {
		return err
	}

	return nil
}

func (git *Git) RevList(args ...string) ([]string, error) {
	stdout, _, err := git.command("rev-list", args...).Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(stdout)), "\n"), nil
}

func (git *Git) command(command string, args ...string) *lexec.Execution {
	args = append([]string{"-C", git.dir, command}, args...)

	return exec.New("git", args...)
}
