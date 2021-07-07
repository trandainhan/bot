#!/bin/sh

cd "$(dirname "$0")"

new_version=$(env -i git describe --abbrev=0 --tags)

echo "Deploy Bot to version $new_version"

for file in *.yaml
do
    sed -i -e "s/VERSION/$new_version/g" "$file"
done

kubectl apply -f ./

# Roll version back to placeholder
for file in *.yaml
do
    sed -i -e "s/$new_version/VERSION/g" "$file"
done
