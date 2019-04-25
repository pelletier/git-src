# git-src

Usage example: `src pelletier/git-src`
Clone this repository if not already present, and navigate to it.

## Installation

```
ln -s git-src ~/bin/
cat >> ~/.zshrc <<SCRIPT
# src
function src {
    new_path=$(git src $*)
    if [ $? = 0 ]; then
        cd ${new_path}
    fi
}
SCRIPT
```
