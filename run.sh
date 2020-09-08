#!/bin/bash
# Current status
echo "Starting run.sh"
date
echo "Current directory:"
pwd

echo "Installing yarn"
npm install -g yarn@^1.19.1

# Install npm packages
echo "Installing required npm packages"
yarn install
echo "npm packages installed"

# Create docker env file
cat << EOF > docker.env
DBDATABASE=$DBDATABASE
DBHOST=$DBHOST
DBPASSWORD=$DBPASSWORD
DBPORT=$DBPORT
DBUSER=$DBUSER
EOF


echo "Starting api"
NODE_ENV=docker yarn run start:dev

# Keep container running
tail -f /dev/null
