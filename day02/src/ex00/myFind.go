package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func walkFunc(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

func symlinkWalkFunc(path string, f os.FileInfo, err error) error {
	if f.Mode()&os.ModeSymlink != 0 {
		resolvedPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			fmt.Println(path, "-> [broken]")
		} else {
			fmt.Println(path, "->", resolvedPath)
		}
	}
	return nil
}

func dirWalkFunc(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		fmt.Println(path)
	}
	return nil
}

func filesWalkFunc(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		fmt.Println(path)
	}
	return nil
}

func extWalkFunc(ext string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(path, ext) {
			fmt.Println(path)
		}
		return nil
	}
}

func combinedWalkFunc(symlinks, dirs, files bool, ext string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}

		isDir := f.IsDir()
		isFile := !isDir
		isSymlink := f.Mode()&os.ModeSymlink != 0
		matchesExt := strings.HasSuffix(path, ext)

		if isSymlink && symlinks {
			resolvedPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				fmt.Println(path, "-> [broken]")
			} else {
				fmt.Println(path, "->", resolvedPath)
			}
		} else if isDir && dirs {
			fmt.Println(path)
		} else if isFile && files && (ext == "" || matchesExt) {
			fmt.Println(path)
		}
		return nil
	}
}

func main() {
	symlinksFlag := flag.Bool("sl", false, "Set sl flag to print only symbolic links.")
	dirFlag := flag.Bool("d", false, "Set d flag to print only directories.")
	filesFlag := flag.Bool("f", false, "Set f flag to print only files.")
	extName := flag.String("ext", "", "Set ext flag to print files with only specified extension (only with f flag).")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Path is required as the last argument.")
	}
	pathToFind := flag.Arg(flag.NArg() - 1)

	if *extName != "" && !*filesFlag {
		log.Fatal("The ext flag should be used only with the f flag.")
	}

	err := filepath.Walk(pathToFind, combinedWalkFunc(*symlinksFlag, *dirFlag, *filesFlag, *extName))
	if err != nil {
		log.Fatal("Error accessing path: ", err)
	}
}
