<p align="center">
  <h2 align="center">LevelDB Dumper</h3>
	
  <p align="center">
  <a href="https://travis-ci.com/github/mdawsonuk/LevelDBDumper" alt="Travis CI">
		<img src="https://img.shields.io/travis/com/mdawsonuk/LevelDBDumper?style=flat-square" /></a>
  <a href="LICENSE" alt="Licence">
		<img src="https://img.shields.io/github/license/mdawsonuk/LevelDBDumper?style=flat-square" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/releases" alt="Releases">
		<img src="https://img.shields.io/github/v/release/mdawsonuk/LevelDBDumper?include_prereleases&style=flat-square&color=blue" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/issues" alt="Issues">
		<img src="https://img.shields.io/github/issues/mdawsonuk/LevelDBDumper?style=flat-square" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/releases" alt="Downloads">
		<img src="https://img.shields.io/github/downloads/mdawsonuk/LevelDBDumper/total?style=flat-square" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/pulse" alt="Maintenance">
		<img src="https://img.shields.io/maintenance/yes/2021?style=flat-square" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/tree/master/src/LevelDBDumper">
		<img src="https://img.shields.io/github/languages/code-size/mdawsonuk/LevelDBDumper?style=flat-square"
			alt="Repo Size"></a>
  </p>
  <p align="center">
    Enumerates all Key values in a LevelDB database and outputs their corresponding Value
    <br />
    <a href="https://github.com/mdawsonuk/LevelDBDumper/issues/new?labels=bug">Report a Bug</a>
    ·
    <a href="https://github.com/mdawsonuk/LevelDBDumper/issues/new?labels=enhancement">Request Feature</a>
  </p>
</p>

### _Like the project? Please give it a star :star: to show your appreciation and watch the project to get notified on new releases_

## Table of Contents

* [About the Project](#about-the-project)
* [Supported Applications](#supported-applications)
* [Usage](#usage)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [To-Do](#to-do)
* [Contributing](#contributing)
* [License](#license)

## About The Project
This project was created out of a lack of a cross-platform tool which was able to enumerate every key in a LevelDB database and output its value. Other Level DB dumpers either were limited to one OS or had a complex installation process, so there was a need for a dumper which can be downloaded and run without dependency or installation issues.

I'm by no means an expert at using Go, so the code might not be as efficient or clean as it could be. If you want to help improve code quality, please consider [contributing](#contributing).

A massive thanks to Harsh Vardhan Singh and his [repo](https://github.com/harshvsingh8/leveldb-reader) which laid the foundations for some of the LevelDB enumeration work.

## Supported Applications

LevelDB Dumper has been tested to work on the following applications and provide the following forensics artifacts:

* [Discord](https://discord.com/) - User email, User token, Recent games, Search history, Draft messages, Collapsed categories/channels, User settings, GIF favourites (Local Storage DB)
* [WhatsApp](https://www.whatsapp.com/) - User name, mobile phone type (Local Storage DB)

_This works on Local Storage LevelDBs - support for IndexedDBs using the idb_cmp1 comparator is coming soon, which will enable support for applications such as Skype, Microsoft Teams, and more_

## Usage

```
Command Line: -h

      h/help              Display this help message.
      d/dir               Directory to recursively process. This is required.
      q/quiet             Don't output all key/value pairs to console. Default will output all key/value pairs
      t/outputType        Output type. Can be "csv" or "json"
      o/outputDir         Directory to save all output results to. Required for any file output
      f/outputFile        Filename to use when saving output. This will be appended with path and date
      b/batch             Combine all output files into one file. Supported by "csv" and "json" file types
      c/clean-output      Clean the file output of non-visual characters, such as \\u001
      no-colour/no-color  Don't colourise output
      no-header           Don't display the header
      u/check-update      Check for updates only

Short options (single letter) are prefixed with a single dash. Long commands are prefixed with two dashes

Examples: LevelDBParser.exe -d "C:\Temp\leveldb"
          LevelDBParser.exe -d "C:\Temp\leveldb" -o "C:\Temp" -q
          LevelDBParser.exe -d "C:\Temp\leveldb" --no-colour --quiet --no-header --clean-output
          LevelDBParser.exe -d "C:\Temp\leveldb" --no-colour -b --outputType json -outputFile Evidence.json
          LevelDBParser.exe -d "C:\Temp\leveldb" -t csv -f LevelDB.csv -o Evidence -b --no-colour --quiet
          LevelDBParser.exe --check-update
          LevelDBParser.exe --help
```

LevelDB Dumper will search recursively from the directory passed to it for LevelDB databases. Upon finding one, it will be queued for dumping. Once it has searched the entire drive, the databases will be enumerated from the item list. By default, the Key/Value pairs are output to the console.

It is recommended to specify an output file for dumping. Using `-o <Directory>` will output a file per LevelDB database found, with the file name containing the timestamp of dumping and path to the LevelDB database. The default format is CSV. If you wish all the output to be in one file, `-b` will batch the outputs and add additional columns for path and the timestamp it was dumped at.

It is worth noting that all Unicode control characters/non-graphics characters are stripped from the output strings but are retained for file output, such as CSV. For applications such as Discord, where null terminators are found in Key names, this is used to improve output formatting. 

There have been issues with Windows 10 where the program is opened in a new window instead of the current Command Line window instance, meaning that the output is not visible. A work-around for this appears to be running the Command Prompt/Powershell as Administrator. However, for analysis of output, the key/value pairs should be output to a file rather than redirecting or analysing through the command line window.

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

LevelDB is written in Go, so it can be run cross-platform. If you are building from source, you will need to [install Golang](https://golang.org/doc/install)

### Installation

#### From Releases

1. Download the latest [release](https://github.com/mdawsonuk/LevelDBDumper/releases) for your platform of choice.

2. That's it!

_Commands available in the README can change from time to time. To get the available commands for the latest release, run with --help_

#### From Repo

1. Clone the repo
```sh
git clone https://github.com/mdawsonuk/LevelDBDumper.git
```

2. Using Go CLI, get the required packages to build
```sh
go get github.com/syndtr/goleveldb/leveldb
go get github.com/hashicorp/go-version
```

If building on Windows, get the Windows sys package
```sh
go get golang.org/x/sys/windows
```

3. Using Go CLI, build the application
```sh
cd src/LevelDBDumper
go build
```

4. That's it! An executable should be created in that directory. View the article [here](https://medium.com/@utranand/building-golang-package-for-linux-from-windows-22fa23764808) for information on cross-platform compilation.

## To-Do
* ~~Recursively parse from directory instead of providing one LevelDB database~~ :heavy_check_mark:
* ~~Improved help dialog with list of arguments and examples~~ :heavy_check_mark:
* ~~Export to CSV file for each discovered LevelDB database~~ :heavy_check_mark:
* ~~Quiet mode to avoid dumping all Key/Values~~ :heavy_check_mark:
* ~~Truncate long Values in non-quiet output~~ :heavy_check_mark:
* ~~Display coloured Key/Value for non-quiet mode~~ :heavy_check_mark:
* ~~Ignore processing empty LevelDB databases to avoid creating empty output files~~ :heavy_check_mark:
* ~~Travis CI builds for Windows and Linux~~ :heavy_check_mark:
* ~~Allow toggling of output colouring~~ :heavy_check_mark:
* ~~Check if user has Administrator/root privileges~~ :heavy_check_mark:
* ~~Detect comparator type from MANIFEST file~~ :heavy_check_mark:
* ~~Update check against GitHub (check latest release)~~ :heavy_check_mark:
* ~~JSON export~~ :heavy_check_mark:
* ~~`--no-header` command to hide header~~ :heavy_check_mark:
* ~~`-u/--update` command to only check for update~~ :heavy_check_mark:
* ~~Allow optional cleaning of output of non-graphical characters~~ :heavy_check_mark:
* Batch output file (All LevelDB dumps into one file)
* Allow customisation of output name
* Implement idb_cmp1 Google comparator for IndexedDB databases
* Allow suppression of warning/error messages through command line e.g. `2>/dev/null`

## Contributing

Want to make the tool better? Improve the code? Pull requests are accepted and very much appreciated.

## License

Distributed under the GPLv3 License. See [LICENSE](LICENSE) for more information.