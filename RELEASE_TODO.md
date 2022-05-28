TODO before going live:
- Change stripe bank account to new personal bank account
- Stripe API key needs to be switched to the live key.
- Remove the "this site is in test mode" message in course_header.html.
- Reset db. There is test data in it that we should not mix with live data. Also, wiping the db will give us a clean slate.

Features in order of importance:
- Hierarchy course system
	- pre-requisite courses
		- DONE
	- course level system
		- TODO