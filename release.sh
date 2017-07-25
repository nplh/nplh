old_version=$(curl -s https://api.github.com/repos/nplh/nplh/releases/latest | jq -r ".tag_name")
echo "New Version:"
read version
sed -i "s/$old_version/$version/g" nplh.go
sed -i "s/$old_version/$version/g" install.sh

git commit -am "Release $version"
git push origin master

git tag -a $version
git push origin $version
