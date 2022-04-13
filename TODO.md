Note: the MOST IMPORTANT todo's are in the RELEASE_TODO.md file, so make sure you finish those first before we go live!

## More important
- support multiple languages other than an `english.md` file
- chat
- projects (final projects and proposals system)
- course hierarchy 
	- search page can have a hierarchy

## Chore Todo's (stuff that can be done later)
- users home page (in `/home`) course links should take user to exact version they are working on
	- if user updates version they are seeing then create "user_current_version"

- generate 5% off coupon when users land on 404 not found page
- course hierarchy system.
	- if Course.Level <= 1 it has to be free.
- in case if user forgets to claim purchase use the stripe webhook to verify and claim purchase for the user
	- much more robust
	- user may close page and/or close the browser and lose their payment, but using webhooks solves this
- project
	- (user_id, state = "proposal" || "completed" || "in-progress")
		- proposalPost (post_id, project_id)
		- projectUsers (project_id, user_id)
- finalProject (links a project to a course release)
	- (release_id, project_id)

- user profile
	- courses
		- Authored
		- Taking
	- projects
		- proposals
		- projects
	- posts
		- series

- dark mode button toggle
	- dark
	- light
	- auto

Use admin interface to manage an "official" courses relational table
	- OfficialCourses
		- (course_id)
	- If your course is official we only take 15% otherwise we take 25% percent