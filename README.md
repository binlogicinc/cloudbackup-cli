# Binlogic CloudBackup Cli
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://github.com/cweill/gotests/blob/master/LICENSE)

---

[CloudBackup](https://www.binlogic.io/) is a tool to orchestrate MySQL, MariaDB, MongoDB and PostgreSQL Backups in the Cloud.

*cloudbackup-cli* is a command line tool to manage and automate Binlogic [CloudBackup](https://www.binlogic.io/) API.


## Installation

- If you already have go installed

```shell
  go get github.com/binlogicinc/cloudbackup-cli
```

- Another alternative is to download the binary from the [Release Section ](https://github.com/binlogicinc/cloudbackup-cli/releases).

## Usage

First and foremost, you need to go to your panel URL, click your username, go to Settings and finally to API Settings.
Here you'll be able to generate the API keys needed to interact with the control panel progragmatically.

Then you need to define three mandatory parameters for all your API calls: `access-key`, `access-secret` and `host` (this is the
same host you use to interact with our panel, like https://YOUR-COMPANY.binlogic.io).

You can pass all these as command line parameters (`--access-key=`, `--access-secret=`, `--host=`), as environment variables
(`BL_ACCESS_KEY`, `BL_ACCESS_SECRET` and `BL_host`) or in a configuration file, by default in $HOME/.cloudbackup-cli.toml, like:
```
access-key = "PUT_YOUR_ACCESS_KEY_HERE"
secret-key = "PUT_YOUR_ACCESS_SECRET_HERE"
host = "https://YOUR_COMPANY.binlogic.IO"
```

After that, you can use the built in command help to explore it's capabilities