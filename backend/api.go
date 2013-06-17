package main

// To build for server CGI with fish shell use:
// begin; set -lx GOOS linux; set -lx GOARCH 386; go build api.go; end

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cgi"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"database/sql"
	_ "github.com/EnthusiasticCode/mysql"
)

// Config object
type Config struct {
	ZipPath                       string
	CsvComma                      string
	ImagesPath, ImagesExtension   string
	DatabaseConnection, TableName string
}

var config = Config{
	ZipPath:            "test/quattroruote/",
	CsvComma:           ";",
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

	fmt.Println(rune("\t"[0]))

	err := cgi.Serve(http.StripPrefix("/cgi-bin/api/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Update database if needed
		err := loadConfig(&config, "config.json")
		if err != nil {
			io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
			return
		}
		err = updateDatabase()
		if err != nil {
			fmt.Println(err.Error())
		}

		p := path.Clean(r.URL.Path)
		encoder := json.NewEncoder(w)

		// api/car
		if m, err := path.Match("car", p); m && err == nil {
			c := make(Cars, 0, 10)
			err := c.loadAll()
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				return
			}
			err = encoder.Encode(c)
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				return
			}
			return
		}

		// api/car/<id>
		if m, err := path.Match("car/*", p); m && err == nil {
			id, err := strconv.Atoi(path.Base(p))
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				return
			}
			var c Car
			err = c.get(id)
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				return
			}
			err = encoder.Encode(c)
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				return
			}
			return
		}
	})))
	//
	if err != nil {
		fmt.Println(err.Error())
	}
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
func updateDatabase() error {
	// Search for zip file
	infos, err := ioutil.ReadDir(config.ZipPath)
	if err != nil {
		return err
	}

	// Early return if no file present
	if len(infos) == 0 {
		return nil
	}

	// Select most recent valid archive infos
	sort.Sort(sort.Reverse(byModTime(infos)))
	var info os.FileInfo
	for _, info = range infos {
		if strings.HasSuffix(info.Name(), ".zip") && !strings.HasPrefix(info.Name(), ".") {
			break
		}
	}

	// Unzip archive
	archive, err := zip.OpenReader(config.ZipPath + info.Name())
	if err != nil {
		return err
	}
	defer archive.Close()
	defer os.Remove(config.ZipPath + info.Name())

	// Open database
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	// Drop old table
	_, err = db.Exec("TRUNCATE TABLE " + config.TableName)
	if err != nil {
		return err
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
				return err
			}

			// Open CSV reader
			c := csv.NewReader(ff)
			c.Comma = rune(config.CsvComma[0])
			c.TrimLeadingSpace = true
			rs, err := c.ReadAll()
			if err != nil {
				return err
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
					return err
				}
			}

			// Extract image
			ff, err := f.Open()
			if err != nil {
				return err
			}
			dest, err := os.Create(config.ImagesPath + f.Name)
			if err != nil {
				return err
			}
			_, err = io.Copy(dest, ff)
			if err != nil {
				return err
			}
			ff.Close()
		}
	}

	return nil
}
