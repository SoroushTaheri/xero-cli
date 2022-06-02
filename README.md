![GitHub Workflow Status](https://img.shields.io/github/workflow/status/SoroushTaheri/xero-cli/Release%20Go%20Binaries)
 [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Xero-CLI
The exclusive command-line interface for XeroCTF competitions.

Table of content:
- [Xero-CLI](#xero-cli)
- [Usage](#usage)
- [Installation](#installation)
  - [Method 1: Download the binary **(Recommended)**](#method-1-download-the-binary-recommended)
  - [Method 2: Build the project yourself](#method-2-build-the-project-yourself)
- [Commands](#commands)
  - [auth](#auth)
    - [login](#login)
    - [status](#status)
  - [challenge](#challenge)
    - [list (ls)](#list-ls)
    - [show](#show)
  - [submit](#submit)
  - [scoreboard](#scoreboard)
  - [rules](#rules)
  - [completion](#completion)
- [License](#license)
- [Contact](#contact)
# Usage
`xero-cli` is your main tool to be able to participate in XeroCTF competitions. You can view challenges, submit your flags or see any challenge's scoreboard.

# Installation
You can either download the CLI binary directly or you can manually clone and build the project yourself. The former is strongly preferred.

## Method 1: Download the binary **(Recommended)**
**Step 1)** Refer to downloads list and download the correct archive according to your operating system and processor architecture.

Make sure to always use the latest version of the project. Otherwise you might encounter some issues while using the CLI. 

[List of available downlads](https://github.com/SoroushTaheri/xero-cli/releases/latest)


**Step 2)** The archive only includes a single file which represents the `xero-cli` binary. You'll need to extract the archive and move the binary to a suitable directory such as `usr/local/bin`:

```
$ wget -qO- "https://github.com/SoroushTaheri/xero-cli/releases/download/v0.1.3/xero-v0.1.3-linux-amd64.tar.gz" | sudo tar xvz -C /usr/local/bin
```

You should be able to use the CLI in your shell. Do so by executing `xero` command:
```
$ xero
Command-line interface for XeroCTF 2022.

Usage:
  xero [command]

Available Commands:
  auth        Auth-related commands (login, status, ...)
  challenge   Get the full list of challenges or view details of a challenge
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  rules       View CTF rules
  scoreboard  Show the scoreboard of the competition
  submit      Here you will submit the flag you earned

Flags:
  -h, --help     help for xero

Use "xero [command] --help" for more information about a command.
```



## Method 2: Build the project yourself
Clone the project:
```
$ git clone https://github.com/SoroushTaheri/xero-cli.git
```

Navigate to project's folder and build the project:
```
$ cd xero-cli
$ go build -v -o xero
```

Move the binary to a suitable location:
```
$ sudo mv ./xero /usr/local/bin
```

Verify that you can use the CLI by executing `xero` command.

# Commands
## auth
In order to use certain commands (e.g., submissing a flag) you must sign in to your [RoboEpics](https://roboepics.com) account.

### login
Use this command to log in to your account using your RoboEpics credentials.
```
$ xero auth login
Username/Email: amghezi@gmail.com
Password: ***********
✅ Successfully logged in as: amghezi@gmail.com
```

### status
Check whether you're logged in or not. If logged in, you'll be able to see your email/username.
```
$ xero auth status
Logged in as "amghezi@gmail.com/MrAmghezi"
```

## challenge
View the full list of challenges or inspect details of a single challenge.

### list (ls)
Shows the list of challenges. Each challenge has a unique identifier which must be used if you want to get any data related to that challenge (e.g., view the scoreboard for a challenge or its description)

```
$ xero challenge list
     
     Challenges
  • Ugupugu     [rsa]

```
In the above example, the challenge `Ugupugu` has the identifier `rsa`. 

### show
Shows details and descriptions of a challenge.

In the following example, we'll inspect the details of `Ugupugu` challenge:
```
$ xero challenge show rsa

     Ugupugu

# Description

Can u decrypt this?
je9VobhwQWIGoNE3ugUBtJWPAYPJnQbaYJiA1BQqc2/6JlYjnN6nyD9gy78n06pjSg0anf7y3+02JbNI9kdksP+ZD+fNfFrSbii1...
```

Note that we used the challenge's identifier (`rsa`) and not its title (`Ugupugu`) in the command line.

## submit
Use this command to submit a flag you've captured!
You need to pass **challenge identifier** and **your flag** as arguments.

In the following example, we'll submit the flag `asc309vlk3m2lvpo` for the `Ugupugu` challenge:
```
$ xero submit rsa asc309vlk3m2lvpo
```

## scoreboard
View the scoreboard for a particular challenge.
```
$ xero scoreboard rsa

     Scoreboard

Total Records: 1

Pos | Team Name  | Captured | Total Submissions | Last Submission
1   | Hallelujah | 1        | 15                | 2022-05-31 13:04:55
```

## rules
Shows specified rules by the organizers.

```
$ xero rules
```

## completion
Generates the autocompletion script for the specified shell.
Available shells are:
- bash
- fish
- powershell
- zsh

For example, to use autocompletion in bash run the following commands:
```
xero completion bash > $HOME/xerocompletion
source $HOME/xerocompletion
```

Obviously you can setup your shell to source the autocompletion file on every reboot so you don't need to source it yourself. For example if you use bash, add a `source` line in your `$HOME/.bashrc` file:

**`.bashrc`**
```
...

source $HOME/<your autocompletion filename>

...
```

# License
Distributed under the MIT License. See [`LICENSE`](https://github.com/SoroushTaheri/xero-cli/blob/master/LICENSE) for more information.

# Contact
Soroush Taheri
- Email: soroushtgh@gmail.com
- Telegram: @soroushtaheri
