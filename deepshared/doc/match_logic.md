# How (should) we do matching

## Background
When one of our shorturl (mainly include token and appid) is clicked on the receiver device inside some browser, two things could happen:

1. We will try to start app on device with the token directly. When app is installed on the device, token will be passed directly to app, which can then use token to retrieve context directly. This is 100% accurate.

2. If failed to open app in some fixed time frame (say 1 second), we will send a http/restful request for binding with the token. Our server will remember the association between user agent we captured from browser, and token supplied in the request. We then send user to app store for installtion, and after installation, app will call our match service with information that help to identify the device, including os and its version, brand (maker, and series). These information, along with ip forms the basis for matching.  

## Matching today
Matching today is fairly simple. When serving bind request, we take information from user agent, mostly include ip, os, os_version, brand(very useful for android), and associate the token with such info. Current implementation basically use concatanation of ip:os:os_version:brand as key, and token as value, and save these in redis during the bind. When match request comes, we again collect info and use ip:os:os_version:brand as key to retrieve token that is associate with it.

Detail:
When multiple device connects to us behind single ip, we might get into trouble. Currently, to reduce the possible confusion, we have a special logic where we try to use time to make match decision. In particular, every time there is a match, we attach the receiver device id (some hashed form), with ip:os:os_version:brand entry. If the same entry is matched again within 15 minutes, we do not return the token. This strategy should be studied and improved, when we have enough time.

Previous implementation also has a bug in match. When match is called with token, the same 15 minutes switch is also applied. This is wrong. In new restful api, there are two kind of matching: /match/match, or /shorturl/<appID>/<token>, where the second match simply always return the context associated with token.
