TODO before going live:
- get webhook signing key
    - set STRIPE_WEBHOOK env variable for payments webhook
- get async payments working with the checkout flow [DONE]
    - https://stripe.com/docs/payments/checkout/fulfill-orders#delayed-notification
 
Features in order of importance:
1. Make new CourseOwned model [DONE]
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

2. Hierarchy course system
    - display prerequisite courses for a course somehow
    - display upcoming courses somehow
    - a way to search and set prerequisite courses for a course

3. Offline viewing of website page (like youtube)
http://diveinto.html5doctor.com/offline.html
    - make course pages viewable offline
        - change sections to be stored in "localStorage" instead of "sessionStorage"
        - add cache manifest so browser will download resources that can be loaded when offline