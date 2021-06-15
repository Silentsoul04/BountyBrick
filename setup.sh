#!/bin/bash

dbName="bountybrick"
echo "Enter the username for the admin user of the db $dbName"
read username

rootPass=$(echo -n $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '' | base64 </dev/urandom | head -c 36 ; echo '') | md5sum | awk '{print $1}')
userPass=$(echo -n $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '' | base64 </dev/urandom | head -c 36 ; echo '') | md5sum | awk '{print $1}')

secretKey=$(echo -n $(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '' | base64 </dev/urandom | head -c 36 ; echo '') | md5sum | awk '{print $1}')

echo "MONGO_URI=mongodb://$username:$userPass@mongo:27017/$dbName" > api/.env
echo "SECRET_KEY=$secretKey" >> api/.env
echo "MAGIC_LINK=https://firebounty.com/?sort=created_at&order=desc&reward=Gift&reward=Reward&type=Bounty&search_field=scopes&search=github" >> api/.env
echo "ROOT_LINK=https://firebounty.com" >> api/.env
echo "GITHUB_API=https://api.github.com/" >> api/.env
echo "DEBRICKED_API=https://app.debricked.com/api/1.0/open/" >> api/.env

sudo docker-compose up -d

sleep 5

echo "COPY THESE INTO THE MONGODB SHELL"
echo "1. use admin"
echo "2. db.createUser({user: \"root\", pwd: \"$rootPass\", roles:[\"root\"]});"
echo "3. use $dbName"
echo "4. db.createUser({user: \"$username\", pwd: \"$userPass\", roles:[{role: \"readWrite\", db: \"$dbName\"}]});"
echo "5. exit"
read
mongo

echo "COPY THESE INTO THE MONGODB SHELL"
echo "1. use $dbName"
echo "2. db.createCollection(\"users\")"
echo "3. db.createCollection(\"programs\")"
echo "4. db.createCollection(\"repos\")"
echo "5. exit"
read
mongo -u $username -p $userPass $dbName

echo "        command: [--auth]" >> docker-compose.yml
sudo docker-compose down
sudo docker-compose up
