# My Vision of the end product
Below I will detail the vision for the final site and its features and user experience.

In one sentence: To make the most addictive learning experience ever, and the more you do it the more it benefits you (because course "progress" is measured by the amount you add to your portfolio).

Imagine a courses site where each course is like a dedicated discord server. You have live chat, voice, and video. But you can also make blog posts to showcase your progress and what you are making. You get live notifications when someone mentions you.

The best part is that you can view the chat and posts without interrupting your course. There is even an overlay DM voice and video system so you can chat directly with your best friend while working on a course together!

You can see all the user's who are currently working the same course as you, and read their bio's and even see their posts. All this while always being able to simply click back (without the page re-loading) and get right back to where you were in your course.

The priorty customer is the END USER (the students). Not the teachers, or the companies using our site for hiring. We want to be like apple. Our students are king, even at the cost of losing a profitable deal with another company. No ad's. We make and sell really good courses. Thats our business model.

Eventually courses will have to start costing more (apple is expensive, but worth it). We can get students hooked by making all lower level courses free, and then allow teachers to charge more money the higher the level of the course is.

## Hierarchy system
The second best part is the course hierarchy system. Every course is part of a massive "tree" of courses. Unlike other sites where you have to look through thousands of random courses, in sparker.com, courses can set pre-requisite courses! That way you know where to start, and what courses to take next!

Courses are also given a "Level" based on the depth of pre-requisite courses it has (Courses at Level 0 are always free!). You "level up" by taking courses that are at higher levels. It gives the community a ranking system, the highest level users are held in awe of their accomplishments.

Without realizing it, you get sucked in and end up building a really awesome portfolio from all the showcase posts you've been making! Now you start posting your profile on job hiring websites.

## Job hiring (once we have 1 million users)
Since our site is kinda becoming the future of online tech portfolios, we add features targeted at helping companies and entreprenours hire our students. Our top priority hould be to benefit and protect the students at all costs. The priorty customer is the END USER (the students). Not the teachers, or the companies using our site for hiring. We want to be like apple. Our students are king, even at the cost of losing a profitable deal with another company. No ad's. We make and sell really good courses. Thats our business model.

Eventually courses will have to start costing. But we can get students hooked by making lower level courses free, and then amping the price as they start climbing the hierarchy and take higher level courses.

# Local Setup
config vars are stored in single key files when running locally. These files must be in the repos root folder. They are ignored in `.gitignore` for security purposes.

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

if you are using postgres you can create a db using the `createdb` command. To find the port number (on MacOS) open the postgres app and click on `Server Settings...`

## stripeconfig
go to stripes docs and you can use the publicly available test key. Or if you are logged in, use the test key provided and paste it into a `stripeconfig` file in the root directory of the app (the files name must be exactly stripeconfig).

## sendgridconfig
go to app.sendgrid.com docs and you can create a free account. Or if you are logged in, use the test key provided and paste it into a `sendgridconfig` file in the root directory of the app.

## githubclientsecret githubclientid
create two files called `githubclientsecret` and `githubclientid` in the root directory. Now create an oauth app on github and register the homepage url as `http://localhost:8080` and the redirect url as `http://localhost:8080/settings/github/connect/return`

now paste the clientid and clientscret into the files. Now github oauth should work locally.

# Running
Use `go run .` once you have the config files created and set.
You may need to install dependancies using `go mod tidy`.