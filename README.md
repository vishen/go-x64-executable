# ELF x86-64 Linux Executable by hand

Generate an ELF 64-bit executable by hand. This is an example
of how to manually create an ELF 64 executable for linux that
will print out some string and then exit.

It will generate the ELF header, then 2 Program Headers, one
for the `.text` segment and one for the `.data` segment. Then the 
`.text` segment and `.data` segment are outputted to the binary.

Currently the string and filename are hardcoded in main.go.

This doesn't currently do back-patching. The data segment is
always known since the ELF header and test segment aren't changing.
This needs to be fixed, but I don't know a good way to do it.

```
$ go run main.go
$ ./comp
Hello World, this is my tiny executable
$ ls -hal ./comp 
-rwxr-xr-x 1 vishen vishen 257 Jun 20 20:08 ./comp*
```

## Resources

- https://en.wikipedia.org/wiki/Executable_and_Linkable_Format
- https://www.hanshq.net/making-executables.html#linux
- https://github.com/nivertech/erl_elf_test/blob/master/manually_creating_elf_executable.md
