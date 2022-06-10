Note: the MOST IMPORTANT todo's are in the RELEASE_TODO.md file, so make sure you finish those first before we go live!

- user profile
	- add users showcase posts per course
	- add review posts section
	- allow user to have their own blog posts
- Likes on showcase posts
- Final Projects?
	- NO?: instead just stick with portfolio posts, and encourage students to team up with another student.
		- will have to use the DM's (Direct Messaging) to communicate, which could be annoying.
		- allow users to add tags to their posts? and posts can be foudn by tag organization?
			- cons:
				- when showcasing a final project how will we display posts from that project?
				- posts will be able to have tags regardless.
				- at first just have tags.
				- then we can add the posts group/playlist system?

## More important
- likes for showcase posts
    - so that the `GetReleasePostsOrderByLikes` function in db/get.go will actually be able to order posts by their liking
- support multiple languages other than an `english.md` file

## Chore Todo's (stuff that can be done later)
Use admin interface to manage an "official" courses relational table
	- OfficialCourses
		- (course_id)
	- If your course is official we only take 15% otherwise we take 20% percent?

- if user owns release, and new release comes out, they get a discount since they own the old release
	- release.UpgradeDiscount
- users home page (in `/home`) course links should take user to exact version they are working on
	- if user updates version they are seeing then create "user_current_version"?

- generate 5% off coupon when users land on 404 not found page
- course hierarchy system.
	- if Course.Level = 0 it has to be free.
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
