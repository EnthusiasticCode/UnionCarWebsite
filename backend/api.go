package main

// To build for server CGI with fish shell use:
// begin; set -lx GOOS linux; set -lx GOARCH 386; go build api.go; end

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
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
	"github.com/gorilla/sessions"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	logFileName  = "log.txt"
	lockFileName = "update.lock"
)

// Config object
type Config struct {
	SessionKey                       string
	SessionName                      string
	ZipPath                          string
	ZipPathKeep                      string
	ZipPrefixes                      []string
	CsvComma                         string
	ImagesPath, ImagesExtension      string
	SMTPHost, SMTPUser, SMTPPassword string
	MailRecipient                    string
	DatabaseConnection, TableName    string
	TableMapping                     []ConfigColumnAlias
	SpecialFields                    []string
	ActivationEmailTemplate          string
	NoticeEmailTemplate              string
	ConfirmationEmailTemplate        string
}

type ConfigColumnAlias struct {
	TableColumn string
	Alias       string
	Transformer string
}

var (
	config = Config{
		SessionKey:                "SecretSessionKey",
		SessionName:               "SessionName",
		ZipPath:                   "test/quattroruote/",
		ZipPathKeep:               "",
		ZipPrefixes:               nil,
		CsvComma:                  ";",
		ImagesPath:                "test/images/",
		ImagesExtension:           ".jpg",
		DatabaseConnection:        "root:root@/unioncars",
		TableName:                 "cars",
		SMTPHost:                  "mail.example.com",
		SMTPUser:                  "user",
		SMTPPassword:              "pass",
		MailRecipient:             "info@example.com",
		TableMapping:              nil,
		SpecialFields:             nil,
		ActivationEmailTemplate:   "activation key: {{.Key}}",
		NoticeEmailTemplate:       "Request sent",
		ConfirmationEmailTemplate: "Account active",
	}
	logger       *log.Logger
	transformers = map[string]func(string) string{
		"SiNo2Boolean": func(v string) string {
			if strings.Contains(v, "SI") {
				return "1"
			}
			return "0"
		},
		"EurDate2SQLDate": func(v string) string {
			parts := strings.Split(v, "/")
			if len(parts) < 3 {
				return v
			}
			return parts[2] + "-" + parts[1] + "-" + parts[0]
		},
		"CleanName": func(v string) string {
			return strings.TrimRight(v, " *")
		},
	}
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
		//w.Header().Set("Connection", "close")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Status", "200 OK")

		p := path.Clean(r.URL.Path)

		// POST api/mail
		if m, err := path.Match("mail", p); m && err == nil {
			err = r.ParseForm()
			if err != nil {
				outputError(w, err)
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
				outputError(w, err)
			}
			io.WriteString(w, "{\"status\": \"ok\", \"sender\": \""+r.PostForm.Get("sender")+"\"}")
			return
		}

		// Get logged user from session if present
		cookieStore := sessions.NewCookieStore([]byte(config.SessionKey))
		session, _ := cookieStore.Get(r, config.SessionName)
		isLoggedInUser := session.Values["loggedIn"] == "special"

		// POST api/login
		if m, err := path.Match("login", p); m && err == nil {
			r.ParseForm()
			user := r.PostForm.Get("email")
			pass := r.PostForm.Get("password")
			authorized, err := isUserEnabled(user, pass)

			if authorized {
				session.Values["loggedIn"] = "special"
				session.Save(r, w)
				io.WriteString(w, "{\"login\": \"ok\"}")
			} else if err != nil {
				io.WriteString(w, "{\"login\": \"invalid\", \"error\": \""+err.Error()+"\"}")
			} else {
				io.WriteString(w, "{\"login\": \"invalid\"}")
			}

			return
		}

		// POST api/register
		if m, err := path.Match("register", p); m && err == nil {
			r.ParseForm()
			email := r.PostForm.Get("email")
			pass := r.PostForm.Get("password")
			if email == "" || pass == "" {
				io.WriteString(w, "{\"registration\": \"error\", \"error\": \"Invalid data\"}")
				return
			}

			key, err := registerUser(email, pass)
			if err != nil {
				io.WriteString(w, "{\"registration\": \"error\", \"error\": \""+err.Error()+"\"}")
				return
			}

			user := struct {
				Name    string
				Address string
				City    string
				Zip     string
				Region  string
				Vat     string
				Email   string
				Key     string
			}{
				Name:    r.PostForm.Get("name"),
				Address: r.PostForm.Get("address"),
				City:    r.PostForm.Get("city"),
				Zip:     r.PostForm.Get("zip"),
				Region:  r.PostForm.Get("Region"),
				Vat:     r.PostForm.Get("vat"),
				Email:   email,
				Key:     key,
			}

			auth := smtp.PlainAuth(
				"",
				config.SMTPUser,
				config.SMTPPassword,
				config.SMTPHost,
			)

			// Templates
			activationTmpl, _ := template.New("activation").Parse(config.ActivationEmailTemplate)
			noticeTmpl, _ := template.New("notice").Parse(config.NoticeEmailTemplate)

			// Activation email
			var activationContent bytes.Buffer
			activationTmpl.Execute(&activationContent, user)
			smtp.SendMail(
				config.SMTPHost+":25",
				auth,
				r.PostForm.Get("email"),
				[]string{config.MailRecipient},
				activationContent.Bytes(),
			)

			// Notice email
			var noticeContent bytes.Buffer
			noticeTmpl.Execute(&noticeContent, user)
			smtp.SendMail(
				config.SMTPHost+":25",
				auth,
				config.MailRecipient,
				[]string{r.PostForm.Get("email")},
				noticeContent.Bytes(),
			)

			io.WriteString(w, "{\"registration\": \"ok\"}")
			return
		}

		// GET /api/activate/<key>
		if m, err := path.Match("activate/*", p); m && err == nil {
			key := path.Base(p)
			activateUserWithKey(key)
			return
		}

		// Update database if needed
		updateIfNeeded()

		encoder := json.NewEncoder(w)

		// GET api/car
		if m, err := path.Match("car", p); m && err == nil {
			cars, err := getCars(isLoggedInUser)
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

			car, err := getCar(id, isLoggedInUser)
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

// Add user to the database and wait activation
func registerUser(email string, password string) (string, error) {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return "", err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM users WHERE email=?", email)
	if err != nil {
		return "", err
	}

	if rows.Next() {
		return "", nil
	}

	hash := md5.New()
	io.WriteString(hash, email)
	io.WriteString(hash, password)
	randToken := make([]byte, 10)
	io.ReadFull(rand.Reader, randToken)
	io.WriteString(hash, string(randToken))
	key := fmt.Sprintf("%x", hash.Sum(nil))

	_, err = db.Exec("INSERT INTO users (email, password, auth_key) VALUES (?, ?, ?)", email, password, key)
	if err != nil {
		return "", err
	}

	return key, nil
}

// Activate a user from a key
func activateUserWithKey(key string) error {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM users WHERE auth_key=?", key)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return nil
	}
	var id int
	if err = rows.Scan(&id); err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET authorized=1 WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

// Get a value indicating if the user is enabled to see special fields
func isUserEnabled(email string, password string) (bool, error) {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT authorized FROM users WHERE email=? AND password=?", email, password)
	if err != nil {
		return false, err
	}

	if !rows.Next() {
		return false, nil
	}

	var authorized bool
	if err = rows.Scan(&authorized); err != nil {
		return false, err
	}

	return authorized, nil
}

// Load
func mapQueryResults(rows *sql.Rows, includeSpecials bool) ([]dbElement, error) {
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

		res := make(dbElement)
		for i, col := range values {
			if !includeSpecials && config.SpecialFields != nil {
				shouldContinue := false
				for _, s := range config.SpecialFields {
					if s == columns[i] {
						shouldContinue = true
						break
					}
				}
				if shouldContinue {
					continue
				}
			}
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

func getCars(includeSpecial bool) ([]dbElement, error) {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + config.TableName)
	if err != nil {
		return nil, err
	}

	mapRows, err := mapQueryResults(rows, includeSpecial)
	if err != nil {
		return nil, err
	}

	return mapRows, nil
}

func getCar(id int, includeSpecial bool) (dbElement, error) {
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM "+config.TableName+" WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	mapRows, err := mapQueryResults(rows, includeSpecial)
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

func getMostRecentFileWithPrefix(dir string, prefix string) (os.FileInfo, error) {
	// Search for files
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Early return if no file present
	if len(infos) == 0 {
		return nil, nil
	}

	// Select most recent valid archive infos
	sort.Sort(sort.Reverse(byModTime(infos)))
	var info os.FileInfo = nil
	for _, i := range infos {
		if strings.HasPrefix(i.Name(), prefix) && strings.HasSuffix(i.Name(), ".zip") {
			info = i
			break
		}
	}

	return info, nil
}

// Search the most recently added zip file in config.ZipPath,
// extracts any image (keeping the path) with config.ImagesExtension in config.ImagesPath
// then read CSV files and replaces database content with the CSV content.
func updateDatabase() error {
	// Open database
	db, err := sql.Open("mysql", config.DatabaseConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete zip files when done
	defer func() {
		if err != nil {
			return // do not delete anything if import failed
		}
		// Delete all files listed earlier
		infos, _ := ioutil.ReadDir(config.ZipPath)
		for _, info := range infos {
			os.Remove(config.ZipPath + info.Name())
		}
	}()

	for idx, prefix := range config.ZipPrefixes {

		for folderIdx, zipPath := range [2]string{config.ZipPath, config.ZipPathKeep} {

			// Get most recent zip file to import
			info, err := getMostRecentFileWithPrefix(zipPath, prefix)
			if err != nil {
				return err
			}
			if info == nil {
				continue
			}

			// Calculate file_source
			file_source := idx + 10*folderIdx

			// Drop old table
			_, err = db.Exec(fmt.Sprintf("DELETE FROM %v WHERE file_source=%v", config.TableName, file_source))
			if err != nil {
				return err
			}

			// Unzip archive
			archive, err := zip.OpenReader(zipPath + info.Name())
			if err != nil {
				return err
			}
			defer archive.Close()

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
					// HAX add `converter.latinToUtf8[151] = []byte("-")` to latinx/part1.go
					ffutf8 := latinx.NewReader(latinx.ISO_8859_1, ff)

					// Open CSV reader
					c := csv.NewReader(ffutf8)
					c.Comma = rune(config.CsvComma[0])
					c.FieldsPerRecord = -1
					c.TrimLeadingSpace = false
					c.TrailingComma = true
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

					// Prepare query string and put csvElements in database
					columnsCount := len(config.TableMapping)
					for _, element := range elements {
						insertQueryString := "INSERT INTO " + config.TableName + " (file_source,"
						values := make([]interface{}, 0, columnsCount+1)
						values = append(values, file_source)
						for _, column := range config.TableMapping {
							// Check if the element has the column
							if _, exist := element[column.Alias]; exist == false {
								continue
							}
							// Insert in query
							insertQueryString += column.TableColumn + ","
							// Transform the column value and store it
							if column.Transformer != "" && transformers[column.Transformer] != nil {
								element[column.Alias] = transformers[column.Transformer](element[column.Alias])
							}
							values = append(values, element[column.Alias])
						}
						// Finalise query string and create query
						insertQueryString = strings.TrimRight(insertQueryString, ",")
						insertQueryString += ") VALUES ("
						valuesCount := len(values)
						for i := 0; i < valuesCount; i++ {
							insertQueryString += "?"
							if i+1 < valuesCount {
								insertQueryString += ","
							}
						}
						insertQueryString += ")"
						insertQuery, err := db.Prepare(insertQueryString)
						if err != nil {
							return err
						}
						// Execute insertion
						_, err = insertQuery.Exec(values...)
						if err != nil {
							return err
						}
						insertQuery.Close()
					}
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

		}
	}

	return nil
}

func outputError(w http.ResponseWriter, err error) {
	w.Header().Set("Status", "500 Internal Server Error")
	io.WriteString(w, "{\"error\": \""+err.Error()+"\"}")
	logger.Fatalln(err)
}
