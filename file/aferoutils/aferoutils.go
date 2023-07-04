package aferoutils

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

// Tree displays the directory tree structure starting from the
// specified directory path in a format similar to the `tree` command.
//
// **Parameters:**
//
// fs: The afero.Fs representing the file system to use.
// dirPath: The path of the directory to display the tree structure for.
// prefix: The prefix string to use for each line of the tree structure.
// indent: The indent string to use for each level of the tree structure.
// out: The io.Writer to write the tree structure output to.
//
// **Returns:**
//
// error: An error if any issue occurs while trying to display the tree structure.
func Tree(fs afero.Fs, dirPath, prefix, indent string, out io.Writer) error {
	dirName := filepath.Base(dirPath)
	fmt.Fprintln(out, dirName)
	return printEntry(fs, dirPath, prefix, indent, out, true)
}

func printEntry(fs afero.Fs, dirPath, prefix, indent string, out io.Writer, isLast bool) error {
	entries, err := afero.ReadDir(fs, dirPath)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		isEntryLast := i == len(entries)-1
		entryName := entry.Name()

		if entry.IsDir() {
			subDirPath := filepath.Join(dirPath, entryName)
			if err := printDir(fs, subDirPath, prefix, indent, out, isEntryLast); err != nil {
				return err
			}
		} else {
			printFile(entryName, prefix, indent, out, isEntryLast)
		}
	}

	return nil
}

func printDir(fs afero.Fs, dirPath, prefix, indent string, out io.Writer, isLast bool) error {
	dirName := filepath.Base(dirPath)

	// Print entry
	fmt.Fprint(out, prefix)
	if isLast {
		fmt.Fprint(out, "└── ")
	} else {
		fmt.Fprint(out, "├── ")
	}
	fmt.Fprintln(out, dirName)

	if isLast {
		prefix += "    " // Increase the spacing for the last entry
	} else {
		prefix += "│   " // Increase the spacing for non-last entries
	}

	if err := printEntry(fs, dirPath, prefix, indent, out, false); err != nil {
		return err
	}

	return nil
}

func printFile(fileName, prefix, indent string, out io.Writer, isLast bool) {
	fmt.Fprint(out, prefix)
	if isLast {
		fmt.Fprint(out, "└── ")
	} else {
		fmt.Fprint(out, "├── ")
	}
	fmt.Fprintln(out, fileName)
}
