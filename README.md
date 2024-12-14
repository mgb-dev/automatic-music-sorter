# Automatic Music Sorter (ams)

Simple CLI tool to sort audio files based on metadata

## Description

Point a directory (sub directories will be ignored) and create subdirectories based on the existing metadata from the analyzed files.

## Table of Contents <a id="toc"></a>

- [Install](#install)
- [Usage](#usage)
- [Credits](#credits)
- [License](#license)

## Install <a id="install"></a>

<small>[:arrow_up: Back to toc](#toc)</small>

```bash
git clone <this-url>
cd autmatic-music-sorter
make build
# Be sure to add ams's directory to your $PATH env variable
```

## Usage <a id="usage"></a>

<small>[:arrow_up: Back to toc](#toc)</small>

> [!WARNING]
> at this time ams has no "undo" (nor I plan to add one) .i.e ams **will make** modifications to your file system, which cannot be reverted by the program itself.

```bash
# Navigate into the directory to be sorted and specify a criteria to be sorted by
# Enjoy the automation
ams . artist
```

## Credits <a id="credits"></a>

<small>[:arrow_up: Back to toc](#toc)</small>

This project would not exist without:

<a id="tag-link" href="https://github.com/dhowden/tag" target="_blank">dhowden/tag@Github</a> : library to read metadata from audio files

## License <a id="license"></a>

<small style="justify-self: 'right'">[:arrow_up: Back to toc](#toc)</small>

MIT
