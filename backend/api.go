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
	"log"
	"net/http"
	"net/http/cgi"
	"net/smtp"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"database/sql"
	_ "github.com/EnthusiasticCode/mysql"
)

const (
	logFileName  = "log.txt"
	lockFileName = "update.lock"
)

// Config object
type Config struct {
	ZipPath                          string
	CsvComma                         string
	ImagesPath, ImagesExtension      string
	DatabaseConnection, TableName    string
	SMTPHost, SMTPUser, SMTPPassword string
	MailRecipient                    string
}

var (
	config = Config{
		ZipPath:            "test/quattroruote/",
		CsvComma:           ";",
		ImagesPath:         "test/images/",
		ImagesExtension:    ".jpg",
		DatabaseConnection: "root:root@/unioncars",
		TableName:          "cars",
		SMTPHost:           "mail.example.com",
		SMTPUser:           "user",
		SMTPPassword:       "pass",
		MailRecipient:      "info@example.com",
	}
	logger *log.Logger
)

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
	// Open the global logger
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return // cannot open log file, abort
	}
	defer logFile.Close()
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Load global configuration
	if err = loadConfig(&config, "config.json"); err != nil {
		logger.Fatalln(err)
	}

	err = cgi.Serve(http.StripPrefix("/cgi-bin/api/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		p := path.Clean(r.URL.Path)

		// POST api/mail
		if m, err := path.Match("mail", p); m && err == nil {
			err = r.ParseForm()
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				logger.Fatalln(err)
			}
			// see https://code.google.com/p/go-wiki/wiki/SendingMail
			// Set up authentication information.
			auth := smtp.PlainAuth(
				"",
				config.SMTPUser,
				config.SMTPPassword,
				config.SMTPHost,
			)
			// Connect to the server, authenticate, set the sender and recipient,
			// and send the email all in one step.
			err := smtp.SendMail(
				config.SMTPHost+":25",
				auth,
				r.PostForm.Get("sender"),
				[]string{config.MailRecipient},
				[]byte(r.PostForm.Get("text")),
			)
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				logger.Fatal(err)
			}
			io.WriteString(w, "{\"status\": \"ok\", \"sender\": \""+r.PostForm.Get("sender")+"\"}")
			return
		}

		// Update database if needed
		go updateIfNeeded()

		encoder := json.NewEncoder(w)

		// GET api/car
		if m, err := path.Match("car", p); m && err == nil {
			c := make(Cars, 0, 10)
			if err := c.loadAll(); err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				logger.Fatalln(err)
			}
			if err = encoder.Encode(c); err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				logger.Fatalln(err)
			}
			return
		}

		// GET api/car/<id>
		if m, err := path.Match("car/*", p); m && err == nil {
			id, err := strconv.Atoi(path.Base(p))
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				logger.Fatalln(err)
			}
			var c Car
			err = c.get(id)
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				logger.Fatalln(err)
			}
			err = encoder.Encode(c)
			if err != nil {
				io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
				logger.Fatalln(err)
			}
			return
		}
	})))
	// On error run locally
	if err != nil {
		updateIfNeeded()
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

func updateIfNeeded() {
	// Managin lock via file
	lockFile, err := os.OpenFile(lockFileName, os.O_CREATE|os.O_EXCL, 0660)
	if err != nil {
		return // someone else is already updating the database
	}
	defer func() {
		lockFile.Close()
		os.Remove(lockFileName)
	}()

	// Updating database
	if err = updateDatabase(); err != nil {
		logger.Println(err)
		return
	}
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
	var info os.FileInfo = nil
	for _, i := range infos {
		if strings.HasSuffix(i.Name(), ".zip") && !strings.HasPrefix(i.Name(), ".") {
			info = i
			break
		}
	}

	// Early exit if no zip file found
	if info == nil {
		return nil
	}

	// Unzip archive
	archive, err := zip.OpenReader(config.ZipPath + info.Name())
	if err != nil {
		return err
	}
	defer archive.Close()
	defer func() {
		if err != nil {
			return // do not delete anything if import failed
		}
		// Delete all files listed earlier
		for _, info = range infos {
			os.Remove(config.ZipPath + info.Name())
		}
	}()

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
			c.FieldsPerRecord = -1
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
		} else if strings.HasSuffix(strings.ToLower(f.Name), config.ImagesExtension) {

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
