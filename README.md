# No Place Like Home: a quick dotfile linker

## Installation

```bash
bash <(curl nplh.cf/dl)
```

## Setup

<details><summary>My dotfiles are in a repo</summary><p>

1. `nplh install yourusername/yourdotfiles` will install your dotfiles to ~/dotfiles
2. Setup your `~/dotfiles/nplh.yml`, as seen in [Configuration](#configuration)
3. Run `nplh link` to link the dotfiles to their proper locations

</p></details>

<details><summary>My dotfiles are not in a repo</summary><p>

1. Make a `~/dotfiles` directory
2. Move the dotfiles you want to keep to your `~/dotfiles` directory
3. As you move them, add them to `~/dotfiles/nplh.yml`, as seen in [Configuration](#configuration)
4. When you're done, run `nplh link` 

</p></details>


## Configuration

`~/dotfiles/nplh.yml` is where you designate the locations that files should be linked to.
The keys in that file are the files in your dotfiles directory, and the values are the
locations in the file system that they should be linked to. The values can be strings or lists of strings.

### Example

```yaml
vimrc:
  - ~/.config/nvim/init.vim
  - ~/.vimrc
zshrc: ~/.zshrc
i3: ~/.config/i3/config
i3status: ~/.config/i3status/config
xinitrc: ~/.xinitrc
zprofile: ~/.zprofile
gtkrc-2.0: ~/.gtkrc-2.0
gtk-3.0: ~/.config/gtk-3.0
agignore: ~/.agignore
gitconfig: ~/.gitconfig
gitignore: ~/.gitignore
bashrc: ~/.bashrc
dunstrc: ~/.config/dunst/dunstrc
pure.zsh: ~/.zfunctions/prompt_pure_setup
async.zsh: ~/.zfunctions/async
```

## Usage

```
COMMANDS:
     init        setup a dotfiles directory
     install, i  install a dotfiles repo
     update, u   update your dotfiles repo
     link, l     link out the files to their corresponding homes
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

```
