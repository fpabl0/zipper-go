package zipper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip uncompress a .zip file.
// Returns
// 		* a rootFolder (if there is not a root folder, root folder is "")
//		* a slice of filenames and foldernames of the zip a file
//		* an error if there is one
func Unzip(src string, dest string) ([]string, []string, error) {

	var rootFolders []string
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return rootFolders, filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return rootFolders, filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		if strings.HasPrefix(fpath, filepath.Join(dest, "__MACOSX")) {
			continue // skip __MACOSX folder
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			if strings.Count(f.Name, "/") == 1 {
				rootFolders = append(rootFolders, fpath) // store root level folders
			}
			os.MkdirAll(fpath, os.ModePerm) // Make Folder
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return rootFolders, filenames, err
		}

		// open zip file
		rc, err := f.Open()
		if err != nil {
			return rootFolders, filenames, err
		}

		// if it is Symbolic Link
		if f.Mode()&os.ModeSymlink != 0 {
			var s strings.Builder
			_, err := io.Copy(&s, rc)
			rc.Close()
			if err != nil {
				return rootFolders, filenames, err
			}
			err = os.Symlink(s.String(), fpath)
			if err != nil {
				return rootFolders, filenames, err
			}
			continue
		}

		// Normal file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return rootFolders, filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close
		// before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return rootFolders, filenames, err
		}
	}
	return rootFolders, filenames, nil
}
