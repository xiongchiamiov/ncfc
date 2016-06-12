package main

import (
	"bytes"
	"fmt"
	"github.com/rthornton128/goncurses"
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

		names = append(names, fmt.Sprintf("%s -- %d", directory.Name(), total))
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
	stdscr, _ := goncurses.Init()
	defer goncurses.End()

	goncurses.Raw(true)
	goncurses.Echo(false)
	goncurses.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

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

	menuItems := []*goncurses.MenuItem{}
	directoryNames := gather_directory_names(directories)
	for _, directoryName := range directoryNames {
		item, err := goncurses.NewItem(directoryName, "")
		if err != nil {
			stdscr.Print(err)
		}
		menuItems = append(menuItems, item)
		defer item.Free()
	}
	spacer, _ := goncurses.NewItem(strings.Repeat("-", 80), "")
	menuItems = append(menuItems, spacer)
	nonDirectoryNames := gather_non_directory_names(nonDirectories)
	for _, nonDirectoryName := range nonDirectoryNames {
		item, err := goncurses.NewItem(nonDirectoryName, "")
		if err != nil {
			stdscr.Print(err)
		}
		menuItems = append(menuItems, item)
		defer item.Free()
	}

	menu, err := goncurses.NewMenu(menuItems)
	if err != nil {
		stdscr.Print(err)
		return
	}
	defer menu.Free()
	menu.Post()

	stdscr.MovePrint(20, 0, "'q' to exit")
	stdscr.Refresh()

	for {
		goncurses.Update()
		ch := stdscr.GetChar()

		switch goncurses.KeyString(ch) {
		case "q":
			return
		case "down", "j":
			menu.Driver(goncurses.REQ_DOWN)
		case "up", "k":
			menu.Driver(goncurses.REQ_UP)
		}
	}
}
