package main

import (
  "crypto/md5"
  "flag"
  "fmt"
  "bufio"
  "os"
  "path"
  "path/filepath"
)

var verbose *bool = flag.Bool("verbose", false, "Print the list of duplicate files.")
var rootDir string = "."
var fullPathsByFilename map[string][]string

type DupeChecker struct{}

func (dc DupeChecker) VisitDir(fullpath string, f *os.FileInfo) bool {
  return true
}

func (dc DupeChecker) VisitFile(fullpath string, f *os.FileInfo) {
  filename := path.Base(fullpath)
  fullPathsByFilename[filename] = append(fullPathsByFilename[filename], fullpath)
}

func MD5OfFile(fullpath string) []byte {
  fi, err := os.Open(fullpath)
  if err != nil { return nil }
  defer fi.Close()

  r := bufio.NewReader(fi)

  buf := make([]byte, 1024)
  md5sum := md5.New()
  for {
    n, err := r.Read(buf)
    if err != nil && err != os.EOF { return nil }
    if n == 0 { break }

    md5sum.Write(buf[:n])
  }

  return md5sum.Sum()
}

func PrintResults() {
  dupes := 0
  for key, value := range fullPathsByFilename {
    if (len(value) < 2) {
      continue
    }
    dupes++
    if (*verbose) {
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
  filepath.Walk(root, DupeChecker{}, nil)
}

func ParseArgs() {
  flag.Parse()
  if (len(flag.Args()) > 0) { 
    rootDir = flag.Arg(0)
  } 
}

func main() {
  ParseArgs()
  FindDupes(rootDir)
  PrintResults()
}