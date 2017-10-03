# flatten
Flatten is a command-line tool to flatten a directory structure.

## Usage
```
Usage:
  flatten [SOURCE] [DESTINATION] [-c | --copy-only] [-f | --force] [--include-source-files] [-s | --simulate-only] [--verbose]
  flatten -h | --help
  flatten -v

Recursively flatten the directory structure from SOURCE to DESTINATION.

Arguments:
  SOURCE                    Optional source directory (default is current directory).
  DESTINATION               Optional destination directory (default is current directory).

Options:
  -c --copy-only            Do not remove anything from the source directory.
  -f --force                Do not propose a simulation first, immediately execute the command.
  --include-source-files    Include the files which are directly located in the SOURCE directory.
  -s --simulate-only        Do not move or copy any files on the system,
                            just output the expected result.
  --verbose                 Explain what is being done.
  -v --version              Show version.
  -h --help                 Show this screen.
```

## Example
Assume we have the following directory strcuture in `/home/goggle/example/`:
```
/home/goggle/example
├── c_progs
│   └── prog01
│       ├── hello
│       └── hello.c
├── data
│   ├── dat001
│   │   ├── data_apples.txt
│   │   ├── data_monkeys.txt
│   │   └── data_trees.txt
│   ├── dat002
│   │   ├── data_apples.txt
│   │   ├── data_monkeys.txt
│   │   └── data_trees.txt
│   ├── dat003
│   │   ├── data_apples.txt
│   │   ├── data_monkeys.txt
│   │   └── data_trees.txt
│   └── dat004
│       ├── data_apples.txt
│       ├── data_monkeys.txt
│       └── data_trees.txt
├── hello
└── hello_1
```

By running `flatten` in `/home/goggle/example` we get the following result:
```
/home/goggle/example
├── data_apples_1.txt
├── data_apples_2.txt
├── data_apples_3.txt
├── data_apples_4.txt
├── data_monkeys_1.txt
├── data_monkeys_2.txt
├── data_monkeys_3.txt
├── data_monkeys_4.txt
├── data_trees_1.txt
├── data_trees_2.txt
├── data_trees_3.txt
├── data_trees_4.txt
├── hello
├── hello_01
├── hello_1
└── hello.c
```
All the files in the subdirectories of `/home/goggle/example/` have been moved into `/home/goggle/example` and the empty directories have been removed. Note, that no regular file has been removed, even though we have file name collisions (e.g. the file `data_apples.txt` exists four times). Flatten does automatically take care of such filename collisions and adds a number to the filename if such a collision happens.

If we want to keep the original files in their subdirectories, we can use the `--copy-only` option. `flatten -c` or `flatten --copy-only` executed in `/home/goggle/example` will lead to the following result:
```
/home/goggle/example
├── c_progs
│   └── prog01
│       ├── hello
│       └── hello.c
├── data
│   ├── dat001
│   │   ├── data_apples.txt
│   │   ├── data_monkeys.txt
│   │   └── data_trees.txt
│   ├── dat002
│   │   ├── data_apples.txt
│   │   ├── data_monkeys.txt
│   │   └── data_trees.txt
│   ├── dat003
│   │   ├── data_apples.txt
│   │   ├── data_monkeys.txt
│   │   └── data_trees.txt
│   └── dat004
│       ├── data_apples.txt
│       ├── data_monkeys.txt
│       └── data_trees.txt
├── data_apples_1.txt
├── data_apples_2.txt
├── data_apples_3.txt
├── data_apples_4.txt
├── data_monkeys_1.txt
├── data_monkeys_2.txt
├── data_monkeys_3.txt
├── data_monkeys_4.txt
├── data_trees_1.txt
├── data_trees_2.txt
├── data_trees_3.txt
├── data_trees_4.txt
├── hello
├── hello_01
├── hello_1
└── hello.c
```
By default, flatten will perform a simulation of its actions first, and ask the user, if they want to continue.
