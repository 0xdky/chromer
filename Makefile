Chromer.app: *.go *.h *.m Makefile Info.plist
	mkdir -p Chromer.app/Contents/MacOS
	go build -i -o Chromer.app/Contents/MacOS/chromer
	cp Info.plist Chromer.app/Contents/Info.plist

install: Chromer.app
	pkill chromer; echo
	rm -fr /Applications/Chromer.app
	cp -rf Chromer.app /Applications/.
	open /Applications/Chromer.app

.PHONY: open
open: Chromer.app
	open Chromer.app

.PHONY: scheme
scheme: Chromer.app
	open https://www.gnu.org/

.PHONY: clean
clean:
	rm -rf Chromer.app

distclean: clean
	pkill chromer; echo
	rm -rf /Applications/Chromer.app
