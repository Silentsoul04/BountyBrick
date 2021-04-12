# BountyBrick

This project aims to create an internal tool which will use the Debricked platform to scan
all the avaiable open-source bug bounty programs, eventually providing a tool which assists the
security researcher to find dependency-related vulnerabilities in massive projects.

## Progress

- [x] Scrape firebounty.com for programs
- [x] Extract and save programs and repositories in database
- [x] Extensive information about repositories collected
- [x] Serve all the information collected
- [x] Able to fork and remove repos on github profile
- [x] Find most optimal way to integrate Debricked API for uploads
- [ ] Subscription system (Debricked API) to scanned repositories -> update database
- [ ] Add more complex regex system to match more repo urls
- [ ] Add sorting based on different filters for programs and repositories
- [ ] Build frontend with Vue

## Setup 

You can use the `setup.sh` script to setup an instance of this api (although some of the steps still have to be done manually). However I recommend for now to look through it before you run it so that it doesn't break anything. In the future I'll make a more effective setup system along with tests.

`api/.env` This file contains all the important variables for the API to work:
```
MONGO_URI=mongodb://splinter:password@mongo:27017/bountybrick
SECRET_KEY=0bded16c51a7809622a195f91895f175

MAGIC_LINK=https://firebounty.com/?sort=created_at&order=desc&reward=Gift&reward=Reward&type=Bounty&search_field=scopes&search=github
ROOT_LINK=https://firebounty.com

GITHUB_API=https://api.github.com/
GITHUB_OAUTH=
GITHUB_ORG=

DEBRICKED_USER=
DEBRICKED_PASS=
DEBRICKED_API=https://app.debricked.com/api/1.0/open/
```

The setup script will take care of all of them except `GITHUB_OAUTH` (OAuth token which has rights the organisation used), `GITHUB_ORG` which is the name of the organization, `DEBRICKED_USER` which is the username (email) used of the debricked account and `DEBRICKED_PASS` which is the password of that account

## API Docs

#### Open Endpoints :earth_americas:

:arrow_forward: `POST /api/login`

Provided username and password as json-data it will return a JWT token

**Valid request:**
```json
{
    "username":"splinter",
    "password":"superpass"
}
```
:heavy_check_mark: **Successful response:**
```json
{
    "message": "success",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InNwbGludGVyIiwiUm9sZSI6InJvb3QiLCJleHAiOjE2MTc5MjEwMzF9.mQ9sd23Bm95kp2rXgMsPj41zAV2h-1GSYmHv8glstRk",
    "username": "splinter"
}
```
:x: **Unsuccessful reponse:**
```json
{
    "message": "Unauthorized."
}
```
#### Endpoints that require token :lock:
All of these requests require a custom header `Token:` + the value of the token obtained at the login endpoint.

:arrow_right: `GET /api/programs`

Returns all the programs in the database

:heavy_check_mark: **Sample response:**
```json
{
    "message": "success",
    "programs": [
        {
            "ID": "606e354439411b67f90222ee",
            "name": "Diem",
            "link": "https://hackerone.com/diem",
            "repos": [
                "606e354539411b67f90222ef"
            ],
            "created": "2021-04-07T22:42:13.125Z",
            "updated": "0001-01-01T00:00:00Z"
        },
    ...
    ]
}
```

:arrow_right: `GET /api/programs/:id`

Returns data about a specific program identified using the `id` parameter

:heavy_check_mark: **Successful response:**
```json
{
    "message": "success",
    "program": {
        "ID": "606e354439411b67f90222ee",
        "name": "Diem",
        "link": "https://hackerone.com/diem",
        "repos": [
            "606e354539411b67f90222ef"
        ],
        "created": "2021-04-07T22:42:13.125Z",
        "updated": "0001-01-01T00:00:00Z"
    }
}
```
:x: **Unsuccessful response**
```json
{
    "message": "No program with id: 606dd68b13e44bc0fea44658 found!"
}
```

:arrow_right: `GET /api/repos/`

Returns all the repositories in the database

:heavy_check_mark: **Sample response:**
```json
{
    "message": "success",
    "programs": [
        {
            "ID": "606e354539411b67f90222f1",
            "name": "bootgen",
            "link": "https://github.com/Xilinx/bootgen",
            "short": "Xilinx/bootgen",
            "brick": "",
            "program": "606e354439411b67f90222ec",
            "program_name": "Xilinx BBP",
            "forked": false,
            "git_forks": 19,
            "git_stars": 11,
            "size": 929,
            "created": "2021-04-07T22:42:13.18Z",
            "updated": "0001-01-01T00:00:00Z"
        },
        ...
    ]
}
```

:arrow_right: `GET /api/repos/:id`

Returns data about a specific repository identified using the `id` parameter

:heavy_check_mark: **Successful response:**
```json
{
    "message": "success",
    "repository": {
        "ID": "606e354539411b67f90222f3",
        "name": "mattermost-server",
        "link": "https://github.com/mattermost/mattermost-server",
        "short": "mattermost/mattermost-server",
        "brick": "",
        "program": "606e354439411b67f90222eb",
        "program_name": "Mattermost",
        "forked": false,
        "git_forks": 4862,
        "git_stars": 20154,
        "size": 428939,
        "created": "2021-04-07T22:42:13.186Z",
        "updated": "0001-01-01T00:00:00Z"
    }
}
```
:x: **Unsuccessful response**
```json
{
    "message": "No repository with id: 606e354539411b67f9022246 found!"
}
```

:arrow_forward: `POST /api/repos/:action`

Executes an action from the `action` parameter on the repositories (id) provided as json data.

**Valid request:**
```json
{
    "repos":[
        "606e354539411b67f90222f1",
        "606e354539411b67f90222f3",
        "606e354539411b67f90222f2"
    ]
}
```

:heavy_check_mark: **Successful response:**
```json
{
    "message": "success",
    "repos": {
        "606e354539411b67f90222f1": "Successfully started action: remove",
        "606e354539411b67f90222f2": "Successfully started action: remove",
        "606e354539411b67f90222f3": "Successfully started action: remove"
    }
}
```
:x: **Unsuccessful response**

```json
{
    "actions": {
        "bookmark": "Bookmark repository to personal profile",
        "fork": "Fork the repository",
        "remove": "Remove repository from github page",
        "scan": "Run a Debricked scan on repository"
    },
    "message": "The action: lol isn't valid!"
}
```

:arrow_forward: `POST /api/programs/:action`

Executes an action from the `action` parameter on the programs (id) provided as json data.

**Valid request:**
```json
{
    "programs":[
        "606e354439411b67f90222ee",
        "606e354439411b67f90222ec",
        "606e354439411b67f90222ed"
    ]
}
```

:heavy_check_mark: **Successful response:**
```json
{
    "message": "success",
    "programs": {
        "606e354439411b67f90222ec": "Successfully started action: fork",
        "606e354439411b67f90222ed": "Successfully started action: fork",
        "606e354439411b67f90222ee": "Successfully started action: fork"
    }
}
```
:x: **Unsuccessful response**

```json
{
    "actions": {
        "bookmark": "Bookmark program to personal profile",
        "fork": "Fork all the repositories in program",
        "scan": "Run a Debricked scan on all the repositories in program"
    },
    "message": "The action: lol isn't valid!"
}
```

:arrow_right: `GET /api/actions`

Returns all actions avaiable for `repos` and `programs`

:heavy_check_mark: **Sample response:**
```json
{
    "message": "success",
    "actions": {
        "bookmark": "Bookmark repository to personal profile",
        "fork": "Fork the repository",
        "remove": "Remove repository from github page",
        "scan": "Run a Debricked scan on repository"
    },
    "programs": "Every action can also be executed on programs, affecting all contained repos"
}
```