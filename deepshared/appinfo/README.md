# App Info

Appinfo is used to store and retrive app related info.
Such as: scheme, download url, package name etc.


## Structure
The producer of Appinfo data should be dashboard, when user register or modify their app.
The consumer of Appinfo data should be JS server, when share link is clicked.
A template html will return to user with the appinfo data injected.

TODO
We may should separate appinfo to an independent service in future.