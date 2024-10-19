# Photo & Video Organizer CLI Tool

## Overview

This CLI tool is designed to help organize photo and video files by renaming them based on EXIF metadata or filename patterns and adjusting their modification dates to match the parsed or EXIF dates. It processes `.jpg`, `.jpeg`, and `.mp4` files, making it easier to manage media files in a directory by ensuring consistency between file names and metadata.

## Features

- **EXIF Data Processing**: Extracts date and time from EXIF metadata of `.jpg` and `.jpeg` files to rename them accurately.
- **Filename Parsing**: Renames files based on dates parsed from their filenames if EXIF data is unavailable.
- **Modification Date Sync**: For `.mp4` video files, checks and updates the modified date to match the date extracted from the filename.
- **Error Handling**: Automatically skips files that are already named correctly or do not contain valid date information in their filenames.
- **Unique Naming**: Ensures file names remain unique when renaming files in the same folder by appending suffixes to avoid conflicts.
- **Cross-Platform Compatibility**: Works on both Windows and Unix-based systems.

## How It Works

The tool operates in two modes depending on the type of file:

1. **Image Files (`.jpg`/`.jpeg`)**:
   - Extracts the "Date Taken" from EXIF metadata, renames the file based on that date, and skips files that are already named correctly.
   - If EXIF data is unavailable, it parses the date from the filename to rename the file.

2. **Video Files (`.mp4`)**:
   - Extracts the date from the filename to rename the file.
   - Updates the file's modified date to match the parsed date if they differ.
   - If the file is already named correctly but has a mismatched modified date, the tool updates the modified date without renaming the file.

## Filename Format

Files are renamed using the following format:

YYYY-MM-DD HH.MM.SS.ext


- `YYYY` = Year
- `MM` = Month
- `DD` = Day
- `HH.MM.SS` = Time (Hour, Minute, Second)
- `ext` = File extension (`.jpg`, `.jpeg`, `.mp4`)

## Installation

To install the tool, clone the repository and build the project using Go.

```bash
git clone https://github.com/bibauporto/photosOrganizer.git
cd photosOrganizer
go build
```

# Usage

After building the tool, run it from the command line to organize files in a specified directory:
./photosOrganizer /path/to/your/folder


Replace /path/to/your/folder with the actual path to the folder containing your photos and videos.
Command Options

    Images: The tool processes .jpg and .jpeg files by reading EXIF metadata or parsing dates from the filenames.
    Videos: The tool processes .mp4 files by parsing dates from the filenames and checking/modifying their modified dates.

# Examples
Renaming a .jpg file with EXIF metadata:

    Original filename: IMG_1234.jpg
    EXIF "Date Taken": 2023-05-12 14:32:05
    Result: 2023-05-12 14.32.05.jpg

Renaming a .mp4 file based on the filename:

    Original filename: video_20231014.mp4
    Parsed date from filename: 2023-10-14
    Result: 2023-10-14 14.00.00.mp4

Updating the modified date of a .mp4 file:

    Original filename: 2023-10-14 14.00.00.mp4
    Current modified date: 2023-10-13 16:00:00
    Updated modified date: 2023-10-14 14:00:00

# Contributing

If you'd like to contribute to the project, feel free to open a pull request or submit issues for bugs and feature requests.
License

This project is licensed under the MIT License. See the LICENSE file for more information.