# prompt
PROMPT="%F{cyan}%D%f %F{cyan}%*%f %B%F{magenta}%n%f%b %B%F{blue}%~%f%b %F{white}%?%f 🐌 "

# git functions
alias g='git'
alias gb='git branch'
gs() { git status $@ }
gl() { git log --graph --pretty=bayou --abbrev-commit $@ }
ga() { git add $@; gs $@ }
gc() { git commit -S $@ }
gp() { git pull --rebase $@ }

