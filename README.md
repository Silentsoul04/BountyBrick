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

## Note

I do not reccomend trying to run it, I haven't put tests in place and to be honest I'm not even sure
the `setup.sh` works... I mean it does but it might break so be prepared to do some manual setup.
