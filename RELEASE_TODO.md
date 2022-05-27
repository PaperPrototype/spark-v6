TODO before going live:
- Stripe API key needs to be switched to the live key.
- Remove the "this site is in test mode" message in course_header.html
- Reset db. There is test data in it and, it will give us a clean slate.
- Set GIN_MODE to "release" in heroku environment variables
- Detect if user has completed onbaording their stripe account https://stripe.com/docs/connect/express-accounts#handle-users-not-completed-onboarding check the details_submitted parameter

Features in order of importance:
- Final Projects? 
	- NO?: instead just stick with portfolio posts, and encourage students to team up with another student.
		- will have to use the DM's (Direct Messaging) to communicate, which could be annoying.
	- YES?: make final project relational table, where users can write proposal posts and showcase posts. this will go on the "projects" section of the users profile, and will show all participants who worked on that final project.
		- private chat for the project
		- post organization?:
			- allow project owner to create groups/playlists for posts in the project?
				- don't like this as much
				- pros:
					- when showcasing a final project author can simply pick a post group to display as the "showcase" posts.
			- allow users to add tags to their posts? and posts can be foudn by tag organization?
				- I like this better
				- cons:
					- when showcasing a final project how will we display posts from that project?
			- DO BOTH?:
				- posts will be able to have tags regardless.
				- at first just have tags.
				- then we can add groups?