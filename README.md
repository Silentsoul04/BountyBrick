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
- [ ] Find most optimal way to integrate Debricked API for uploads
- [ ] Integrate Github webhooks to keep the forks updated
- [ ] Subscription system (Debricked API) to scanned repositories -> update database
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
```

The setup script will take care of all of them except `GITHUB_OAUTH` (OAuth token which has rights the organisation used) and `GITHUB_ORG` which is the name of the organization.

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
:heavy_check_mark: Returns all the programs in the database

**Sample response:**
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

:arrow_forward: `POST /api/repos/:id?action=`
Executes an action from the `action` query parameter on a specific repository identified with the `id` parameter

:heavy_check_mark: **Successful response:**
```json
// /api/repos/606e354539411b67f90222f2?action=fork
{
    "message": "Successfully executed action: fork on skale-consensus"
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

:arrow_right: `GET /api/actions`
Returns all actions avaiable for `repos` and `programs`

:heavy_check_mark: **Sample response:**
```json
{
    "message": "Success",
    "program_actions": {
        "bookmark": "Bookmark program to personal profile",
        "fork": "Fork all the repositories in program",
        "scan": "Run a Debricked scan on all the repositories in program"
    },
    "repo_actions": {
        "bookmark": "Bookmark repository to personal profile",
        "fork": "Fork the repository",
        "remove": "Remove repository from github page",
        "scan": "Run a Debricked scan on repository"
    }
}
```