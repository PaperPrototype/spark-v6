### The most important things you'll hear me say are here
This document is very important. Please to not edit it without talking to one of the founders (Abdiel, David).

Do not edit QUESTIONS.md unless you are adding a question
Feel free to edit RELEASE_TODO.md

> Scroll to the bottom for setup instructions to run the app locally.
> Check the RELEASE_TODO.md for critical features that need to be added.

# Lets get started
If you are part of the team it's gonna help to know what exactly what it is we are building so lets get started. 

It's not just courses, it's a learning experience. We want to solve online learnings biggest problems, with ridiculously simple solutions. I'll use Udemy as an example for comparison since they are our biggest competitor.

1.
    Problem: No organization to courses or logical pathway to the ultimate learning goal (eg. Make minecraft game, make multiplayer MMO, make a game engine. Everyone who codes has one, it's our job to inspire them to remember). Essentially, I could take an intro to scripting course (if I can even find a good one) but then what? **There's isn't a logical path to the next course to take based on what I'm interested in.** There is an incredibly simple solution to this. 
    
    Solution: Allow a course to set other courses as pre-requisites. Now as more courses are added a hierarchy or "graph" starts to form. We can now display "next courses to take", and "pre-requisite courses". When someone makes a new course allow them to select courses to set as pre-requisites courses.

    - teachers know what the student has already learned based on selected pre-requisite courses
    - students don't waste time learning stuff they've already learned
    
    **Marketing:**

    I realize that trying to sell this concept directly would be bad marketing. People buy things that they believe they need. So we should make "paths" (groups of courses with a final goal that users are asking for) because that is what customers will buy! Like say a minecraft series of courses, we can present the entire series as a "path" of courses. They'll look at the last courses in the series(planetary terrain, multiplayer, minecraft clone, minecraft trees, procedural voxel destruction, voxel terrain LOD, marching cubes, Unity native mash API) and buy them (as well as the pre-requisites) because they want to get to the end exciting stuff! I'd say 90% of sales will happen if we can paint this perspective of "series" right (a dedicated section in browse for "series", and a "series" page that displays a series of course and their "tree structure" of connections).

2. 
    Course Communities + Assistence and help when studetns get stuck. 
    
    Udemy has a "Q&A" section for each course. They've attacked the problem in the most direct way possible. But what are people using even more than the Q&A system? Discord servers. With Discord you can make more than just a Q&A section, you can organize your community, make announcements, threads, showcase channels. Discord is a place for your student to come tegther as a community. So what we are we going to do? Each course is going to be like it's own dedicated discord server.

    **Marketing:**

    Marketing is selling a good to saisfy a need. People need a place to ask questions and get help. They also will create or find a community to satisfy conversing with poeple with similar interests.
    Each course should display what students are currently live online (like discord). This gives a sense of not being alone. And its encouraging to know other people are working on the same thing you are.
3. 
    Posts system like medium. The purpose is for the Profile system to create the coolest portfolio ever. Utilizing showcase posts the user has made. I'm still stuck on how to keep continuity between the posts system and the messages system in the discord like chat system.

    **Marketing:**

    I have no idea. But first we need to figure out why students would want to make a post vs use our chat system. Somehow I need to change it from "sparker needs ppl to make posts" to "students want/need to post because..." 

    **Marketing Ideas:**
    I say it's best used as a method of showcasing progress, and users who want to indirectly add tutorials or information they learned to a course. Its also a more permanent way of showcasing work/sharing thoughts than a fast-paced chat system where conversations get forgotten.

And thats our MVP! So far none of the above have been implemented, and the only things working are:
- courses hosted through github
    - We have 2 courses that are active, and 1 of them is paid.
- payments for course releases
    - each course can have multiple releases (like editions of a book). Each release is pruchased separately.
- payment errors cause an email to be sent to the author, and gift the course for free to the customer
    - customer satifaction because something went wrong, we take the hit and make sure to satisfy the customer
        - someone should look at that code and find ways to ensure those errors never can happen
- onboarding, login / signup
    - email verification (not automatically sent when user signs up) if user wants to publish a course they have to verify their email through `/settings`. check `/settings` URL in `router2` package to see what I mean

Now you know what needs done. Go do it! Ask for help, and get your hands dirty. Remember someone else may already be working on the thing you want to work on, so ask around (and if someone asks you what you're working on, make sure to answer them cause we don't need poeple wasting time). 

less critical MVP features
- Course progress system. When you finish a section its marked as complete
- Ranking system. Courses with more pre-requisites get a higher "level" (think "levels" like in a video game). This will could bring out the competitive nature of some pestudents, and make a system for evaluating how deep a students knowledge has gone. Basically the higher level course a student takes the higher their rank gets.

Make sure to read the [QUESTIONS.md](/QUESTIONS.md) file in the repo as it outlines questions that we should have answers for!

There is also a list of the most important questions that I have **not** figured out:

1. 
    Who is gonna make the course other than me? Should we hire teachers? I personally like making courses, but this is not going to work long term.
2. 
    How are we going to break even while we continue to develop the site. How can we increase our runway as much as possible!
    - google ads
        - this needs to be done as soon as possible.
    - a few high demand niche courses
        - voxel planets <= a course that I am working on right now
3. 
    Password resetting through email. Yeah, still need to get that one working.
4. 
    Customer support. Not even sure where to start with that one. Currently the Discord server
5. 
    Server costs with google ads. Thinking of popping in some **unabtrusive** ads along the sides of the course and at the bottom to help cover server costs.
6. 
    Accountability structure for me and other founders to keep us from going off the rail.
7. 
    What is Unit economics?
8. 
    I'm putting the follwing here because its important, and I want help figuring this out properly. It's very important that I figure it out as soon as possible since it *could* become a very BIG problem if it's not figured out (I'm also going to be very blunt and not hide what I'm thinking): 5% of the company is David Spooners (currently), and 3% is my parents (currently). I need to clarify what exactly this means for the company. I also need to know what this will mean legally for the company and what 3% evem means. Stock? Shareholders stock? Voting stock? Equity? Profit share? Once I figure it out I need to make a it legal and write  operating agreement with my parents and uncle. I want to my uncles and parents percentage to be used for VC's (not all of it, but I'm not sure I want them holding 9% of the company). Does that mean buying some of it back from them? Should I add something to the operating agreement? Like maybe pay them back the money they've invested plus let them have 1% of the company? I want them to get their fair share, but I also don't want to regret giving 9% to them. Reason I gave it in the first place? My parents are paying the server fee's and my uncle is paying for the domain name and SSL certificate. Without them I'd still be stuck without a domain and dedicated servers. Also I have 2 emails `abdiellopez@sparker3d.com` and `info@sparker3d.com`. We need them to keep paying for server fee's and domain names since they are our runway right now. Literally they are funding this thing. Server fee's are 8$ per month (parents). domain name  + SSL + 3 emails all through namecheap (paid by uncle) are about 84.15$ per year.

    Here is a link to a reddit that has some good answers related to this https://www.reddit.com/r/legaladvice/comments/y3rp3n/my_mom_wants_13_of_my_business_and_13_or_more_of/ (not that I am having trouble with my parents. They will probably be happy to just give up their percentage, but I want them to be treated fairly)

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