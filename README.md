# dmenu_alias

This is a simple program that was written to provide a method to "alias" a dmenu
output string into a different string.

## Setup
First install [dmenu](https://tools.suckless.org/dmenu/):

```
git clone https://git.suckless.org/dmenu
cd dmenu
make
sudo make install
```

Install this program:
```
go build .
sudo ./install
```

Update `dmenu_run` to the following: 
```
dmenu_path | dmenu_alias -i | dmenu | dmenu_alias -o | sh
```

By default, this program searches in `$HOME/.config/dmenu_alias_list` so
create a file there with your aliases.

If you wish to uninstall this program, run `sudo ./uninstall`

## dmenu_alias_list
The syntax is quite simple: `string1=string2`. There must not be a space
between the equal sign and the two strings. Each alias should be seperated by
a newline.

Example:
```
vi=xterm -e vi 
```

Here, if the user were to enter `vi` into dmenu, then it can then be
piped into `dmenu_alias` to be turned into `xterm -e vi` which can then
be given as a command to the shell so that `vi` can be run in a terminal
window (this makes `dmenu_alias` useful for running terminal apps directly from
dmenu).

### Escape Sequences

If your text contains `=` or `\`, use `\=` and `\\` to represent them.

## Usage
```
dmenu_alias [-i|-o] [alias list file]
```

If `-i` is given to `dmenu_alias`, then it will read from `stdin` and output
what it reads along with the aliases as input into dmenu so that aliases actually
appear.

If `-o` is given to `dmenu_alias` then it will read from `stdin` and output
what the input should be aliased to. This allows the program to take dmenu
output and convert it to its alias. `-o` is the default option.

## License 
MIT
