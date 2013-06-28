package main

// To build for server CGI with fish shell use:
// begin; set -lx GOOS linux; set -lx GOARCH 386; go build api.go; end

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
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

	"github.com/bjarneh/latinx"

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
	SMTPHost, SMTPUser, SMTPPassword string
	MailRecipient                    string
	DatabaseConnection, TableName    string
	TableMapping                     []ConfigColumnAlias
}

type ConfigColumnAlias struct {
	TableColumn string
	Alias       string
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
		TableMapping:       nil,
	}
	logger *log.Logger
)

// CSV and Database elements
type csvElement map[string]string
type dbElement map[string]interface{}

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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Status", "200 OK")

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
			cars, err := getCars()
			if err != nil {
				outputError(w, err)
			}

			if err = encoder.Encode(cars); err != nil {
				outputError(w, err)
			}
			return
		}

		// GET api/car/<id>
		if m, err := path.Match("car/*", p); m && err == nil {
			id, err := strconv.Atoi(path.Base(p))
			if err != nil {
				outputError(w, err)
			}

			car, err := getCar(id)
			if err != nil {
				outputError(w, err)
			}

			err = encoder.Encode(car)
			if err != nil {
				outputError(w, err)
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

// Load
func mapQueryResults(rows *sql.Rows) ([]dbElement, error) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Create result array
	result := make([]dbElement, 0)

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		res := make(dbElement)
		for i, col := range values {
			if col == nil {
				res[columns[i]] = nil
			} else {
				// TODO parse time.Time
				res[columns[i]] = string(col)
			}
		}

		// Add to result
		result = append(result, res)
	}

	return result, nil
}

func getCars() ([]dbElement, error) {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + config.TableName)
	if err != nil {
		return nil, err
	}

	mapRows, err := mapQueryResults(rows)
	if err != nil {
		return nil, err
	}

	return mapRows, nil
}

func getCar(id int) (dbElement, error) {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM "+config.TableName+" WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	mapRows, err := mapQueryResults(rows)
	if err != nil {
		return nil, err
	}

	if len(mapRows) == 0 {
		return nil, nil
	}

	return mapRows[0], nil
}

// Lock the archives folder and updates the database
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

			// Reader from ISO-8859-1
			ffutf8 := latinx.NewReader(latinx.ISO_8859_1, ff)

			// Open CSV reader
			c := csv.NewReader(ffutf8)
			c.Comma = rune(config.CsvComma[0])
			c.FieldsPerRecord = -1
			c.TrimLeadingSpace = false
			rs, err := c.ReadAll()
			if err != nil {
				return err
			}

			// Convert csv to slice of csvElement
			elements := make([]csvElement, 0, 50)
			var columns []string
			for i, r := range rs {
				if i == 0 {
					// Columns
					columns = r
				} else {
					// Content
					element := make(csvElement)
					for c, cv := range r {
						element[columns[c]] = cv
					}
					elements = append(elements, element)
				}
			}
			ff.Close()

			// Build insert prepared query
			columnsCount := len(config.TableMapping)
			insertQueryString := "INSERT INTO " + config.TableName + " ("
			for i, column := range config.TableMapping {
				insertQueryString += column.TableColumn
				if i+1 < columnsCount {
					insertQueryString += ","
				}
			}
			insertQueryString += ") VALUES ("
			for i := 0; i < columnsCount; i++ {
				insertQueryString += "?"
				if i+1 < columnsCount {
					insertQueryString += ","
				}
			}
			insertQueryString += ")"
			insertQuery, err := db.Prepare(insertQueryString)
			if err != nil {
				return err
			}
			// Put csvElements in database
			for _, element := range elements {
				values := make([]interface{}, 0, columnsCount)
				for _, column := range config.TableMapping {
					values = append(values, element[column.Alias])
				}
				_, err = insertQuery.Exec(values...)
				if err != nil {
					return err
				}
			}
			insertQuery.Close()
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

func outputError(w http.ResponseWriter, err error) {
	w.Header().Set("Status", "500 Internal Server Error")
	io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
	logger.Fatalln(err)
}
