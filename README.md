## Tempus 

My TUI to-do manager for WebDAV

![](./.media/new.gif)

### Info

Uses TODO instead of calendar events from WebDAV (e.g. Nextcloud). Code can be reporpused to sync/consume .isc files.

Password is stored in desktop secret-services (should be compatable with MacOS keyring)

Currently shows today tasks

# Intallation

`go install github.com/AbandonwareDev/tempus@latest`

## TODOs

BUGs:
 - TODO won't appear if its time set to 0:00 of today (other apps may do that to make notifications)
 - Repetative tasks - won't complete (instead destroys them???), can't create, no indication 
 - not always shows full description?
 - while Adding TODO - hitting 'q' may(?) close program (and destroy new TODO) (???)
 - tempus may consume 100% cpu core for unknown reason (day passed // system sleep // network change/issues // database locked/inaccessable but we need creds // ???)

Packaging: make nixos flake building and installing easy

todo:
 - in new TODO date chooser - add days of week representing closest one (and Sun1 - through 1 sunday)
 - check TODOs in code
 - TODO creation/overview - can't add multiline description?
 - recently completed/edited search/list - to return accidentally completed task
 - Add task searcher, next week tasks, today+tomorrow, missed tasks (past week, all missed)
