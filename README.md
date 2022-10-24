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
    Poeple learn things with an end goal in mind.

    Sites like Udemy have no organization to courses or logical pathway to the ultimate learning goal (Everyone who codes has one, and we can teach for those goals to maximize impact), they are just a massive collection of random courses. I see that their disadvantage. 
    
    Udemy Example: I could take an intro to scripting course (if I can even find a good one) but then what? **There's isn't a logical path to the next course to take based on what I'm interested in.** There is an incredibly simple solution to this. Allow a course to set other courses as pre-requisites. 
    
    With this system as more courses are added to ur site a natural hierarchy or "graph" of courses starts to form. We can now display "next courses to take", and "pre-requisite courses" when someone is trying to find their way to their final goal.

    - teachers know what the student has already learned based on selected pre-requisite courses
    - students don't waste time learning stuff they've already learned
    
    **Marketing:**

    I realized that trying to sell the "hierarchy" as a concept sounds cool, but it doesn't show how it solves a users **need**. Marketing is getting poeple to believe they need something so that they buy it. 
    
    Instead I propose we make "paths" (groups of courses with a final goal that users are asking for) tglained with a specific end goal... because that is what customers will buy! Like say a minecraft series of courses dubbed "Make Minecraft". 
    
    The "tagline" obviously will be based on the last course in the series (the final end goal: planetary terrain, multiplayer, minecraft clone, minecraft trees, procedural voxel destruction). The result is that poeple will have to buy the pre-requisite courses first because they want to get to the end-goal! 
    
    I'd say 90% of sales will happen if we can paint this perspective of "series" right (eg: a dedicated section in browse for "series". A "series" page that displays a series of courses. Maybe alos show a "tree structure" of connections?).

2. 
    Getting people to use the site!!! We can have an amazing platform, but if no one uses it then we will fail.

    Courses. We need to make courses that poeple will want to take. 
    
    We could aim for high interest catalyst and exciting but simple to learn niche' subjects that the internet has (so far) failed to fill.

    catalyst example (clickbaitable, something that is exciting and not common)
        make a game engine => simulate billions of players at once in multiplayer (but also make a game engine)
        voxel terrain => destructible planet scale voxel terrain

        - blog posts and tutorials to attract customers
            - "How to build planet scale terrain" announcing planetary voxel terrain course
        - game development subjects
            - planet scale terrain [almost done. low investment high return if we can make it viral. Try to hook PippenFTP build the earth in minecraft youtuber]
            - [high investment low return?] make a game engine (simplified, very simple, non time consuming, and not frustrating)
        - web development subjects
            - make an e-commerece store/website with golang, stripe, heroku, and postgres
            - make a discord clone (with cloudflare workers, D1, R2, DurableObjects, Ember.js)

    Recruit teachers> We need to figure out compnesation (salary, company stocks compensation, profit share with courses?). Should we have "Official courses"? Alternative. Invest in recruiting other users to build courses on our site! Make the process of building courses user friendly and offer "making your first course" tutorial that guides users to make a course (simple, connect to github, make single README.md file section)
    
    Obviously we need courses to kickstart the process. Teachers will teach where there are students.

    Most of the below "series" reuse existing courses!

    On "series" page, have a "next course to buy" that the user can click "buy" on

    - Series: Make Minecraft Planets (quickest and soonest course we can make)
        - Minecraft Terrain Basics (C# Unity) [DONE]
        - Minecraft Worlds - Terrain Generation with Jobs (C# Unity) [ALMOST]
        - Planetary LOD Voxel Terrain [PROTOTYPE]

    - Series: Make Games with Unity
        - Scripting in JavaScript - Make a Flappy Bird game [TODO]
            - temporary: "Learn to code in javascript Khan Academy"
        - Basics of Unity - Make Flappy Bird + (some other popular game) and publish on Itch.io store [TODO]
            - temporary: "Learn to code by making games GameDev.tv"

    - Series: Make Planet Scale Terrain in Unity
        **first half**
        - Scripting in JavaScript - Make a Flappy Bird game [TODO]
            - temporary: "Learn to code in javascript Khan Academy"
        - Basics of Unity - Make Flappy Bird + (some other popular game) and publish on Itch.io store [TODO]
            - temporary: "Learn to code by making games GameDev.tv"
        - CS50 - Computer Science (C) [DONE,THIRD-PARTY]
        - Minecraft Terrain Basics (C# Unity) [DONE]
        - Minecraft Worlds - Terrain Generation with Jobs (C# Unity) [ALMOST]
        - Planetary LOD Voxel Terrain [PROTOTYPE]
        
        **second half, do later**
        - Marching Cubes Voxel Terrain with Jobs (C# Unity)
        - Planetary LOD Marching Cubes Terrain (prerequisite "Planetary LOD voxel terrain" + "Marching Cubes Voxel Terrain")
    
    - Series: Massively Destructible Planet Terrain in Unity
        - Scripting in JavaScript - Make a Flappy Bird game [TODO]
            - temporary: "Learn to code in javascript Khan Academy"
        - Basics of Unity - Make Flappy Bird + (some other popular game) and publish on Itch.io store [TODO]
            - temporary: "Learn to code by making games GameDev.tv"
        - CS50 - Computer Science (C) [DONE,THIRD-PARTY]
        - Minecraft Terrain Basics (C# Unity) [DONE]
        - Minecraft Worlds - Terrain Generation with Jobs (C# Unity) [ALMOST]
        - Planetary LOD Voxel Terrain [PROTOTYPE]
        - Planet Scale Terrain Destruction
            - place + destroy blocks at a certain LOD level (and all the complexities of that)
            - simple inventory
            - collison based destruction
                - if block at low-res LOD gets edited cascade changes to high-res LODS
                - start at highest LOD level for destruction slowly cascade the destruction to hiegher res LOD levels
                - jobified/C#-tasks-based destroy blocks queue
            - throw a sphere at the planet and watch destruction!
                - particle systems

    - Series: Make Destructible Terrain in Unity
        - Scripting in JavaScript - Make a Flappy Bird game [TODO]
            - temporary: "Learn to code in javascript Khan Academy"
        - Basics of Unity - Make Flappy Bird + (some other popular game) and publish on Itch.io store [TODO]
            - temporary: "Learn to code by making games GameDev.tv"
        - CS50 - Computer Science (C) [DONE,THIRD-PARTY]
        - Minecraft Terrain Basics (C# Unity) [DONE]
        - Minecraft Worlds - Terrain Generation with Jobs (C# Unity) [ALMOST]
        - Destructible Voxel Terrain
            - digging + placing individual blocks
            - collision based destruction

    - Series: Make Minecraft
        - Scripting in JavaScript - Make a Flappy Bird game [TODO]
            - temporary: "Learn to code in javascript Khan Academy"
        - Basics of Unity - Make Flappy Bird + (some other popular game) and publish on Itch.io store [TODO]
            - temporary: "Learn to code by making games GameDev.tv"
        - CS50 - Computer Science (C) [DONE,THIRD-PARTY]
        - Minecraft Terrain Basics (C# Unity) [DONE]
        - Minecraft Worlds - Terrain Generation with Jobs (C# Unity) [ALMOST]
        - Destructible Voxel Terrain
            - digging + placing individual blocks
                - simple inventory system (like minecrafts)
            - collision based destruction
        - Minecraft Trees and Villages


    - Series: Make Minecraft
        - Scripting in JavaScript - Make a Flappy Bird game [TODO]
            - temporary: "Learn to code in javascript Khan Academy"
        - Basics of Unity - Make Flappy Bird + (some other popular game) and publish on Itch.io store [TODO]
            - temporary: "Learn to code by making games GameDev.tv"
        - CS50 - Computer Science (C) [DONE,THIRD-PARTY]
        - Minecraft Terrain Basics (C# Unity) [DONE]
        - Minecraft Worlds - Terrain Generation with Jobs (C# Unity) [ALMOST]
        - Destructible Voxel Terrain
            - digging + placing individual blocks
            - collision based destruction
        - Minecraft Trees and Villages
    
    - Series: Make a game engine with Zig (long term, less revenue, definitely a niche)
        - Scripting in JavaScript - Make Flappy Bird
        - CS50 - Computer Science (C)
        - Rendering and 3D math - Make a Renderer (Processing.js)
        - Entity Component System - Managing Game Data (C)
        - Harnessing the Graphics Card - Make a Renderer (Processing.js then C)
        - Physics and 3D math - Make a Physics system (C)

3. 
    Humans seek out other humans with similar interests. They are also more likely to succeed if they are part of a community of people who can help them and encourage them.
    
    To solve this Udemy has a "Q&A" section for each course. **They've attacked the problem in the most direct way possible.** But what are people using even more than the Q&A system? Discord servers! With Discord you can make more than just a Q&A section, you can organize your community, make announcements, threads, showcase channels. Discord is a place for your student to come together as a community. Youtube is also probably the most exciting place when it comes to discovering cool things programmers have built. So what we are we going to do? Each course is going to be like it's own dedicated discord server. We aren't just selling courses, we are offering the whole social experience, online.

    **Marketing:**

    Marketing is selling a good to saisfy a need. People need a place to ask questions and get help. They also have a need for recogniition. They also will create or find a community to satisfy conversing with poeple with similar interests.
    Each course should display what students are currently online and which ones aren't online (ege: like discord). This does something important. It shows you aren't alone in what your doing. This gives a sense of not being alone. And its encouraging to know other people are working on the same thing you are.
4. 
    Posts system like medium. Not as important but the purpose is for course progress to not just be a checkmark, but to help you create an online persona and portfolio that you can use to sell youself to others and show them what you are capable of. Also it's just so cool to see what people have done. 
    
    I'm think of it like a bunch of mini master thesis'.
    
    I'm still stuck on how to keep continuity between the posts system and chat system. But they should each have a purpose.

    Video uploading? Best thing I can think of: In broswer recording tool? after you record it is stored in the cloud and you can share that?

    **Marketing:**

    I have no idea. But first we need to figure out why students would want to make a post vs use our chat system. Somehow I need to change it from "sparker needs ppl to make posts" to "students want/need to post because..." 

    **Marketing Ideas:**
    I say it's best used as a method of showcasing progress, and users who want to indirectly add tutorials or information they learned to a course. Its also a more permanent way of showcasing work/sharing thoughts than a fast-paced chat system where conversations get forgotten (don't delete this).

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