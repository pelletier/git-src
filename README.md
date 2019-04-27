# git-src

Usage example: `src pelletier/git-src`
Clone this repository if not already present, and navigate to it.

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
