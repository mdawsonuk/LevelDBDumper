<p align="center">
  <h2 align="center">LevelDB Dumper</h3>
	
  <p align="center">
    <a href="LICENSE" alt="Licence">
		<img src="https://img.shields.io/github/license/mdawsonuk/LevelDBDumper?style=flat-square" /></a>
	<a alt="Releases">
		<img src="https://img.shields.io/github/v/release/mdawsonuk/LevelDBDumper?include_prereleases&style=flat-square&color=blue" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/issues" alt="Issues">
		<img src="https://img.shields.io/github/issues/mdawsonuk/LevelDBDumper?style=flat-square" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/releases" alt="Downloads">
		<img src="https://img.shields.io/github/downloads/mdawsonuk/LevelDBDumper/total?style=flat-square" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/pulse" alt="Maintenance">
		<img src="https://img.shields.io/maintenance/yes/2020?style=flat-square" /></a>
	<a href="https://github.com/mdawsonuk/LevelDBDumper/">
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

## Table of Contents

* [About the Project](#about-the-project)
* [TODO](#todo)
* [Usage](#usage)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
  * [Installation](#installation)
* [Contributing](#contributing)
* [License](#license)

## About The Project
This project was created out of a lack of a cross-platform tool which was able to enumerate every key in a LevelDB database and output its value. Other Level DB dumpers either were limited to one OS or had a complex installation process, so there was a need for a dumper which can be downloaded and run without dependency or installation issues.

I'm by no means an expert at using Go, so the code might not be as efficient or clean as it could be. If you want to help improve code quality, please consider [contributing](#contributing).

A massive thanks to Harsh Vardhan Singh and his [repo](https://github.com/harshvsingh8/leveldb-reader) which did much of the LevelDB enumeration work for me.

## TODO
* ~~Recursively parse from directory instead of providing one LevelDB database~~ :heavy_check_mark:
* ~~Improved help dialog with list of arguments and examples~~ :heavy_check_mark:
* ~~Export to CSV file for each discovered LevelDB database~~ :heavy_check_mark:
* ~~Quiet mode to avoid dumping all Key/Values~~ :heavy_check_mark:
* ~~Truncate long Values in non-quiet output~~ :heavy_check_mark:
* ~~Display coloured Key/Value for non-quiet mode~~ :heavy_check_mark:
* Travis CI builds for Windows, Linux, OS X
* Check if user has Administrator/root privileges
* Allow customisation of CSV output name
* Batch CSV file (All LevelDB dumps into one file)
* Ignore processing empty LevelDB databases to avoid creating empty output files
* JSON export
* Text export

## Usage

To dump all of the contents to a file, run the program on the Command Line, providing the full path to the LevelDB database.

```sh
LevelDBDumper.exe C:\Users\User\AppData\Roaming\Discord\Local Storage\leveldb
```

Redirecting the output to a file will remove line wrapping and make it easier to view and search for keys in a text document

```sh
LevelDBDumper.exe C:\Users\User\AppData\Roaming\Discord\Local Storage\leveldb > DiscordDBDump.txt
```

It is worth noting that all Unicode control characters/non-graphics characters are stripped from the output strings. For applications such as Discord, where null terminators are found in Key names, this is used to improve output formatting.

There have been issues with Windows 10 where the program is opened in a new window instead of the current Command Line window instance. A work-around for this appears to be running the Command Prompt/Powershell as Administrator.

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

LevelDB is written in Go, so it can be run cross-platform. If you are building from source, you will need to [install Golang](https://golang.org/doc/install)

### Installation

#### From Releases

1. Download the latest [release](https://github.com/mdawsonuk/LevelDBDumper/releases).

2. Unzip the compressed release to a folder of your choice.

3. That's it!

#### From Repo

1. Clone the repo
```sh
git clone https://github.com/mdawsonuk/LevelDBDumper.git
```

2. Using Go CLI, get the GoLevelDB package
```sh
go get github.com/syndtr/goleveldb/leveldb
```

3. Using Go CLI, build the application
```sh
cd src/LevelDBDumper
go build
```

4. That's it! An executable should be created in that directory. View the article [here](https://medium.com/@utranand/building-golang-package-for-linux-from-windows-22fa23764808) for information on cross-platform compilation.

## Contributing

Want to make the tool better? Improve the code? Pull requests are accepted and very much appreciated.

## License

Distributed under the GPLv3 License. See [LICENSE](LICENSE) for more information.