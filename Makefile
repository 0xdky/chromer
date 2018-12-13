# -*-Makefile-*-
# Time-stamp: <2018-11-05 11:37:59 dky>

BUILDROOT=build
BUILDTIME=$(shell date +%d_%b_%Y)

all: $(BUILDROOT)/Chromer.app

$(BUILDROOT)/Chromer.app: chromer Info.plist
	cp Info.plist $(BUILDROOT)/Chromer.app/Contents/Info.plist

chromer: $(BUILDROOT)/Chromer.app/Contents/MacOS/chromer

$(BUILDROOT)/Chromer.app/Contents/MacOS/chromer: *.go *.h *.m Makefile frameworks
	mkdir -p $(BUILDROOT)/Chromer.app/Contents/MacOS
	go build -i -o $(BUILDROOT)/Chromer.app/Contents/MacOS/chromer

frameworks: $(BUILDROOT)/.frameworks

$(BUILDROOT)/.frameworks: macosx-frameworks
	mkdir -p $(BUILDROOT)
	macosx-frameworks $(BUILDROOT)

.Phony: install
install: $(BUILDROOT)/Chromer.app
	pkill chromer; echo
	rm -fr /Applications/Chromer.app
	cp -rf $(BUILDROOT)/Chromer.app /Applications/.
	open /Applications/Chromer.app

zip: $(BUILDROOT)/Chromer.app
	cd build && zip -r -9 ../Chromer.app.$(BUILDTIME).zip Chromer.app

.PHONY: clean
clean:
	rm -rf $(BUILDROOT)/Chromer.app

.PHONY: distclean
distclean: clean
	pkill chromer; echo
	rm -rf /Applications/Chromer.app
	rm -rf $(BUILDROOT)
