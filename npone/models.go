package npone

import (
	"fmt"
)

// Course represents one course in the system
type Course struct {
	ID          uint   `gorm:"primary_key;auto_increment;"`
	Name        string `gorm:"type:varchar(256)"`
	Description string `gorm:"type:varchar(1024)"`

	Lessons []Lesson
}

// PrintNicely prints a course in a nice way
func (c Course) PrintNicely() {
	fmt.Printf("Course: (%d)\n", c.ID)
	fmt.Printf("Lessons:\n")
	for _, lesson := range c.Lessons {
		lesson.PrintNicely(false)
	}
}

// Lesson represents one lesson in the system
type Lesson struct {
	ID          uint   `gorm:"primary_key;auto_increment;"`
	Name        string `gorm:"type:varchar(256)"`
	Description string `gorm:"type:varchar(1024)"`
	CourseID    uint   `sql:"type:int(10) unsigned REFERENCES courses(id);index"`

	Course Course
}

// PrintNicely prints a course in a nice way
func (l Lesson) PrintNicely(courseDetails bool) {
	fmt.Printf(" Lesson: (%d)\n", l.ID)
	if courseDetails {
		l.Course.PrintNicely()
	}
}
