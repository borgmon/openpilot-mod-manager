# Openpilot Mod Manager
Openpilot Mod Manager(OMM) is an Openpilot modding tool for both users and developers. It manipulates source code directly as a text file(instead of using language features) in order to customize behaviors.
 
## Install
Make sure you are on the official fork release branch.
```sh
git clone https://github.com/commaai/openpilot -b v0.8.13
```
 
Download and install OMM
```sh
sh -c "$(curl -fsSL {TBA})"
```
## Usage
### For users
Most commands you can get from `omm --help`. Here are some tips
 
#### Initialize Repo
This command will prepare the repo for modding and generate an `omm.yml` file in `~/.omm`, we will talk about it later.
```sh
omm init
```
 
#### Install Mod
```sh
omm install https://github.com/borgmon/omm-no-disengage_on_gas
```
 
You can also install from a local directory
```sh
omm install /home/usr/my-awesome-mod
```
 
#### Reset
```sh
omm reset
```
This will uninstall all mods, including
 
**For list of mods please check WIKI**
#### omm.yml
This is the core of the modding process. You will find one like this:
```yml
opVersion: v0.8.13
mods:
  - url: https://github.com/borgmon/omm-no-disengage_on_gas
    version: v0.8.13-1
```
As you can see, it records what's your base branch, and all your mods.
You can make changes by using commands, or edit this file directly. After your change, just run
 
```sh
omm apply
```
 
#### Load order
Some mods may require higher load order(dependency for other mods). This can be achieved by putting it higher up in the list.
 
### For developers
#### Step 1: Create manifest
Developing a mod is very easy. Let's start by making a new directory and type
```sh
omm mod init
```
 
This will generate your manifest file `manifest.yml`. This file must be placed in root of your folder
```yml
name: omm-my-mod # same as github project name
displayName: my mod
repoURL: https://github.com/myname/omm-my-mod
version: v0.8.13-1
publisher: my name
description: my first mod!
dependencies: # planned. not in use right now
  - url: https://github.com/borgmon/some-UI-framework
    version: v0.8.13-18
```
 
#### Step 2: Adding patches
MMO works by injecting your code block into the source code.
For example, if you want to add a new param to `selfdrive/common/params.cc`, find the line number you want to append to. For example [line 177 line](https://github.com/commaai/openpilot/blob/v0.8.13/selfdrive/common/params.cc#L177)
 
Create the patch file by run
```sh
omm mod new selfdrive/common/params.cc
```
This will create an empty file with that path. Open newly created file and enter:
 
```c++
// >>>#177
    {"MyNewParam", PERSISTENT},
```
The `>>>` will tell omm which line to append after. Pretty neat yeah? You can add multiple code blocks like that in the same file. Also noticed the indentation? MMO will append the code block as-is so make sure it's right, especially for python scripts.(for python the comment token is `#`, omm only cares about operands)

##### Operands
currently 2 operands are supported:
> `>>>` is for appending code block. Appending order is the load order.

> `<<<` is for replacing code block. If multiple mods use replace operands on the same line number will cause **conflict**! So use this cautiously!

Keep repeating step 2 until you finish all your changes.
 
#### Step 3: Testing
Install your mod by go to openpilot dir and type
```sh
omm install {absolute path to your local mod}
```
MMO will create a new branch with the name `omm-timestemp`, and commit all your changes here.
 
Here you can test your mod and see if it's working.
 
#### Step 4: Publish
When you push your mod repo to github, you can publish your code by adding tags.
 
It has to be in this format:
```
{openpilot_version}-{mod_version}
```
for example `v0.8.13-1`.
 
And you are done! Share your mod url and let people try it out!

## Q&A
### What is `omm-core`?
`omm-core` is a build-in mod and will be downloaded when you run `oom init`. This is required to setup the repo and should always be placed on the top. This also allow omm run independently across all versions.

## To-Do
- [ ] dependency management
- [ ] auto reset on system error
- [ ] `omm mod test` to test mod
- [ ] a desktop wrapper similar to [OpenpilotToolkit](https://github.com/spektor56/OpenpilotToolkit), allowing user ssh into devices and manage mod remotely with UI
- [ ] use worker pool for injection and mod loading
- [ ] upgrades
- [ ] on device options

