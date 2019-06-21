package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	virtualStartAddress     uint64 = 0x400000
	dataVirtualStartAddress uint64 = 0x600000
	alignment               uint64 = 0x200000
)

type Builder struct {
	o []byte
}

func (b *Builder) WriteBytes(bs ...byte) {
	b.o = append(b.o, bs...)
}

func (b *Builder) WriteValue(size int, value uint64) {
	buf := make([]byte, size)
	binary.LittleEndian.PutUint64(buf, value)
	b.WriteBytes(buf...)
}

func buildELF(textSection, dataSection []byte) []byte {
	textSize := uint64(len(textSection))
	// Size of ELF header + 2 * size program header?
	textOffset := uint64(0x40 + (2 * 0x38))

	var o Builder

	// Build ELF Header
	o.WriteBytes(0x7f, 0x45, 0x4c, 0x46) // ELF magic value

	o.WriteBytes(0x02) // 64-bit executable
	o.WriteBytes(0x01) // Little endian
	o.WriteBytes(0x01) // ELF version
	o.WriteBytes(0x00) // Target OS ABI
	o.WriteBytes(0x00) // Further specify ABI version

	o.WriteBytes(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Unused bytes

	o.WriteBytes(0x02, 0x00)             // Executable type
	o.WriteBytes(0x3e, 0x00)             // x86-64 target architecture
	o.WriteBytes(0x01, 0x00, 0x00, 0x00) // ELF version

	// 64-bit virtual offsets always start at 0x400000?? https://stackoverflow.com/questions/38549972/why-elf-executables-have-a-fixed-load-address
	// This seems to be a convention set in the x86_64 system-v abi: https://refspecs.linuxfoundation.org/elf/x86_64-SysV-psABI.pdf P26
	o.WriteValue(8, virtualStartAddress+textOffset)

	o.WriteBytes(0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Offset from file to program header
	o.WriteBytes(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Start of section header table
	o.WriteBytes(0x00, 0x00, 0x00, 0x00)                         // Flags
	o.WriteBytes(0x40, 0x00)                                     // Size of this header
	o.WriteBytes(0x38, 0x00)                                     // Size of a program header table entry - This should always be the same for 64-bit
	o.WriteBytes(0x02, 0x00)                                     // Length of sections: data and text for now
	o.WriteBytes(0x00, 0x00)                                     // Size of section header, which we aren't using
	o.WriteBytes(0x00, 0x00)                                     // Number of entries section header
	o.WriteBytes(0x00, 0x00)                                     // Index of section header table entry

	// Build Program Header
	// Text Segment
	o.WriteBytes(0x01, 0x00, 0x00, 0x00) // PT_LOAD, loadable segment. Both data and text segment use this.
	o.WriteBytes(0x05, 0x00, 0x00, 0x00) // Flags: 0x4 executable, 0x2 write, 0x1 read
	o.WriteValue(8, 0)                   // textOffset)          // Offset from the beginning of the file. These values depend on how big the header and segment sizes are.
	o.WriteValue(8, virtualStartAddress)
	o.WriteValue(8, virtualStartAddress) // Physical address, irrelavnt on linux.
	o.WriteValue(8, textSize)            // Number of bytes in file image of segment, must be larger than or equal to the size of payload in segment. Should be zero for bss data.
	o.WriteValue(8, textSize)            // Number of bytes in memory image of segment, is not always same size as file image.
	o.WriteValue(8, alignment)

	dataSize := uint64(len(dataSection))
	dataOffset := uint64(textOffset + textSize)
	dataVirtualAddress := dataVirtualStartAddress + dataOffset

	// Build Program Header
	// Data Segment
	o.WriteBytes(0x01, 0x00, 0x00, 0x00) // PT_LOAD, loadable segment. Both data and text segment use this.
	o.WriteBytes(0x07, 0x00, 0x00, 0x00) // Flags: 0x4 executable, 0x2 write, 0x1 read
	o.WriteValue(8, dataOffset)          // Offset address.
	o.WriteValue(8, dataVirtualAddress)  // Virtual address.
	o.WriteValue(8, dataVirtualAddress)  // Physical address.
	o.WriteValue(8, dataSize)            // Number of bytes in file image.
	o.WriteValue(8, dataSize)            // Number of bytes in memory image.
	o.WriteValue(8, alignment)

	// Output the text segment
	o.WriteBytes(textSection...)
	// Output the data segment
	o.WriteBytes(dataSection...)
	return o.o
}

var (
	outputBinaryName = flag.String("output", "tiny-x64", "output binary executable name")
	wordToOutput     = flag.String("word", "Hello World, this is my tiny executable", "word to output in binary sys_write")
)

func main() {
	flag.Parse()

	// data section with word in it
	dataSection := []byte(*wordToOutput)
	wordLen := byte(len(*wordToOutput)) // TODO: Length must be able to fit into a single byte at the moment.

	// https://defuse.ca/online-x86-assembler.htm#disassembly
	textSection := []byte{
		// Sys write
		0x48, 0xC7, 0xC0, 0x04, 0x00, 0x00, 0x00, // mov rax, 0x04
		0x48, 0xC7, 0xC3, 0x01, 0x00, 0x00, 0x00, // mov rbx, 0x01
		0x48, 0xC7, 0xC2, wordLen, 0x00, 0x00, 0x00, // mov rdx, <wordLen>

		0x48, 0xC7, 0xC1, 0xDA, 0x00, 0x60, 0x00, // mov rdx, 0x6000da (HARD CODED at the moment)

		0xcd, 0x80, // int 0x80

		// Sys exit
		0xb8, 0x01, 0x00, 0x00, 0x00, // mov rax, 0x1
		0xbb, 0x00, 0x00, 0x00, 0x00, // mov rbx, 0x0
		0xcd, 0x80, // int 0x80
	}

	data := buildELF(textSection, dataSection)
	if err := ioutil.WriteFile(*outputBinaryName, data, 0755); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote binary to %s\n", *outputBinaryName)
}
