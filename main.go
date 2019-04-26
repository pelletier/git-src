package main

import (
	"os"
	"fmt"
	"os/user"
	"path"
	"strings"
	"os/exec"
)

type environment interface {
	Getenv(key string) string
}

func getFromEnvironment(e environment, key string, def string) string {
	val := e.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

type osEnvironment struct {
}

func (e *osEnvironment) Getenv(key string) string {
	return os.Getenv(key)
}

type config struct {
	SrcRootPath string // Root path of src. fully expended.
	DefaultOrg string // Default organization.
	WorkingDirectory string // Current working directory
}

func configFromEnv(env environment) *config {
	usr, _ := user.Current()
	defaultSrcRootPath := path.Join(usr.HomeDir, "src")
	return &config{
		SrcRootPath: getFromEnvironment(env, "GIT_SRC_ROOT", defaultSrcRootPath),
		DefaultOrg: getFromEnvironment(env, "GIT_SRC_DEFAULT_ORG", "github.com"),
		WorkingDirectory: getFromEnvironment(env, "PWD", ""),
	}
}

func (c *config) parseWorkingDirectory() (org string, owner string) {
	if !strings.HasPrefix(c.WorkingDirectory, c.SrcRootPath) {
		return
	}
	noRoot := strings.TrimPrefix(c.WorkingDirectory, c.SrcRootPath)
	noLeadingSlash := strings.TrimLeft(noRoot, "/")
	noTrailingSlash := strings.TrimRight(noLeadingSlash, "/")
	components := strings.SplitN(noTrailingSlash, "/", 3)
	org = components[0]
	if len(components) > 1 {
		owner = components[1]
	}
	return
}

type gitOps interface {
	clone(repo string, target string) error
	exists(path string) bool
	create(path string) error
}

type gitOpsCli struct {
}

func (g *gitOpsCli) clone(repo string, target string) error {
	cmd := exec.Command("git", "clone", "--recursive", repo, target)
	stdoutStderr, err := cmd.CombinedOutput()
	errPrintln(string(stdoutStderr)) // would be better as a stream
	if err != nil {
		return err
	}
	return nil
}

func (g *gitOpsCli) exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (g *gitOpsCli) create(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func errPrintln(args... interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func usage(reason string) {
	if len(reason) > 0 {
		errPrintln("fatal:", reason)
		errPrintln()
	}
	errPrintln(`usage: git src [repo]

When [repo] is omitted, git-src returns $GIT_SRC_ROOT.

[repo] is in the format [[org/]owner/]repository.

For example: github.com/pelletier/git-src.

[org] and [owner] can be omitted when the current working directory is set
appropriately. For example:

PWD=$GIT_SRC_ROOT/github.com/pelletier git-src  => github.com/pelletier/git-src
PWD=$GIT_SRC_ROOT/github.com/ pelletier/git-src => github.com/pelletier/git-src`)
}

func src(env environment, git gitOps, args []string) (string, int) {
	config := configFromEnv(env)

	// Print root path when no argument provided, so that you can just type `src` to get
	// there. Not sure it's the best behavior yet.
	if len(args) == 0 {
		return config.SrcRootPath, 0
	}

	if len(args) > 1 {
		usage("git-src only accepts one argument")
		return "", 129
	}

	repoArg := args[0]
	components := strings.Split(repoArg, "/")

	if len(components) > 3 {
		usage("malformed repository: too many /")
		return "", 129
	}

	wdOrg, wdOwner := config.parseWorkingDirectory()

	organization := config.DefaultOrg
	if wdOrg != "" {
		organization = wdOrg
	}

	owner := wdOwner
	repo := ""

	switch len(components) {
	case 3:
		organization = components[0]
		owner = components[1]
		repo = components[2]
	case 2:
		owner = components[0]
		repo = components[1]
	case 1:
		repo = components[0]
	}

	if organization == "" {
		usage("could not figure out organization")
		return "", 129
	}

	if owner == "" {
		usage("could not figure out owner")
		return "", 129
	}

	ownerPath := path.Join(config.SrcRootPath, organization, owner)
	target := path.Join(ownerPath, repo)

	if !git.exists(target) {
		origin := fmt.Sprintf("git@%s:%s/%s.git", organization, owner, repo)
		errPrintln("repository not found. cloning from", origin)
		err := git.create(ownerPath)
		if err != nil {
			errPrintln("could not create directory", ownerPath, "because:", err)
			return "", 1
		}

		err = git.clone(origin, target)
		if err != nil {
			errPrintln("could not clone:", err)
			return "", 1
		}
	}

	return target, 0
}

func main() {
	env := new(osEnvironment)
	git := &gitOpsCli{}
	output, exitCode := src(env, git, os.Args[1:])
	fmt.Println(output)
	os.Exit(exitCode)
}
