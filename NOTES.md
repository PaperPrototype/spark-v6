TODO
- generate 5% off coupon when users land on 404 not found page
- course hierarchy system.
	- if Course.Level <= 1 it has to be free.
- in case if user forgets to claim purchase use stripe webhook to verify and claim purchase for the user
	- much more robust
	- user may close page and/or close the browser and lose their payment, but using webhooks solves this
- chat
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