## Description ##

Using a work laptop and keeping work and personal workflows is a problem.
Chromer is a little helper to open URLs in Google Chrome browser in specific
profile based on pattern matching the URL.

This is a Mac OS specific solution but can easily be extended to any other
platform. The Mac OS specifics are the round-about approach in getting the
clicked URL on Mac.

You create a config file and specify the profile and patterns to match. Chromer
takes care of the rest. Chromer can handler 'http', 'https' and 'file' URL
schemes.

## Getting it running ##

### Build & Installation ###

* Build Chromer (will provide a pre-built version soon) by checking out this
  repository and `make` (Needs Go, Xcode + MacOS SDK)
* `make install`
* Assuming everything goes fine, you should find '/Applications/Chromer.app'
  and a small icon in the Mac dock
* Go to 'System Preferences', 'General' and 'Default Web Browser' drop down.
  Select 'Chromer' as the default browser.


### TODO ###

* Add an event listener to handler application quite event. My Objective-C
  skills are non-existent, will take time

### Bugs ###

* Chromer does not exit when you try to shut down the laptop and you get a
  popup message. You need to open the terminal and kill the 'chromer' process.
  Ex: `pkill -9 chromer`

### Disclaimer ###
	I have built and tested on a single machine, your mileage might vary. I
	will be happy to help if required.
