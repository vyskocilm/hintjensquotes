
.PHONY: hintjens.com
hintjens.com:
	cd hintjens.com/ && \
	wget \
		--wait=2 \
		--limit-rate=20K \
		--recursive \
		--page-requisites \
		--no-parent \
		--convert-links \
		--adjust-extension \
		http://hintjens.com
