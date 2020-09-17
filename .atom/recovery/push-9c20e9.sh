# create a temporary named pipe
PIPE=$(mktemp -u)
mkfifo $PIPE
# attach it to file descriptor 3
exec 3<>$PIPE
# unlink the named pipe
rm $PIPE
# extract working-branch name
BRANCH=$(git branch | grep "*" | cut -d" " -f2)
# push branch to origin and pipe output into descriptor 3
git push -u origin $BRANCH >&3

# read every line of the output
while read line; do
  # attempt to find a match for the remote-message containing the MR-URL
  MATCH=$(echo ${line} | grep "remote:   http" | cut -d" " -f2)
  if [ -z "$MATCH"]
  then
    open -a "Google Chrome" $MATCH
  fi
  echo "${line}"
done <&3
