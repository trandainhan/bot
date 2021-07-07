#!/bin/sh

cd "$(dirname "$0")"

echo "================="
echo "Tag and push new tag to production to buil new image"
echo "Make sure you have latest code ( pull master ) and latest tag ( pull --tags )"
echo "================="

read -r -p "are you standing at master branch: (Y/n)" atMaster
case $atMaster in
  [yY] ) isAtMaster=true;;
  [nN] ) isAtMaster=false;;
esac

if [ "$isAtMaster" != true ]; then
    exit 0
fi

echo "Current version is:" $(git describe --abbrev=0 --tags)

echo "Which version do you want to tag, please remember follow version convention, e.g v1.2.3"
read -p "Which version?: " version

echo "Tagging new verison"
git tag $version

echo "Pushing new tag to gitlab"
git push --tags

echo "Done push new tag to gitlab, wait for the build success: https://gitlab.com/fiahub/coingiatot/-/pipelines"

read -p "Press anykey after the build is done" waiting

read -p "Type production to deploy: " confirm

if [ "$confirm" != production ]; then
    exit 0
fi
echo "Deploying to production"

echo "Apply new config and secrets if there are changes"
kubectl apply -f ./deployments/config/configmap.yaml

./deployments/app/update_version_and_deploy.sh
