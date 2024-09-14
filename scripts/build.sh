#! /bin/sh

set -e

# create a frontend build

# copy the build to backend public folder, make the folder if not present

# invoke the docker build command

# tagging the image??

if [ $# -lt 2 ]; then
    echo "Correct Input Order <<ENVIRONMENT>> <<TAG REF>>"
    exit 1
fi

echo "Building Frontend"
cd ./frontend
yarn install --frozen-lockfile
if [ -d "./build" ]; then
    rm -r build
fi
echo "building $1 build"
if [ $1 == "production" ]; then
    yarn build:prod
else
    yarn build
fi
cd ..

echo "Copying frontend build to backend"
rm -R ./backend/public
if [ -d "./backend/public" ]; then
    echo "backend public exists"
    cp -R ./frontend/build/* ./backend/public
else
    echo "backend public not exists"
    mkdir ./backend/public && cp -R ./frontend/build/* ./backend/public
fi

echo "Building backend"
cd ./backend
yarn install --frozen-lockfile
docker build . -t neo73/trekyourworld:$2 --target $1
docker tag neo73/trekyourworld:$2 neo73/trekyourworld:latest

echo "Pushing images to docker registry"
docker push neo73/trekyourworld:$2
docker push neo73/trekyourworld:latest