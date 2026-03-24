package goappimagetool

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/yalue/elf_reader"
)

func readELF(path string) (elf_reader.ELFFile, error) {
	// Print the section names in path. This code will work on both 32-bit
	// and 64-bit systems.
	raw, e := os.ReadFile(path)
	if e != nil {
		fmt.Printf("Failed reading %s: %s\n", "AppImage", e)
		return nil, e
	}
	elf, e := elf_reader.ParseELFFile(raw)
	if e != nil {
		fmt.Printf("Failed parsing ELF file: %s\n", e)
		return nil, e
	}

	return elf, nil
}

func getSectionHeaderByName(path string, section string) (elf_reader.ELFSectionHeader, error) {
	elf, err := readELF(path)
	if err != nil {
		return nil, err
	}

	var count uint16 = elf.GetSectionCount()
	for i := range count {
		if i == 0 {
			continue
		}

		name, err := elf.GetSectionName(i)
		Check(err)

		if name == section {
			header, err := elf.GetSectionHeader(i)
			Check(err)
			return header, nil
		}
	}

	return nil, fmt.Errorf("Section not found")
}

func hashEngine(path string) (hash.Hash, int64) {
	var offset int64 = 0
	hash := md5.New()

	elf, err := readELF(path)
	Check(err)

	var count uint16 = elf.GetSectionCount()
	for i := range count {
		// First section never has a name.
		if i == 0 {
			continue
		}

		name, err := elf.GetSectionName(i)
		Check(err)
		switch name {
		case ".bss":
			// This only exists in memory, and
			// does not have contents which need to be hashed.
			continue
		case ".digest_md5", ".sha256_sig", ".sig_key":
			// Skip them entirely as if they are not even here.
			continue
		}

		// fmt.Printf("Section %d name: %s\n", i, name)
		content, err := elf.GetSectionContent(i)
		Check(err)
		bytesWritten, err := hash.Write(content)

		offset += int64(bytesWritten)
	}

	return hash, offset
}

func CalculateMD5(path string) []byte {
	// First read the start of the file and hash its contents.
	hash, offset := hashEngine(path)

	// Now read the rest of the file and hash its contents.
	file, err := os.Open(path)
	Check(err)

	fileOffset, err := file.Seek(offset, io.SeekStart)
	Check(err)

	if fileOffset != int64(offset) {
		Check(fmt.Errorf("No squashfs was appended"))
	}

	var buf []byte = make([]byte, 4096)
	for {
		size, err := file.Read(buf)
		if size > 0 {
			// Write only as many bytes as i need (don't
			// also hash padding on last iteration)
			hash.Write(buf[:size])
		}
		if err != nil {
			break
		}
	}
	finalHash := hash.Sum(nil)

	return finalHash
}

func UpdateMD5(path string, hash []byte) error {
	return OverwriteSection(path, ".digest_md5", hash)
}

func CalculateSha256(path string) []byte {
	hash := sha256.New()

	file, err := os.Open(path)
	Check(err)

	var buf []byte = make([]byte, 4096)
	for {
		size, err := file.Read(buf)
		if size > 0 {
			// Write only as many bytes as i need (don't
			// also hash padding on last iteration)
			hash.Write(buf[:size])
		}
		if err != nil {
			break
		}
	}

	return hash.Sum(nil)
}

func UpdateSha256(path string, hash []byte) error {
	return OverwriteSection(path, ".sha256_sig", hash)
}
