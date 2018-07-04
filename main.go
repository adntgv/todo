package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Task model for gorm
type Task struct {
	gorm.Model
	Text string
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

func listTasks(db *gorm.DB) {
	var tasks []Task
	db.Find(&tasks, &[]Task{})
	if len(tasks) < 1 {
		fmt.Println("There are no tasks for now")
	} else {
		for _, task := range tasks {
			fmt.Println(task)
		}
	}
}

func main() {
	list := flag.Bool("list", false, "list all tasks")
	new := flag.Bool("new", false, "add a task")
	flag.Parse()

	dbParams, err := readParams("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open("postgres", dbParams)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if db.HasTable(&Task{}) == false {
		db.CreateTable(&Task{})
	}
	switch {
	case *list:
		{
			listTasks(db)
		}
	case *new:
		{
			task := Task{Text: strings.Join(flag.Args(), " ")}
			fmt.Println(task)
			db.Create(&task)
		}
	}
}
