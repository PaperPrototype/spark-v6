TODO before going live:
- Change stripe bank account to new personal bank account [DONE]
- Stripe API key needs to be switched to the live key.
- Remove the "this site is in test mode" message in course_header.html.
- Create a new db, update DB_URL to new DB, and then delete the old db.
    - WRONG: Reset db. There is test data in it that we should not mix with live data.
 
Features in order of importance:
- Make new CourseOwned model
    - currently using CoursePurchase which is bad because it contains sensitive information related to a users payment
    - why:
        - when course started vs when purchased
        - after user purchases, leave the course as "not started" and we can track and store progress in CourseOwned
            ```
            CourseOwned {
                UserID    uint
                Started   bool
                StartedAt time
                
                // course release has 'postsToCompleteCourse'
                // should we store progress as 0 to 100? or the number of posts the user has made '0 out of 2 posts'?
                // DO BOTH?
                Progress   float
                PostsCount uint

                Completed bool
                
                CourseID  uint
                ReleaseID uint
            }
            ```

- Hierarchy course system
    - TODO (in order of importance)
        - when searching use /api/courses endpoint (not /api/level/:level endpoint)
        - upcoming courses in course.html page
        - rename "course" to "overview" in course.html page
    - DONE
        - course hierarchy
	        - search page has levels/hierarchy
            - pre-requisite courses in course.html page