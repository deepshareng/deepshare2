# Generate URL

URL generator is used to generate a share link with all data attached to it.

## API.GenURLPostBody
All data that the user can provide

## API.ShareLink_Info
The data related to the share link, when share link is clicked, these data would be used.
These data should be handled sharelink front, and should not send to match service.

## API.Context
The data that user want to use for their context restore.

## API.IsShort
This is to determine if the returned as a short url.

If true, the returned url should be like this:
##### short url format:
http://fds.so/u/7713337217A6E150/:seg

If false, the returned url should be like this:
##### raw Url format
http://fds.so/u/7713337217A6E150?download_title=aaa&download_msg=bbb&redirect_url=ccc&inapp_data=y

## API.SenderID
If this field is set, it means the action of generating URL is from Sender part.

## API.ForwardedSenderID
If this field is set, it means the action of generating URL is from Receiver part.
