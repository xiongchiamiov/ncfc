package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func gather_directory_names(directories []os.FileInfo) []string {
	names := []string{}

	for _, directory := range directories {
		find := exec.Command("find", directory.Name())
		wc := exec.Command("wc", "-l")

		reader, writer := io.Pipe()
		find.Stdout = writer
		wc.Stdin = reader

		var buffer bytes.Buffer
		wc.Stdout = &buffer

		find.Start()
		wc.Start()
		find.Wait()
		writer.Close()
		wc.Wait()

		total, _ := strconv.Atoi(strings.TrimSpace(buffer.String()))

		names = append(names, fmt.Sprintf("%s\t%d", directory.Name(), total))
	}

	return names
}

func gather_non_directory_names(nonDirectories []os.FileInfo) []string {
	names := []string{}

	for _, nonDirectory := range nonDirectories {
		names = append(names, nonDirectory.Name())
	}

	return names
}

func main() {
	// If this starts to become a performance problem when inspecting
	// directories with large numbers of direct children, look into using
	// os.File.Readdir instead to avoid the needless sort.
	files, _ := ioutil.ReadDir(".")

	directories := []os.FileInfo{}
	nonDirectories := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file)
		} else {
			nonDirectories = append(nonDirectories, file)
		}
	}

	directoryNames := gather_directory_names(directories)
	for _, directoryName := range directoryNames {
		fmt.Println(directoryName)
	}
	fmt.Println(strings.Repeat("-", 80))
	nonDirectoryNames := gather_non_directory_names(nonDirectories)
	for _, nonDirectoryName := range nonDirectoryNames {
		fmt.Println(nonDirectoryName)
	}
}
