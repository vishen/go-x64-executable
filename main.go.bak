package main

type Builder struct {
	// o bytes.Buffer
	o []byte
}

func (b *Builer) Len() int {
	return len(b.o)
}

func (b *Builder) WriteBytes(bs ...byte) {
	// b.o.Write(bs)
	b.o = append(b.o, bs...)
}

func buildELF(dataSection, textSection []byte) {
	var o Builder

	// Build ELF Header
	// TODO: Do these need to be changed to little endian?
	o.WriteBytes(0x7f, 0x45, 0x4c, 0x46) // ELF magic value

	o.WriteBytes(0x02) // 64-bit executable
	o.WriteBytes(0x01) // Little endian
	o.WriteBytes(0x01) // ELF version
	o.WriteBytes(0x00) // Target OS ABI
	o.WriteBytes(0x00) // Further specify ABI version

	o.WriteBytes(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Unused bytes

	o.WriteBytes(0x02, 0x00) // Executable type
	o.WriteBytes(0x3e, 0x00) // x86-64 target architecture
	o.WriteBytes(0x01)       // ELF version

	// TODO: 8 bytes for memory address of entry point, this
	// needs to be the virtual address. 
	o.WriteBytes(0x??,0x??, 0x??,0x??,0x??,0x??,0x??,0x??)	// TODO: ?? Should also be the virtual starting address for text segment
	o.WriteBytes(0x40) // Offset from file to program header
	o.WriteBytes(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Start of section header table
	o.WriteBytes(0x00, 0x00, 0x00, 0x00) // Flags
	o.WriteBytes(0x40, 0x00)// Size of this header
	o.WriteBytes(0x38, 0x00) // Size of a program header table entry - This should always be the same for 64-bit 
	o.WriteBytes(0x02, 0x00) // Length of sections: data and text for now
	o.WriteBytes(0x00, 0x00) // Size of section header, which we aren't using
	o.WriteBytes(0x00, 0x00) // Number of entries section header
	o.WriteBytes(0x00, 0x00) // Index of section header table entry
	

	// 64-bit virtual offsets always start at 0x400000?? https://stackoverflow.com/questions/38549972/why-elf-executables-have-a-fixed-load-address
	// This seems to be a convention set in the x86_64 system-v abi: https://refspecs.linuxfoundation.org/elf/x86_64-SysV-psABI.pdf P26
	VIRTUAL_OFFSET := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00}
	// ALTHOUGH, protoc and go binaries are starting at 0x00 00 00 00 00 40 00 40??

	// Build Program Header
	// Text Segment
	o.WriteBytes(0x01, 0x00, 0x00, 0x00) // PT_LOAD, loadable segment. Both data and text segment use this.
	o.WriteBytes(0x05, 0x00, 0x00, 0x00) // Flags: 0x4 executable, 0x2 write, 0x1 read
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Offset from the beginning of the file. These values depend on how big the header and segment sizes are, since we don't want to overlap there. TODO: Should be `textSegmentOffset` defined below
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Virtual address
	o.WriteBytes(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Physical address, irrelavnt on linux.
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Number of bytes in file image of segment, must be larger than or equal to the size of payload in segment. TODO: Length of data segment?
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Number of bytes in memory image of segment, is not always same size as file image??
	o.WriteBytes(0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Alginment. I am unsure how to set this, apparently setting to 1 is valid?


	// Build Program Header
	// Data Segment
	o.WriteBytes(0x01, 0x00, 0x00, 0x00) // PT_LOAD, loadable segment. Both data and text segment use this.
	o.WriteBytes(0x07, 0x00, 0x00, 0x00) // Flags: 0x4 executable, 0x2 write, 0x1 read
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Offset from the beginning of the file. These values depend on how big the header and segment sizes are, since we don't want to overlap there.  TODO: Should be `dataSegmentOffset`, defined below.
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Virtual address, TODO: WTF to do here?
	o.WriteBytes(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Physical address, irrelavnt on linux.
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Number of bytes in file image of segment, must be larger than or equal to the size of payload in segment. TODO: Length of data segment?
	o.WriteBytes(0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??, 0x??) // Number of bytes in memory image of segment, is not always same size as file image??
	o.WriteBytes(0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Alginment. I am unsure how to set this, apparently setting to 1 is valid?

	// TODO: Append the text and data segments at the very end, one after the other, with some padding inbetween?

	// ELF Header plus program header for text and program header for data
	headerLength := o.Len()
	
	// TODO: Pad bytes to nearest power of 2??
	// o.WriteBytes(0x00, 0x00, 0x00, 0x00, 0x00) // Pad bytes

	textSegmentOffset := o.Len()
	// Output the text segment
	o.WriteBytes(textSegment...)

	// TODO: Pad bytes for safety?

	dataSegmentOffset := o.Len()
	// Output the data segment
	o.WriteBytes(dataSegment...)
}

func main() {
	// Linux "Hello World" text section, possibly in Little Endian?
	textSection := []byte{
		0xb8, 0x04, 0x00, 0x00, 0x00, // mov rax, 0x04
		0xbb, 0x01, 0x00, 0x00, 0x00, // mov rbx, 0x01
		0x48, 0xb9, 0xd8, 0x00, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, // mov rcx, 0x600d8 (TODO: This needs to be the actual position in the binary...?)
		0xba, 0x0c, 0x00, 0x00, 0x00, // mov edx, 0xc
		0xcd, 0x80, // int 0x80
		0xb8, 0x01, 0x00, 0x00, 0x00, // mov rax, 0x1
		0xbb, 0x00, 0x00, 0x00, 0x00, // mov rbx, 0x0
		0xcd, 0x80, // int 0x80
	}

	// "Hello World" data section
	dataSection := []byte{0x60, 0x00, 0xd8, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x0a}
	fmt.Println(textSection)
	fmt.Println(dataSection)
}
