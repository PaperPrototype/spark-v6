TODO before going live:
- get webhook signing key
    - set STRIPE_WEBHOOK env variable for payments webhook
- get async payments working with the checkout flow [DONE]
    - https://stripe.com/docs/payments/checkout/fulfill-orders#delayed-notification
 
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