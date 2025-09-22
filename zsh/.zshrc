# Path to your oh-my-zsh installation.
export ZSH=$HOME/.oh-my-zsh

ZSH_THEME="af-magic"

plugins=(git)

source $ZSH/oh-my-zsh.sh

# Vim keybindings
bindkey -v
bindkey "^[OA" up-line-or-beginning-search
bindkey "^[OB" down-line-or-beginning-search

# Custom Settings
setopt noincappendhistory
setopt nosharehistory

# Custom Aliases
alias graph="snakeviz -H `hostname` -p 8082 -s"
alias vim="nvim"

# Custom Environment Variables
export GOPATH=$HOME/golang
export PATH=$PATH:$HOME/golang/bin:/usr/local/go/bin:$HOME/bin:/opt/homebrew/bin
export EDITOR=vim
export VISUAL=vim
