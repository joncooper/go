package main

import (
  "io/ioutil"
  "fmt"
  "os"
  "path"
  "hash"
  "path/filepath"
  "crypto/md5"
)

var rootDir = "."
var fullPathsByFilename map[string][]string

func MD5OfFile(fullpath string) []byte {
  if contents, err := ioutil.ReadFile(fullpath); err == nil { 
    var md5sum hash.Hash = md5.New()
    md5sum.Write(contents)
    return md5sum.Sum()
  }
  return nil
}

type DupeChecker struct{}

func (dc DupeChecker) VisitDir(fullpath string, f *os.FileInfo) bool {
  return true
}

func (dc DupeChecker) VisitFile(fullpath string, f *os.FileInfo) {
  filename := path.Base(fullpath)
  fullPathsByFilename[filename] = append(fullPathsByFilename[filename], fullpath)
}

func PrintResults() {
  dupes := 0
  for key, value := range fullPathsByFilename {
    if (len(value) < 2) {
      continue
    }
    dupes++
    println(key, ":")
    for _, filename := range value {
      println("  ", filename)
      fmt.Printf("    %x\n", MD5OfFile(filename))
    }
  }
  println("Total duped files found:", dupes)
}

func FindDupes(root string) {
  fullPathsByFilename = make(map[string][]string)
  filepath.Walk(root, DupeChecker{}, nil)
}

func main() {
  FindDupes(rootDir)
  PrintResults()
}
