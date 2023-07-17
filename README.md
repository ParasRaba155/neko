## Neko

This is a cheap knock off of majestic linux command `cat`

#### USAGE

```
Usage: neko [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

  -b, --number-nonblank    number nonempty output lines, overrides -n
  -e, --show-ends          display $ at end of each line
  -n, --number             number all output lines
  -t, --show-tabs          display TAB characters as ^I
  -v, --show-nonprinting   use ^ and M- notation, except for LFD and TAB
      --help     display this help and exit

Examples:
  neko f - g  Output f's contents, then standard input, then g's contents.
  neko        Copy standard input to standard output.

```
