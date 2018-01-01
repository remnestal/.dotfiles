PACKAGES = apm list --installed --bare | grep '^[^@]\+' -o

.PHONY: version
version:
	apm --version

.PHONY: install
install:
	apm install --packages-file packages.txt

.PHONY: save
save:
	$(PACKAGES) > packages.txt

.PHONY: clean
clean:
	$(PACKAGES) > uninstalled-packages-$$(date +%s).txt
	apm uninstall .
