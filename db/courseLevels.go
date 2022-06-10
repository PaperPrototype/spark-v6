package db

import "log"

/*
	DB QUERY'S
*/

// get all the courses that are at a certain level
// TODO pagination and max limit
func GetCoursesAtLevelPreload(level interface{}) ([]Course, error) {
	courses := []Course{}

	err := gormDB.Model(&Course{}).Preload("User").Preload("Release", orderByNewestCourseRelease).Where("public = ?", true).Where("level = ?", level).Find(&courses).Error

	return courses, err
}

/*
	UPDATING THE LEVEL OF A COURSE (all below)
*/

// this function does the following
// update the "level" of its course by checking the "depth" of prerequisites it has
// update the "level" of any course that has set this course as its prerequisite
func UpdateCourseLevel(courseID uint64) error {
	log.Println("db/courseLevels updating course level based on its pre-requisites")

	// get all pre-requisites
	preqs := []Prerequisite{}
	gormDB.Model(&Prerequisite{}).Where("course_id = ?", courseID).Find(&preqs)

	// check the depth of each pre-requisite, and find the longest branch
	depths := []uint32{}
	for _, preq := range preqs {
		// recursively check the depth of pre-requisites, and find the longest branch
		// add it to a list
		depths = append(depths, getDepthRecursive(preq.PrerequisiteCourseID, 1))
	}

	// find the greatest depth
	greatestDepth := greatest(depths)

	// set the course level = greatestDepth
	err1 := gormDB.Model(&Course{}).Where("id = ?", courseID).Update("level", greatestDepth).Error
	if err1 != nil {
		log.Println("db/courseLevels ERROR updating course level?", err1)
		return err1
	}

	return nil
}

func getDepthRecursive(preqCourseID uint64, depth uint32) uint32 {
	// get all pre-requisites for this pre-requisite course
	var preqCourses []Prerequisite
	err := gormDB.Model(&Prerequisite{}).Where("course_id = ?", preqCourseID).First(&preqCourses).Error
	if err != nil {
		log.Println("db/courseLevels ERROR getting course preq (recursive function found its end, so possibly ignore this), ", err)
		return depth
	}

	// store all the depths
	depths := []uint32{}
	for _, preq := range preqCourses {
		depths = append(depths, getDepthRecursive(preq.PrerequisiteCourseID, depth+1))
	}

	// give back the greatest depth
	return greatest(depths)
}

func greatest(depths []uint32) uint32 {
	var greatest uint32 = 0
	for _, depth := range depths {
		// if depth greater than 0
		if greatest < depth {
			// increase greatest
			greatest = depth
		}
	}

	return greatest
}

// When updating a courses pre-requisites we should update the level of other courses
// that have this course as their pre-requisite! since the pre-requisite depth/level may now be different!
func UpdateDependantCourseLevels(courseID uint64) error {
	// find courses that depend on this course
	var dependants []Prerequisite
	err := gormDB.Model(&Prerequisite{}).Where("prerequisite_course_id = ?", courseID).Find(&dependants).Error
	if err != nil {
		log.Println("db/courseLevels ERROR getting preq from db (you can probably ignore this error since a recursive function has just found it end),", err)
		// don't give back an error since nothing wen't wrong and this is expected to happen
		// when we reach the end of our recursive function
		return nil
	}

	// update the level of courses that depend on us
	for _, dependant := range dependants {
		// ignore error since it is
		UpdateCourseLevel(dependant.CourseID)
		// update course levels of any courses that depend on this course
		err1 := UpdateDependantCourseLevels(dependant.CourseID)
		if err1 != nil {
			return err1
		}
	}

	return nil
}
