config vars are stored in single key files. These files must be in the repos root folder. They are ignored in `.gitignore` for escurity purposes.

```
dbconfig
stripeconfig
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

## running
Use `go run .` once you have the config files created and set.
You may need to install dependancies using `go mod tidy`.

TODO
- generate coupon for 5% off on 404 not found page