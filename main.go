package main

import (
	"flag"
	"fmt"

	"github.com/arpitbbhayani/npone-impact/npone"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	host := flag.String("host", "localhost", "host of mysql")
	port := flag.Int("port", 3306, "port of mysql")
	user := flag.String("user", "root", "user of mysql")
	pass := flag.String("pass", "", "password of mysql")
	dbname := flag.String("db", "npone", "db of mysql")

	coursesCount := flag.Int("courses", 100, "total number of courses")
	maxLessonsPerCourse := flag.Int("maxlessons", 15, "max number of lessons per course")
	courseFetchCount := flag.Int("coursefetch", 10, "total number of courses to be fetched")

	flag.Parse()

	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *pass, *host, *port, *dbname)

	db, err := gorm.Open("mysql", connectionStr)
	if err != nil {
		panic(err)
	}
	db.LogMode(false)

	courses := npone.Populate(db, *coursesCount, *maxLessonsPerCourse)
	fmt.Printf("%d courses will be fetched per call\n", *courseFetchCount)

	withoutPreloadOutput := npone.Benchmark(db, courses, *courseFetchCount, npone.GetCoursesByIDsWithoutPreload)
	withPreloadOutput := npone.Benchmark(db, courses, *courseFetchCount, npone.GetCoursesByIDsWithPreload)

	fmt.Println(withoutPreloadOutput)
	fmt.Println(withPreloadOutput)

	defer db.Close()
}
