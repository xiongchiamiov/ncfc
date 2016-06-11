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

		fmt.Printf("%s\t%d\n", directory.Name(), total)
	}
	fmt.Println(strings.Repeat("-", 80))
	for _, nonDirectory := range nonDirectories {
		fmt.Println(nonDirectory.Name())
	}
}
