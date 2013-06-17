package main

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
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

// const (
// 	zipPath         = "test/quattroruote/"
// 	imagesPath      = "test/images/"
// 	imagesExtension = ".jpg"
// 	tableName       = "cars"
// )

// Config object
type Config struct {
	ZipPath                       string
	ImagesPath, ImagesExtension   string
	DatabaseConnection, TableName string
}

var config = Config{
	ZipPath:            "test/quattroruote/",
	ImagesPath:         "test/images/",
	ImagesExtension:    ".jpg",
	DatabaseConnection: "root:root@/unioncars",
	TableName:          "cars",
}

// os.FileInfo sorter by ModTime
type byModTime []os.FileInfo

func (i byModTime) Len() int {
	return len(i)
}

func (i byModTime) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i byModTime) Less(a, b int) bool {
	return i[a].ModTime().Before(i[b].ModTime())
}

// Car row in database
type Car struct {
	Id    int
	Brand string
}

func (car Car) get(id int) error {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ? WHERE id=?", config.TableName, id)
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(
			&car.Id,
			&car.Brand,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Cars collection in database
type Cars []Car

func (c *Cars) loadAll() error {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM ?", config.TableName)
	if err != nil {
		return err
	}

	cars := make(Cars, 0, 10)
	for rows.Next() {
		car := Car{}
		err = rows.Scan(
			&car.Id,
			&car.Brand,
		)
		if err != nil {
			return err
		}
		cars = append(cars, car)
	}

	c = &cars
	return nil
}

func main() {
	loadConfig(&config, "config.json")
	updateDatabase()
}

// Load a json configuration file structured as a Config struct
func loadConfig(c *Config, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err = dec.Decode(&c); err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Search the most recently added zip file in config.ZipPath,
// extracts any image (keeping the path) with config.ImagesExtension in config.ImagesPath
// then read CSV files and replaces database content with the CSV content.
func updateDatabase() {
	// Search for zip file
	infos, err := ioutil.ReadDir(config.ZipPath)
	if err != nil {
		fmt.Println("Error reading dir: " + err.Error())
		return
	}
	sort.Sort(sort.Reverse(byModTime(infos)))
	var info os.FileInfo
	for _, info = range infos {
		if strings.HasSuffix(info.Name(), ".zip") {
			break
		}
	}

	// Unzip archive
	archive, err := zip.OpenReader(config.ZipPath + info.Name())
	if err != nil {
		fmt.Println("Error reading zip file [" + config.ZipPath + info.Name() + "]: " + err.Error())
		return
	}
	defer archive.Close()

	// Open database
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		fmt.Println("Error DB opening: " + err.Error())
		return
	}
	defer db.Close()

	// Drop old table
	_, err = db.Exec("TRUNCATE TABLE " + config.TableName)
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
		} else if strings.HasSuffix(f.Name, config.ImagesExtension) {

			// Create direcotry
			dir := filepath.Dir(f.Name)
			if len(dir) > 0 {
				err = os.MkdirAll(config.ImagesPath+filepath.Dir(f.Name), os.FileMode(0777))
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
			dest, err := os.Create(config.ImagesPath + f.Name)
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
