# Install Oh-My-ZSH
sh -c "$(wget https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh -O -)"

# Install github.com/junegunn/vim-plug
mkdir -p ~/.vim/autoload/
wget -O ~/.vim/autoload/plug.vim https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim

# Backup existing dotfiles if they exist.
mv $HOME/.zshrc $HOME/.zshrc.old
mv $HOME/.vimrc $HOME/.vimrc.old
mv $HOME/.bashrc $HOME/.bashrc.old
mv $HOME/.tmux.conf $HOME/.tmux.conf.old

# Symlink our dotfiles to the right location.
ln -s $PWD/zshrc $HOME/.zshrc
ln -s $PWD/vimrc $HOME/.vimrc
ln -s $PWD/bashrc $HOME/.bashrc
ln -s $PWD/tmux.conf $HOME/.tmux.conf
