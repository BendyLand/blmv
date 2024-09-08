# blmv

A basic batch move tool written in Go!

**Note: This tool has only been tested in controlled settings. Behavior may be unpredictable. Please do NOT use for important data.**

## Usage

This tool expects a file called "config.txt" to be found in the directory where it is run from.

 - The first line of this file should be the source directory relative to the root directory.
 - The second line of this file should be the destination directory relative to the root directory.
 - Config files that do not conform to this structure may cause unpredictable behavior.
 
Personally, I suggest building the binary using `go build .` . Then, move the resulting executable to a place where your system can use it from (`sudo mv blmv /usr/local/bin` or similar, depending on your OS), or just keep a copy where you may need it. Then, simply run the executable. 
 
## Example

In the root directory of this project, there is a folder called `example_dir`. It has the following structure:
```bash
example_dir/
├── config.txt
├── dst/
└── src/
    ├── file0.txt
    ├── file1.txt
    ├── ...
    └── file99.txt
```

The contexts of `config.txt` are as follows:
```txt
src/
dst/
```

If the command is set up in such as way that it can be run from anywhere in the terminal, you can simply run:
```bash
example_dir % blmv
```

Otherwise, you may need to keep a copy of the executable in the directory and run it like you would any other executable file (.exe):
```batch
blmv.exe
```
