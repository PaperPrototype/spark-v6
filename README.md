# My Vision of the end product
> scroll to the bottom for setup instructions to run the app locally.

If your part of the team it's gonna help to know what exactly what you are building. Well, it's not just courses. It's a learning experience. And I want to solve some big hickups, with ridiculously simple solutions. I'll use Udemy as an example since they are the biggest competitor.

1. No organization to courses or logical pathway to the ultimate learning goal (eg. Make minecraft game, make multiplayer MMO, make a game engine. Everyone who codes has one, it's our job to inspire them to remember). Essentially, I could take an intro to scripting course (if I can even find a good one) but then what? **There's isn't a logical path to the next course to take based on what I'm interested in.** There is an incredibly simple solution to this. Any course can set other courses as pre-requisites. And with that your courses start to form a hierarchy or "graph", that you can use to display "next courses to take", and "pre-requisite courses".
   - teachers know what the student has already learned based on selected pre-requisite courses
   - students don't waste time learning stuff they've already learned

2. Assistence and help when you get stuck. Udemy has a "Q&A" section for each course. They've attacked the problem in the most direct way possible. But what are people using even more than the Q&A system? Discord. With Discord you can make more than just a Q&A section, you can organize your community, make announcements, threads, showcase channels. Discord is a place for your student to come tegther as a community. So what we are we going to do? Each course is going to be like it's own dedicated discord server.

There is a few trivial but important features that need to exist
- Course progress system. When you finish a section its marked as complete
- Profile system that looks like the coolest portfolio ever.
    - The Discord system also comes with simple (think of medium.com) posts for things like showcases. I'm still stuck on how to keep continuity between the posts and messages in the discord like chat system.


## Lets take a visionary tour of what I'm invisioning
Imagine a courses site where each course is like a dedicated discord server. You have live chat, voice, and video. As well as posts that showcase your progress and skills and what you are capable of. 

Just like discord, you get live noifications to help you find where people have mentioned you. The best part is that you can view the chat and posts without interrupting your work on the course, because the chat *is* part of the course! There is even an overlay DM voice and video system so you can chat directly with your best friend while working on a course together!

The priorty customer is the END USER (the students). Not the teachers, or the companies using our site for hiring. We want to be like apple. Our students are king, even at the cost of possibly losing a profitable deal with another company. No ad's (well at the beginning we may use google ads to keep from going under). We make and sell really good courses. Thats our business model.

Eventually courses will have to start costing more (as an example: apple is expensive, but worth it). We can get students hooked by making all lower level courses free, and then allow teachers to charge more money the higher the level of the course is (the level of a course depends on the number of prerequisites the student has to go through to get to your course).

### Hierarchy system
The best part is the course hierarchy system. Every course is part of a massive "tree" of courses. Unlike other sites where you have to look through thousands of random courses, at sparker3d.com, courses can set pre-requisite courses! That way you know where to start, and what courses to take next!

Courses are also given a "Level" based on the depth of pre-requisite courses it has (Courses at Level 0 are always free!). You "level up" by taking courses that are at higher levels. It gives the community a ranking system, the highest level students are held in awe by their peers because of their "rank".

Without realizing it, you get sucked in (like a video game!) and end up building a really awesome portfolio from all the showcase posts you made! Now you start posting your profile on job hiring websites. No actually, sparker3d.com is a hiring website, based on skill, not certification.

### Job hiring (once we have 1 million users)
Since our site is kinda becoming the future of online tech portfolios, we add features targeted at helping companies and entreprenours hire our students. Our top priority is to benefit and protect the students at all costs. Again, the priorty customer is the END USER (the students). Not the teachers, or the companies using our site for hiring. We want to be like apple. Our students are king.

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

## stripewebhook
from stripe create a new webhook that points to `/stripe/webhooks` and set it up to listen for the events locally. This is for payments

# Running
Use `go run .` once you have the config files created and set.
You may need to install dependancies using `go mod tidy`.