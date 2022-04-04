config vars are stored in single key files. These files must be in the repos root folder. They are ignored in `.gitignore` for escurity purposes.

```
dbconfig
stripeconfig
sendgridconfig
```

## dbconfig
place the database connection info into the dbconfig file

(example)
```
port = 5432
dbname = spark-v6
sslmode = disable
```

if you are using postgres you can create a db using the `createdb` command. TO find the port number (on MacOS) open the postgres app and click on `Server Settings...`

## stripeconfig
go to stripes docs and you can use the publicly available test key. Or if you are logged in, use the test key provided and paste it into a `stripeconfig` file


## sendgridconfig
go to app.sendgrid.com docs and you can create a free account. Or if you are logged in, use the test key provided and paste it into a `sendgridconfig` file

# Running
Use `go run .` once you have the config files created and set.
You may need to install dependancies using `go mod tidy`.

TODO
- generate 5% off coupon when users land on 404 not found page
- course hierarchy system.
	- if Course.Level <= 1 it has to be free.
- in case if user forgets to claim purchase use stripe webhook to verify and claim purchase for the user
	- much more robust
	- user may close page and/or close the browser and lose their payment, but using webhooks solves this
- github oauth
	- github linkify
		- when user edits course contents, allow special button for "commit changes to github"
		- github based course serving??
- blog posts, course post playlists
- course settings
- profile
- chat
- final project proposal posts
	- posts chat
	- proposer can accept participants.
	- once more than 1 participant joins, they can begin project. 
	- "Final Project" flow? db data layout?
		- ideas
			- allow for projects to be made outside of course, specifics
			- projects go on "Projects" page of user
			- project can just be a blog post? Added to a specific db relation "PostToProject"?
			- private chat for PostToProject?
			- Post (id, user_id)
				- PostAuthors (post_id, user_id)
				- PostToRelease (post_id, release_id)
				- PostToProject (post_id, proposal)
					- proposal: if the post is currently only a proposal
				- Chat (post_id, private, id)
					- Comments
						- Comment (chat_id, user_id, )
- user profile
	- courses
		- Authored
		- Taking
	- projects
		- proposals
		- projects
	- posts
		- playlists

- dark mode button toggle
	- dark
	- light
	- auto