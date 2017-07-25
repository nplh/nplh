old_version=$(curl -s https://api.github.com/repos/nplh/nplh/releases/latest | jq -r ".tag_name")
echo "New Version: (currently $old_version)"
read version
sed -i "s/$old_version/$version/g" *.*

git commit -am "Release $version"
git push origin master

git tag -a $version
git push origin $version
