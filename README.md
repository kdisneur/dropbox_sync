# Dropbox Sync

> :warning: this application add/update/remove files and folders on dropbox and
> on your file system. Use at your own risk.

Synchronize a Dropbox and local folder, each time a modification happens.

## Install

1. Create a [Dropbox OAuth2 application][DROPBOX_OAUTH_DOC]
2. Create a configuration file in `~/.config/dropbox_sync/config`

```toml
[authentication]
client_id = "<dropbox_client_id>"
client_secret = "<dropbox_client_secret>"

[[folder]]
remote_path = "/path/to/dropbox/folder"
local_path = "~/Documents/here"

[[folder]]
remote_path = "/path/to/dropbox/another/folder"
local_path = "~/Documents/somewhere/else"
```

## Usage

```
Usage of dropbox_sync:
      --debug     enable debug logging
  -h, --help      show the current message
  -v, --version   show version number
```

[DROPBOX_OAUTH_DOC]: https://www.dropbox.com/developers/reference/oauth-guide
