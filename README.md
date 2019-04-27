# git-src

git-src helps jumping around git repositories in your shell.

## Usage

*Usage:* `src <repo>`: `cd` to a repository, cloning it if it is missing.

First you need to define a `GIT_SRC_ROOT` path. For example: `$HOME/src`.

From then on, `git-src` assumes that all your repository are located in
`$HOME/src/<organization>/<owner>/<repo>`. Same as [Go workspaces][go-workspaces].
For example: `$HOME/src/github.com/pelletier/git-src`.

`git-src` will clone (if necessary) and `cd` into any repository in that
structure, based on the provided `<repo>` argument and your current working
directory. Here are some examples:

```
# From anywhere, go to github.com/pelletier/git-src
$ src github.com/pelletier/git-src

# From anywhere within github.com/pelletier/
# Go to github.com/pelletier/foo
$ src foo

# From anywhere within github.com/
# Go to github.com/hello/world
$ src hello/world
```


## Installation

1. Get the latest [git-src binary][binaries] for your system.
2. Rename the binary to `git-src`.
3. Move it to somewhere on your `PATH`. Make sure it is executable.
4. Add the following to your shell config (e.g. `.bashrc`):

```
# src
function src {
    new_path=$(git src $*)
    if [ $? = 0 ]; then
        cd ${new_path}
    fi
}
```

[binaries]: https://github.com/pelletier/git-src/releases/tag/master
[go-workspaces]: https://golang.org/doc/code.html#Workspaces
