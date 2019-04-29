package main

import (
	"testing"
)

type mapEnv struct {
	data map[string]string
}

func makeMapEnv(kvs ...string) *mapEnv {
	if len(kvs)%2 != 0 {
		panic("makeMapEnv: arguments must go in pairs")
	}
	m := map[string]string{}
	for i := 0; i < len(kvs)/2; i++ {
		k := kvs[i*2]
		v := kvs[i*2+1]
		m[k] = v
	}

	return &mapEnv{
		data: m,
	}
}

func (e *mapEnv) Getenv(key string) string {
	val, ok := e.data[key]
	if !ok {
		return ""
	}
	return val
}

var defaultEnv environment = makeMapEnv("GIT_SRC_ROOT", "/blah")

type fakeGit struct {
	cloned       bool
	repo         string
	target       string
	parentExists bool
	err          error
}

func (g *fakeGit) clone(repo string, target string) error {
	g.cloned = true
	g.repo = repo
	g.target = target
	return g.err
}

func (g *fakeGit) exists(path string) bool {
	return g.parentExists
}

func (g *fakeGit) create(path string) error {
	return nil
}

func makeGit() *fakeGit {
	return new(fakeGit)
}

func assertSrc(t *testing.T, output string, expectedOutput string, exitCode int,
	expectedExitCode int) {
	if output != expectedOutput {
		t.Error("expected output", expectedOutput, "got", output)
	}
	if exitCode != expectedExitCode {
		t.Error("expected exit code", expectedExitCode, "got", exitCode)
	}
}

func TestNoArg(t *testing.T) {
	output, exitCode := src(defaultEnv, makeGit(), []string{})
	assertSrc(t, output, "/blah", exitCode, 0)

}

func TestTooManyArgs(t *testing.T) {
	output, exitCode := src(defaultEnv, makeGit(), []string{"hello", "world"})
	assertSrc(t, output, "", exitCode, 129)
}

func TestFullRepo(t *testing.T) {
	git := makeGit()
	output, exitCode := src(defaultEnv, git, []string{"github.com/pelletier/go-src"})
	assertSrc(t, output, "/blah/github.com/pelletier/go-src", exitCode, 0)
	if !git.cloned {
		t.Error("repository should have been cloned")
	}
	if git.repo != "git@github.com:pelletier/go-src.git" {
		t.Error("cloned wrong repo. expected go-src, got", git.repo)
	}
	if git.target != "/blah/github.com/pelletier/go-src" {
		t.Error("cloned to wrng target. expected , got", git.target)
	}
}

func TestExtraSlashes(t *testing.T) {
	output, exitCode := src(defaultEnv, makeGit(), []string{"github.com/pelletier/go-src/whatever"})
	assertSrc(t, output, "", exitCode, 129)
}

func TestPartialInWd(t *testing.T) {
	e := makeMapEnv(
		"GIT_SRC_ROOT", "/home/me/src",
		"PWD", "/home/me/src/github.com/pelletier/other-repo/in/the/weeds",
	)
	output, exitCode := src(e, makeGit(), []string{"go-src"})
	assertSrc(t, output, "/home/me/src/github.com/pelletier/go-src", exitCode, 0)
	output, exitCode = src(e, makeGit(), []string{"pelletier/go-src"})
	assertSrc(t, output, "/home/me/src/github.com/pelletier/go-src", exitCode, 0)
	output, exitCode = src(e, makeGit(), []string{"other/repo"})
	assertSrc(t, output, "/home/me/src/github.com/other/repo", exitCode, 0)
	output, exitCode = src(e, makeGit(), []string{"mycorp.com/something/else"})
	assertSrc(t, output, "/home/me/src/mycorp.com/something/else", exitCode, 0)
}

func TestOutOfWd(t *testing.T) {
	e := makeMapEnv(
		"GIT_SRC_ROOT", "/home/me/src",
		"PWD", "/home/me/somewhere/else",
	)
	output, exitCode := src(e, makeGit(), []string{"mycorp.com/something/else"})
	assertSrc(t, output, "/home/me/src/mycorp.com/something/else", exitCode, 0)
}

func TestOutOfWdNoOrg(t *testing.T) {
	e := makeMapEnv(
		"GIT_SRC_ROOT", "/home/me/src",
		"PWD", "/home/me/somewhere/else",
	)
	output, exitCode := src(e, makeGit(), []string{"something/else"})
	assertSrc(t, output, "/home/me/src/github.com/something/else", exitCode, 0)
}

func TestEmptySlashes(t *testing.T) {
	e := makeMapEnv(
		"GIT_SRC_ROOT", "/home/me/src",
		"PWD", "/home/me/somewhere/else",
	)
	output, exitCode := src(e, makeGit(), []string{"//"})
	assertSrc(t, output, "", exitCode, 129)
	output, exitCode = src(e, makeGit(), []string{"/"})
	assertSrc(t, output, "", exitCode, 129)
}

func TestEmptySlashesInWd(t *testing.T) {
	e := makeMapEnv(
		"GIT_SRC_ROOT", "/home/me/src",
		"PWD", "/home/me/src/github.com/pelletier/other-repo/in/the/weeds",
	)
	output, exitCode := src(e, makeGit(), []string{"//"})
	assertSrc(t, output, "", exitCode, 129)
	output, exitCode = src(e, makeGit(), []string{"/"})
	assertSrc(t, output, "", exitCode, 129)
}
