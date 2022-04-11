# Openpilot Mod Manager
[![Test](https://github.com/borgmon/openpilot-mod-manager/actions/workflows/go.yml/badge.svg)](https://github.com/borgmon/openpilot-mod-manager/actions/workflows/go.yml)

Openpilot Mod Manager(OMM) is an Openpilot modding tool for both users and developers. It manipulates source code directly as a text file(instead of using language features) in order to customize behaviors.


## Install
Even tho OMM works for every forks, it is recommended to use official branch for dependency tracking.
```sh
git clone https://github.com/commaai/openpilot -b v0.8.13
```
 
Download and install OMM
```sh
curl -fsSL https://github.com/borgmon/openpilot-mod-manager/raw/main/bin/install | sh
```

## Quick Start
### Initialize Repo
Go to your openpilot folder and type

```sh
omm init
```
This command will prepare the repo for modding and generate an `omm.yml` file

### Install Mod
```sh
omm install https://github.com/borgmon/omm-no-disengage_on_gas
```
This mod does not exist yet :P

You can checkout all the commands via `omm --help` or the [Wiki](https://github.com/borgmon/openpilot-mod-manager/wiki/2.-Commands).

## Developer guide
Developing a mod is very easy. Let's start by making a new directory and type
```sh
omm mod init
```
This will generate your manifest file `manifest.yml`. This file must be placed in root of your folder
```yml
name: omm-my-mod # same as github project name
displayName: my mod
repoURL: https://github.com/myname/omm-my-mod
version: v0.8.13-1.0
publisher: my name
description: my first mod!
dependencies: # planned. not being used right now
  - url: https://github.com/borgmon/some-UI-framework
    version: v0.8.13-18
```


Full guide is available in the [Wiki](https://github.com/borgmon/openpilot-mod-manager/wiki/3.-Developer-Guide).