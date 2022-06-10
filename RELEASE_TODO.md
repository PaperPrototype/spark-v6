TODO before going live:
- Change stripe bank account to new personal bank account
- Stripe API key needs to be switched to the live key.
- Remove the "this site is in test mode" message in course_header.html.
- Reset db. There is test data in it that we should not mix with live data. Also, wiping the db will give us a clean slate.

Features in order of importance:
- Hierarchy course system
    - TODO (in order of importance)
        - when searching use /api/courses endpoint (not /api/level/:level endpoint)
        - upcoming courses in course.html page
        - rename "course" to "overview" in course.html page
    - DONE
        - pre-requisite courses
        - course level system in search page
   