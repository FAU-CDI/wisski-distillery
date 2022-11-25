# News

This file contains signficant news items for the distillery.

# Automatic Password Checking (2022-11-25)
- Implemented automatic password checking

# Login using Distillery Administration (2022-11-23)
- The admin interface now allows login to individual user accounts

# Showing Statistics (2022-11-16)
- The distillery nows shows generic statistics on the public homepage 
- detailed statistics can be found on the admin interface

# Refactored SSH Support (2022-11-12)
- Fully refactored ssh for users to use a real OpenSSH Server along with a small custom proxy in between
- It is now possible for developers to directly use e.g. [VSCode Remote SSH](https://code.visualstudio.com/docs/remote/ssh) to develop new WissKI features directly on the distillery.
- No configuration beyond regular ssh access is neccessary

# Migration to Traefik & Support for HTTP3 (2022-10-12)
- We have migrated the entry point from nginx to [traefik](https://traefik.io/traefik/)
- This enables much cleaner support for automatically fetching and renewing SSL certificates 
- It is now possible to turn on http3 support for the entire distillery

# Addition of a Global Distillery Resolver (2022-10-05)
- We have added a global WissKI Resolver, that functions similarly to the WissKI resolver under `/wisski/`.
- It can resolve WissKI URIs for the entire distillery, and redirect to their view page.
- Can be called exactly like the `/wisski/get?uri=` route of individual WissKIs, but toplevel on the distillery.

# Addition of an administrative server (2022-09-09)
- We have added a new web route under `/dis/`
- Allows administrators to manage WissKI Instances and see their status
- Administrators can e.g. download pathbuilders and make backups and snapshots
- At this point it is only for administrators
- A future (public) server with statistics will follow

# Migration to go (2022-09-08)
- We have ported the distillery from a set of bash scripts to a self-contained go executable.
- This makes future development easier, and allows us to develop new features easier.