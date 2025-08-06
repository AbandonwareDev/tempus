## Tempus 

My TUI to-do manager (with support to sync .isc files from WebDAV)

## TODOs

BUGs:
 - while Adding TODO - hitting 'q' will close program (and destroy new TODO)
 - TODO won't appear if its time set to 0:00 of today (other apps may do that to make notifications)
 - Repetative tasks - won't complete (instead destroys them???), can't create, no indication 
 - tempus may consume 100% cpu core for unknown reason (day passed // system sleep // network change/issues // database locked/inaccessable but we need creds // ???)

Packaging: make nixos flake building and installing easy

todo:
 - in new TODO date chooser - add days of week representing closest one (and Sun1 - through 1 sunday)
 - check TODOs in code
 
