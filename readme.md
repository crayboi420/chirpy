# Chirpy

## Introduction

This is a project where I attempt to build a webserver using go. The design is supposed to be similar to twitter with users who can post `chirps' that are less than 140 characters long.

 I learnt many things like authentication, authorization, data storage and how to build RESTful APIs during this project. I also learnt about webhooks and deploying my app to docker in this project.

 ## Usage

Firstly, we need a .env file with some environment variables needed.

```
JWT_SECRET (A randomly generated key for generating JWT tokens)
POLKA_KEY (An api authentication key for using the polka webhook)
```

Then, change the ```PORT``` which is 8080 by default and the path to the database which is ```./database.json``` by default to your choice.

 Finally build and deploy the app with:
 ```bash
go build . && ./chirpy
```

### Endpoints

The app has many endpoints for you to explore. The first thing to do is make a user. Use:
```url
POST http://localhost:PORT/api/users
```
With the structure:
```json
{"email": user's email,
"password": user's password,}
```
Then you can login with ```localhost:PORT/api/login``` and create chirps.
