" Include the system settings
:if filereadable( "/etc/vimrc" )
   source /etc/vimrc
:endif

" If we are using ZSH set the theme to solarized
:if $ZSH == "$HOME/.oh-my-zsh"
"colorscheme solarized
:endif

" If this is not above `higlight` it overwrites it.
set background=dark
"syntax on
set nocompatible
filetype off

" These are commented out so I get the lighter navy color for the columns
:highlight ColorColumn guibg=lightblue
:highlight ColorColumn ctermbg=blue

" Matches extra whitespace at end of line
:highlight ExtraWhitespace ctermbg=red guibg=red

:match ExtraWhitespace /\s\+$/

" Installs vim-plug
if empty(glob('~/.vim/autoload/plug.vim'))
   silent !curl -fLo ~/.vim/autoload/plug.vim --create-dirs
      \ https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
   autocmd VimEnter * PlugInstall --sync | source $MYVIMRC
endif

call plug#begin()

Plug 'fatih/vim-go'
Plug 'majutsushi/tagbar'
Plug 'yegappan/lid'
Plug 'Shougo/vimshell.vim'
Plug 'Shougo/vimproc.vim', {'do': 'make'}
Plug 'scrooloose/nerdtree'
Plug 'Xuyuanp/nerdtree-git-plugin'
Plug 'Shougo/neocomplete', {'do': 'make'}
Plug 'w0rp/ale'

call plug#end()

" neocomplete settings
let g:neocomplete#enable_at_startup = 1
let g:neocomplete#enable_smart_case = 1
let g:neocomplete#sources#syntax#min_keyword_length = 3

" Tagbar settings
let g:tagbar_left = 1
let g:tagbar_autofocus = 1
let g:tagbar_width = 40

" ALE settings
nmap <silent> <C-k> <Plug>(ale_previous_wrap)
nmap <silent> <C-j> <Plug>(ale_next_wrap)

set tags+=/golang/tags,/src/gated/tags,/src/Bgp/tags,/src/MultiProtocolRouting/tags,/src/TAGS,/src/tags

" Some visual settings
set ruler
set number
set hlsearch
"set foldmethod=syntax

"filetype plugin indent on

" Settings for vim-go
let g:go_term_enabled = 1

" Removes whitespaces from end of lines.
function CleanWhitespace()
   exe '%s/ \+$//g'
endfunction

" Rebindings
nnoremap CleanWhitespace :CleanWhitespace<CR>
nnoremap <leader>tree :NERDTree<CR>
nnoremap <leader>id   :Lid<CR>

" Settings for untyped files
au BufWinEnter *      set ts=3 shiftwidth=3 colorcolumn=86,101
au BufWinEnter *.log  set ts=4 shiftwidth=4 colorcolumn=0
au BufWinEnter *.go   set ts=4 shiftwidth=4 colorcolumn=101
au BufWinEnter *.c    set ts=4 shiftwidth=4 colorcolumn=101
au BufWinEnter *.h    set ts=4 shiftwidth=4 colorcolumn=101
au BufWinEnter *.py   set ts=3 shiftwidth=3 colorcolumn=86
au BufWinEnter *.html set ts=4 shiftwidth=4 colorcolumn=101
