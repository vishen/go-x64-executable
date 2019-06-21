# ELF x86-64 Linux Executable by hand

Generate an ELF 64-bit executable by hand. This is an example
of how to manually create an ELF 64 executable for linux that
will print out some string and then exit using `sys_write` and
`sys_exit` linus system calls.

It will generate the ELF header, then 2 Program Headers (one
for the `.text` segment and one for the `.data` segment). Then the 
`.text` segment and `.data` segment are outputted to the binary.

The `.text` segment will do a `sys_write` using the data offset
from where the `.data` segment is located and then it will do
a `sys_exit`.

The `.data` segment contains just the string data.

This doesn't currently do back-patching. The data segment is
always known since the ELF header and text segment aren't changing.
This needs to be fixed, but I don't know a good way to do it.

```
$ go run main.go
wrote binary to tiny-x64
$ ./tiny-x64
Hello World, this is my tiny executable
$ ls -hal ./tiny-x64
-rwxr-xr-x 1 vishen vishen 257 Jun 21 08:58 ./tiny-x64*

# Or run with a desired string to output in the binary
$ go run main.go -output hello-world -word "Hello, World!"
wrote binary to hello-world
$ ./hello-world
Hello, World!
$ ls -hal ./hello-world
-rwxr-xr-x 1 vishen vishen 231 Jun 21 08:59 ./hello-world*
```

## Resources

- https://en.wikipedia.org/wiki/Executable_and_Linkable_Format
- https://www.hanshq.net/making-executables.html#linux
- https://github.com/nivertech/erl_elf_test/blob/master/manually_creating_elf_executable.md
- http://www.sco.com/developers/gabi/latest/ch5.pheader.html
