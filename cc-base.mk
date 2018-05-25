add-mk-include:
	git subtree add --prefix mk-include git@github.com:confluentinc/cc-mk-include.git master --squash

update-mk-include:
	git subtree pull --prefix mk-include git@github.com:confluentinc/cc-mk-include.git master --squash

bats:
	find . -name *.bats -exec bats {} \;
