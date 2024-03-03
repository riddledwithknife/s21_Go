package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func rotateLog(logFile, archiveDir string) error {
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		return err
	}

	mtime := fileInfo.ModTime().Unix()
	baseName := filepath.Base(logFile)
	archiveName := fmt.Sprintf("%s_%d.tar.gz", baseName, mtime)

	if archiveDir != "" {
		if err := os.MkdirAll(archiveDir, 0755); err != nil {
			return err
		}
		archiveName = filepath.Join(archiveDir, archiveName)
	}

	tarFile, err := os.Create(archiveName)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return err
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func main() {
	archivePath := flag.String("a", "", "Custom archive directory")
	flag.Parse()
	logFiles := flag.Args()

	if len(logFiles) == 0 {
		log.Fatal("No log files specified.")
	}

	var wg sync.WaitGroup
	for _, logFile := range logFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := rotateLog(file, *archivePath)
			if err != nil {
				fmt.Printf("Error rotating log file %s: %v\n", file, err)
			}
		}(logFile)
	}

	wg.Wait()
}
