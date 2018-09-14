package npone

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	"github.com/jamiealquiza/tachymeter"
	"github.com/jinzhu/gorm"
)

func shuffle(vals []Course) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}

// ResetDatabase empties the database and creates required Tables
func ResetDatabase(db *gorm.DB) {
	db.DropTableIfExists(&Course{}, &Lesson{})
	db.CreateTable(&Course{}, &Lesson{})
}

// PopulateCourses populates the courses table with `count` number of
// courses
func PopulateCourses(db *gorm.DB, count int) []Course {
	courses := make([]Course, count, count)
	for i := 0; i < count; i++ {
		course := Course{
			Name:        fake.Sentence(),
			Description: fake.Paragraph(),
		}
		db.Save(&course)
		courses[i] = course
	}
	return courses
}

// PopulateLessons populates the lessons table with `count` number of
// lessons per course
func PopulateLessons(db *gorm.DB, count int, courses []Course) int {
	maxLessons := count * len(courses)
	lessonIndex := 0
	for i := 0; i < maxLessons; i++ {
		course := courses[rand.Int()%len(courses)]
		lesson := Lesson{
			Name:        fake.Sentence(),
			Description: fake.Paragraph(),
			Course:      course,
		}
		db.Save(&lesson)
		lessonIndex++
	}
	return lessonIndex
}

// GetCoursesByIDsWithoutPreload returns a courses given its IDs without using preloading
func GetCoursesByIDsWithoutPreload(db *gorm.DB, courseIDs []uint) []Course {
	var courses []Course
	db.Find(&courses, "id IN (?)", courseIDs)
	for i := range courses {
		db.Model(courses[i]).Related(&courses[i].Lessons)
	}
	return courses
}

// GetCoursesByIDsWithPreload returns a courses given its IDs using preloading
func GetCoursesByIDsWithPreload(db *gorm.DB, courseIDs []uint) []Course {
	var courses []Course
	db.Preload("Lessons").Find(&courses, "id IN (?)", courseIDs)
	return courses
}

// Populate populates the db for given configuration
func Populate(db *gorm.DB, coursesCount, maxLessonsPerCourse int) []Course {
	ResetDatabase(db)
	courses := PopulateCourses(db, coursesCount)
	lessonsCount := PopulateLessons(db, maxLessonsPerCourse, courses)
	fmt.Printf("%d courses created\n", len(courses))
	fmt.Printf("%d lessons created\n", lessonsCount)
	return courses
}

// Benchmark runs benchmark for given configuration and prints the benchmark output
func Benchmark(db *gorm.DB, courses []Course, courseFetchCount int, fn func(*gorm.DB, []uint) []Course) string {
	t := tachymeter.New(&tachymeter.Config{Size: 100})
	for i := 0; i < 100; i++ {
		for k := 0; k < 5; k++ {
			shuffle(courses)
		}
		courseIDs := make([]uint, courseFetchCount, courseFetchCount)
		for j := 0; j < courseFetchCount; j++ {
			courseIDs[j] = courses[:courseFetchCount][j].ID
		}

		start := time.Now()
		fn(db, courseIDs)
		t.AddTime(time.Since(start))
	}

	return t.Calc().String()
}
