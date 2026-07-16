dotfiles
===

## Overview

This repo manages my dotfiles across various environments using a combination of `Stow` and Go Templates.

### Configuration

Our configuration file is in the following format:
```yaml
# List of valid environments that we can choose from on first launch.
environments:
  - home-pc
  - home-laptop
  - ...

# These globals will be available to all applications during template rendering.
globals:
  foo: true

# Per-app configurations go here.
apps:
  nvim:
    # Applications can be enabled/disabled on a per-environment basis.
    enabled: {{ eq .Environment home-pc }}
    # These variables are only passed to the given app during template rendering.
    vars:
        bar: false
```


### Templates

Since I don't want every application installed exactly the same way on all of my development environments we can manipulate the configs (and determine whether we even bother installing the app in the first place) through the use of the above configuration file and per-app templates.

Any file containing a `.tmpl` suffix will be passed through our template rendering step with the following available fields:
```go
type AppTemplateData struct {
    Environment string
    Globals map[string]any
    Vars map[string]any
}
```

These can then be used to conditionally tweak a configuration depending on where it's supposed to be running.

### Installation Flow

For each app we run the following steps:

1. Run `pre.sh`
    - Typically this will do things like installing the given application (nvim, zellij, etc..)

2. "Render" our files into the /build directory in
    - Any template files will be passed through our templating script and written to the build directory with the `.tmpl` suffix removed.
    - Any other files will be copied over exactly as they are.
    - Note that we skip-over `pre.sh` and `post.sh`

3. Use `stow` to link our files to our home directory.

4. Run `post.sh`
    - This typically takes care of any post-install steps (such as installing oh-my-zsh or downloading app related plugins)

## Usage

```bash
# Build and install configs for the current environment
# (you'll be prompted to choose an environment on first run)
make install

# Render only (skip symlinking)
make render

# View pending changes before installing
make diff
```
