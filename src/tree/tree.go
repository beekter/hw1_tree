package tree

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type fileSequence []os.FileInfo

func (fs fileSequence) Len() int {
	return len(fs)
}

func (fs fileSequence) Less(i, j int) bool {
	return fs[i].Name() < fs[j].Name()
}

func (fs fileSequence) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs fileSequence) dirs() fileSequence {
	var files fileSequence

	for _, file := range fs {
		if file.IsDir() {
			files = append(files, file)
		}
	}

	return files
}

func Print(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, "")
}

func printDir(out io.Writer, path string, printFiles bool, parentTab string) error {
	dir, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}

	var files fileSequence

	files, err = dir.Readdir(0)
	//сортировка по имени
	sort.Sort(files)

	//если файлы не нужны, то отфильтровываем только папки
	if !printFiles && files != nil {
		files = files.dirs()
	}

	for i := range files {
		var file = files[i]
		var lastNode = i == len(files)-1
		var childTab = parentTab + "│\t"

		if lastNode {
			childTab = parentTab + "\t"
		}

		if err = printNode(out, parentTab, lastNode, file); err != nil {
			return err
		}
		if file.IsDir() {
			childPath := path + string(os.PathSeparator) + file.Name()
			if err = printDir(out, childPath, printFiles, childTab); err != nil {
				return err
			}
		}
	}

	return dir.Close()
}

func printNode(out io.Writer, parentTab string, lastNode bool, info os.FileInfo) error {
	var sizeString string

	branchSign := "├"
	if lastNode {
		branchSign = "└"
	}

	if !info.IsDir() {
		bytes := info.Size()
		if bytes == 0 {
			sizeString = " (empty)"
		} else {
			sizeString = " (" + strconv.FormatInt(bytes, 10) + "b)"
		}
	}

	if _, err := fmt.Fprintf(out, "%s%s───%s%s\n", parentTab, branchSign, info.Name(), sizeString); err != nil {
		return err
	}
	return nil
}
