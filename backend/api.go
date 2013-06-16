package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"database/sql"
	_ "github.com/EnthusiasticCode/mysql"
)

const (
	zipPath         = "test/quattroruote/"
	imagesPath      = "test/images/"
	imagesExtension = ".jpg"
	tableName       = "cars"
)

// os.FileInfo sorter by ModTime
type ByModTime []os.FileInfo

func (i ByModTime) Len() int {
	return len(i)
}

func (i ByModTime) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i ByModTime) Less(a, b int) bool {
	return i[a].ModTime().Before(i[b].ModTime())
}

// type Car struct {
// 	Id    int
// 	Brand string
// }

// type Cars []Car

func main() {
	updateDatabase()
}

func updateDatabase() {
	// Search for zip file
	infos, err := ioutil.ReadDir(zipPath)
	if err != nil {
		fmt.Println("Error reading dir: " + err.Error())
		return
	}
	sort.Sort(sort.Reverse(ByModTime(infos)))
	var info os.FileInfo
	for _, info = range infos {
		if strings.HasSuffix(info.Name(), ".zip") {
			break
		}
	}

	// Unzip archive
	archive, err := zip.OpenReader(zipPath + info.Name())
	if err != nil {
		fmt.Println("Error reading zip file [" + zipPath + info.Name() + "]: " + err.Error())
		return
	}
	defer archive.Close()

	// Open database
	db, err := sql.Open("mysql", "root:root@/unioncars")
	if err != nil {
		fmt.Println("Error DB opening: " + err.Error())
		return
	}
	defer db.Close()

	// Drop old table
	_, err = db.Exec("TRUNCATE TABLE " + tableName)
	if err != nil {
		fmt.Println("Error truncating DB table: " + err.Error())
		return
	}

	// Reading files in archive
	for _, f := range archive.File {
		if strings.HasPrefix(filepath.Base(f.Name), ".") {
			continue
		}

		if strings.HasSuffix(f.Name, ".csv") {

			// Reading CSV file
			ff, err := f.Open()
			if err != nil {
				fmt.Println("Error reading CSV file [" + f.Name + "]: " + err.Error())
				return
			}

			// Open CSV reader
			c := csv.NewReader(ff)
			c.Comma = ';'
			c.TrimLeadingSpace = true
			rs, err := c.ReadAll()
			if err != nil {
				fmt.Println("Error reading CSV content: " + err.Error())
				return
			}

			// Convert csv to sql
			for i, r := range rs {
				if i == 0 {
					// Columns
					fmt.Println(r) // []string
				} else {
					// Content
					fmt.Println(r)
				}
			}
			ff.Close()
		} else if strings.HasSuffix(f.Name, imagesExtension) {

			// Create direcotry
			dir := filepath.Dir(f.Name)
			if len(dir) > 0 {
				err = os.MkdirAll(imagesPath+filepath.Dir(f.Name), os.FileMode(0777))
				if err != nil {
					fmt.Println("Error creating directory for images: " + err.Error())
					return
				}
			}

			// Extract image
			ff, err := f.Open()
			if err != nil {
				fmt.Println("Error reading image file [" + f.Name + "]: " + err.Error())
				return
			}
			dest, err := os.Create(imagesPath + f.Name)
			if err != nil {
				fmt.Println("Error creating new image file [" + f.Name + "]: " + err.Error())
				return
			}
			_, err = io.Copy(dest, ff)
			if err != nil {
				fmt.Println("Error unzipping image file [" + f.Name + "]: " + err.Error())
				return
			}
			ff.Close()
		}
	}
}
