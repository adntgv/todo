package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Task model for gorm
type Task struct {
	gorm.Model
	Text     string
	Done     bool
	Deadline time.Time
}

func readParams(configFile string) (string, error) {
	var user, dbname, password string
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		switch {
		case strings.Contains(line, "user"):
			idx := strings.Index(line, "=") + 1
			user = line[idx:]
		case strings.Contains(line, "dbname"):
			idx := strings.Index(line, "=") + 1
			dbname = line[idx:]
		case strings.Contains(line, "password"):
			idx := strings.Index(line, "=") + 1
			password = line[idx:]
		}
	}
	return fmt.Sprintf("user=%s dbname=%s password=%s", user, dbname, password), nil
}

func main() {
	dbParams, err := readParams("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open("postgres", dbParams)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
