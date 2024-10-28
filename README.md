# **my-ls: A Custom Implementation of the ls Command in Go**

## **Objective**

The objective of this project is to create a custom implementation of the `ls` command in Go, which displays the files and folders of a specified directory or the current directory if no directory is specified. The project aims to replicate the behavior of the original `ls` command with the following variations:

- Incorporating at least the following flags: `-l`, `-R`, `-a`, `-r`, and `-t`
- Displaying files and folders in a specific order

## **Features**

- Displays files and folders of a specified directory or the current directory
- Supports the following flags:
  - `-l`: displays files and folders in a long format
  - `-R`: recursively displays files and folders in subdirectories
  - `-a`: displays all files and folders, including hidden ones
  - `-r`: reverses the order of files and folders
  - `-t`: sorts files and folders by modification time
- Written in Go, following good practices and using only allowed packages:
  - `fmt`
  - `os`
  - `os/user`
  - `strconv`
  - `strings`
  - `syscall`
  - `time`
  - `math/rand`
  - `errors`
  - `io/fs`
- Includes test files for unit testing
- Does not use the `os/exec` package

## **Implementation Notes**

- The project takes into account the implications of the `-R` flag from the beginning of the code.
- The order of files and folders is carefully considered and implemented.
- The project consults the `ls` command manual to ensure accuracy and consistency.

## **Learning Outcomes**

This project helps to learn about:

- Unix system
- Ways to receive data
- Ways to output data
- Manipulation of strings
- Manipulation of structures

## **Usage**

To use the `my-ls` command, simply compile the Go code and run the executable in your terminal. You can use the command with various flags, such as:

VerifyOpen In EditorEditCopy code

`1my-ls -l 2my-ls -R 3my-ls -a 4my-ls -r 5my-ls -t`

You can also specify a directory as an argument, such as:

VerifyOpen In EditorEditCopy code

`1my-ls /path/to/directory`

## **License**

This project is licensed under the MIT License. See the `LICENSE` file for details.
