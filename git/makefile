.PHONY: version
version:
	git --version

.PHONY: origin
origin:
	git config --list --show-origin

.PHONY: config
config:
	cp .gitconfig ~/.gitconfig

.PHONY: alias
alias:
	cat alias.sh >> ~/.bashrc
