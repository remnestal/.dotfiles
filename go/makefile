GOSOURCE=https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz
TARBALL=/tmp/go.tar.gz

INSTALLPATH=/usr/local

.PHONY: version
version:
	go version

.PHONY: install
install:
	wget $(GOSOURCE) -O $(TARBALL)
	sudo tar -C $(INSTALLPATH) --verbose -xzf $(TARBALL)
	echo "# Added by dotfiles/go/makefile" >> ~/.profile
	echo "export PATH=$$PATH:$(INSTALLPATH)/go/bin" >> ~/.profile
	# Run 'source ~/.profile' to reload current shell
