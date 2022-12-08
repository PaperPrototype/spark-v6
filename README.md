### The most important things you'll hear me say are here
This document is very important. Please to not edit it without talking to one of the founders (Abdiel, or David).

Do not edit QUESTIONS.md (unless you are adding a question)
Feel free to edit RELEASE_TODO.md

> Scroll to the bottom for setup instructions to run the app locally.
> Check the RELEASE_TODO.md for critical features that need to be added.

# Lets get started
If you are part of the team it's gonna help to know what exactly what it is we are building so lets get started. 

The only things working are:
- courses hosted through github
    - We have 2 courses that are active, and 1 of them is paid.
- payments for course releases
    - each course can have multiple releases (like editions of a book). Each release is pruchased separately.
- payment errors cause an email to be sent to the author, and gift the course for free to the customer
    - customer satifaction because something went wrong, we take the hit and make sure to satisfy the customer
        - someone should look at that code and find ways to ensure those errors never can happen
- onboarding, login / signup

Now you know what needs done. Go do it! Ask for help, and get your hands dirty. Remember someone else may already be working on the thing you want to work on, so ask around (and if someone asks you what you're working on, make sure to answer them cause we don't need poeple wasting time). 

less critical MVP features
- Course progress system. When you finish a section its marked as complete
- When there is more than 1 release, the course page breaks/stops working

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
    Video or text based courses? Right now we are forced to do mostly text. I've been inspired by Descript's text-first video editing software. 

    Text-first appeals to teachers.
    Video appeals to students.

    text-first article creation with video is an idea I love.

    When editing a courses section think of medium.com post editing + discord channel management, but you can drag and drop video files right into the post.
    - for now embed youtube videos when needed using markdown syntax.
    - will have to wait until we have more capital and resources. Continue to invest in github based courses until then.

    An embeddable processing.js system like khan academy has for its javascript course is awesome. This for will be defered until we have more resources to build something like this. 
    1. Defer students to the khan academy course?
    2. Do a video course and have students follow along in editor.p5js.com?

4. 
    Password resetting through email. Yeah, still need to get that one working.
5. 
    Customer support. Not even sure where to start with that one. Currently the Discord server
6. 
    Server costs with google ads. Thinking of popping in some **unabtrusive** ads along the sides of the course and at the bottom to help cover server costs.
7. 
    Accountability structure for me and other founders to keep us from going off the rail.
8. 
    What is Unit economics?
10. 
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