package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

var verbose *bool = flag.Bool("verbose", false, "Print the list of duplicate files.")
var rootDir string = "."
var fullPathsByFilename map[string][]string

func Visit(fullpath string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if err != nil {
		fmt.Println(err)
	}
	filename := path.Base(fullpath)
	fullPathsByFilename[filename] = append(fullPathsByFilename[filename], fullpath)
	return nil
}

func MD5OfFile(fullpath string) []byte {
	fi, err := os.Open(fullpath)
	if err != nil {
		return nil
	}
	defer fi.Close()

	r := bufio.NewReader(fi)

	buf := make([]byte, 1024)
	md5sum := md5.New()
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return nil
		}
		if n == 0 {
			break
		}

		md5sum.Write(buf[:n])
	}

	return md5sum.Sum(nil)
}

func PrintResults() {
	dupes := 0
	for key, value := range fullPathsByFilename {
		if len(value) < 2 {
			continue
		}
		dupes++
		if *verbose {
			println(key, ":")
			for _, filename := range value {
				println("  ", filename)
				fmt.Printf("    %x\n", MD5OfFile(filename))
			}
		}
	}
	println("Total duped files found:", dupes)
}

func FindDupes(root string) {
	fullPathsByFilename = make(map[string][]string)
	filepath.Walk(root, Visit)
}

func ParseArgs() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		rootDir = flag.Arg(0)
	}
}

func main() {
	ParseArgs()
	FindDupes(rootDir)
	PrintResults()
}
