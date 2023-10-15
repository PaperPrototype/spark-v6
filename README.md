For stripe payment routes see `stripe_webhooks` `stripe_buy` `stripe_connect` `stripe_checkout` in the `router2` folder.

# Developer info
`router` and `templates` and `resources` were the old site. I've kept them for legacy purposes since I use them for example code of things I've done before. `templates2` and `router2` and `resources2` are the replacements.

# Local Setup
config vars are stored in single key files when running locally. These files must be in the repos root folder. They are ignored in `.gitignore` for security purposes.

(`.gitignore` file)
```
dbconfig
stripeconfig
sendgridconfig
githubclientsecret
githubclientid
stripewebhook
```

## dbconfig
place the database connection info into the dbconfig file

(example)
```
port = 5432
dbname = spark-v6
sslmode = disable
```

if you are using postgres you can create a db using the `createdb` command. To find the port number (on MacOS) open the postgres app and click on `Server Settings...`

## stripeconfig
go to stripes docs and you can use the publicly available test key. Or if you are logged in, use the test key provided and paste it into a `stripeconfig` file in the root directory of the app (the files name must be exactly stripeconfig).

## sendgridconfig
go to app.sendgrid.com docs and you can create a free account. Or if you are logged in, use the test key provided and paste it into a `sendgridconfig` file in the root directory of the app.

## githubclientsecret githubclientid
create two files called `githubclientsecret` and `githubclientid` in the root directory. Now create an oauth app on github and register the homepage url as `http://localhost:8080` and the redirect url as `http://localhost:8080/settings/github/connect/return`

now paste the clientid and clientscret into the files. Now github oauth should work locally.

## stripewebhook
from stripe create a new webhook that points to `/stripe/webhooks` and set it up to listen for the events locally. This is for payments

# Running
Use `go run .` once you have the config files created and set.
You may need to install dependancies using `go mod tidy`.
