gnome-terminal -x sh -c "
cd /home/raz/IdeaProjects1/projects1/gas_05_2017
git add -u
git add *
git commit -m '$(date)'
git push --repo https://github.com/raizul8/gas_05_2017.git

echo '\nDone! Press Enter to close.\n'
read v
"
