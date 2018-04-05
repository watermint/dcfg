# About DCFG

[![Build Status](https://travis-ci.org/watermint/dcfg.svg?branch=master)](https://travis-ci.org/watermint/dcfg)
[![Coverage Status](https://coveralls.io/repos/github/watermint/dcfg/badge.svg?branch=master)](https://coveralls.io/github/watermint/dcfg?branch=master)

DCFG allows syncing between Google Apps and Dropbox Business. DCFG automates provisioning and deprovisioning of both users and groups. DCFG reads users and groups on Google Apps, then create/update/delete Dropbox accounts and groups.

# Feature

## user-provision

Invite users to Dropbox Business team who are in Google Apps but NOT in Dropbox. Users are identified by primary email address.

## user-deprovision

Delete users from Dropbox Business team who are NOT in Google Apps but in Dropbox. Users are identified by primary email address.

## group-provision

Create/update Dropbox Groups by refering Google Group. Google Group can have nested groups, DCFG expands all nested group, then add all members to Dropbox Group.
In some usecase, Google Groups have tons of unmanaged groups. Most IT admins don't want to sync all group into Dropbox. DCFG requires white list of Google Groups (group emali address).

### example

Sample structure of Google Groups.
```
japan@example.com (Name: Japan)
  |
  +- taro@example.com (USER)
  +- tokyo@example.com (GROUP)
  
tokyo@example.com (Name: Tokyo)
  |
  +- kevin@example.com (USER)
```

If IT admin add `japan@example.com` to the white list, DCFG creates Dropbox Group named "Japan". And the group has both `taro@example.com` and `kevin@example.com`.

# Requirements

* Administrator priviledges required for both Google Apps and Dropbox Business.

# How to use: Preparation

## Prepare directory

1. DCFG requires directory for configuration files and writing logs. Create *DCFG directory* somewhere on your machine.
2. Download DCFG binary into *DCFG directory*

## Create Project on Google APIs

1. Login to Google Apps as administrator.
2. Create project on [Google APIs Portal](https://console.developers.google.com/iam-admin/projects).
3. Enable [Admin SDK](https://console.developers.google.com/apis/api/admin/overview).
4. Create credentials. Select "Admin SDK" for "Which API are you using?", "Other UI (e.g. Windows, CLI tool)" for "Where will you be calling the API from?". Choose "User data" for "What data will you be accessing?".
5. Create Client ID.
6. Then download client secret JSON file.
7. Rename `client_secret_*.json` to `google_client_secret.json`
8. Move `google_client_secret.json` to *DCFG directory*.

## Authorise and store token of Google Apps

1. `dcfg -path *DCFG directory* -auth google 
2. Open link, which displayed by above command.
3. Approve and copy code.
4. Paste code into dcfg

## Create App on Dropbox Business

1. Login to Dropbox Business as team admin.
2. Create app on [My Apps](https://www.dropbox.com/developers/apps).
3. Choose `Dropbox Business API`
4. Choose `Team member management` for access type.
5. Name your app.
6. Generate access token by clicking "Generate" button.

## Store token of Dropbox Business

1. `dcfg -path *DCFG directory* -auth dropbox`
2. Paste generated token

## Create Google Group white list (optional)

Create list of Google Group emails. Format is like below:

```
japan@example.com
singapore@example.com
```

# How to use: Provisioning, deprovisioning

## Dryrun

DCFG runs as dryrun by default. If you don't need to sync groups, `-group-provision-list *white list file*` is not required.

```
dcfg -path *DCFG directory* -group-provision-list *white list file* -sync user-provision,group-provision,user-deprovision
```

## Run

add option `-dryrun=false`

# Build

```bash
$ docker build -t dcfg . && rm -fr /tmp/dist && docker run -v /tmp/dist:/dist:rw --rm dcfg
```

# LICENSE

Copyright (c) 2016 Takayuki Okazaki

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
