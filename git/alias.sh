
########################
# Start of git-functions
# status
gs() { git status $1; }
# add
ga() { git add $1; git status .; }
# commit
gc() { git commit $1; }
# amend commit
gca() { git commit --amend --no-edit $1; }
# pull
gp() { git pull --rebase; }
# diff
gd() { git diff --stat; }
# log
gl() { git log --max-count=10 --pretty=florida; }
# branch
gb() { git branch; }
# End of git-functions
########################
