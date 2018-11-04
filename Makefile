Chromer.app: chromer Info.plist
	cp Info.plist Chromer.app/Contents/Info.plist

chromer: Chromer.app/Contents/MacOS/chromer

Chromer.app/Contents/MacOS/chromer: *.go *.h *.m Makefile
	mkdir -p Chromer.app/Contents/MacOS
	go build -i -o Chromer.app/Contents/MacOS/chromer

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
