# ZIP Password Cracker

A high-performance command-line tool for cracking password-protected ZIP archives using wordlist-based brute-force attacks. Built with Go for speed and efficiency.

## Overview

ZIP Password Cracker is designed to recover lost or forgotten passwords for encrypted ZIP files. The tool iterates through a list of candidate passwords and attempts to extract files from the archive, reporting success when a valid password is found.

## Features

- Fast password cracking with Go's concurrency capabilities
- Real-time progress bar showing cracking status
- Efficient file validation and error handling
- Support for large wordlists (automatically handles large file buffers)
- Clear and detailed console output
- Automatic detection of encrypted files within archives

## Requirements

- Go 1.11 or higher
- Standard library packages: bufio, fmt, os, strings

## Dependencies

The project uses two external Go packages:

1. **alexmullins/zip** - Enhanced ZIP file handling with password support
2. **schollz/progressbar/v3** - Visual progress bar for tracking cracking progress

## Installation

### Prerequisites

Ensure you have Go installed on your system. Download it from [golang.org](https://golang.org/dl/) if needed.

### Build from Source

1. Clone the repository:
   ```
   git clone https://github.com/gagaltotal/zipfile-cracker-tot.git
   cd zipfile-cracker-tot
   ```

2. Install dependencies:
   ```
   go mod init zipfile-cracker-tot
   go mod tidy
   OR
   go get github.com/alexmullins/zip
   go get github.com/schollz/progressbar/v3
   ```

3. Build the executable:
   ```
   go build -o zippcrack zip_cracker.go
   ```

## Usage

### Basic Syntax

```
zippcrack <zip_file> <wordlist>
```

### Parameters

- `<zip_file>` - Path to the password-protected ZIP archive
- `<wordlist>` - Path to the wordlist file containing candidate passwords (one password per line)

### Examples

Crack a ZIP file using the included wordlist:
```
zippcrack secret.zip wordlist.txt
```

Using a custom wordlist:
```
zippcrack archive.zip passwords.txt
```

### Output

The tool provides clear feedback during execution:

- Target ZIP file and wordlist paths being used
- File being tested for password validation
- Number of passwords in the wordlist
- Real-time progress bar with iteration count and speed
- Success message with found password or failure notification

## Wordlists

The project includes a `wordlist.txt` file with common passwords for testing purposes.

For more comprehensive cracking attempts, consider using larger wordlists:

- RockYou wordlist: [rockyou.txt](https://github.com/brannondorsey/naive-hashcat/releases/download/data/rockyou.txt)
- SecLists: [https://github.com/danielmiessler/SecLists](https://github.com/danielmiessler/SecLists)

## Project Structure

```
zipfile-cracker-tot/
├── images               Sample screenshot
├── zip_cracker.go       Main application source code
├── wordlist.txt         Sample wordlist for testing
├── zippcrack            Compiled executable (after building)
├── go.mod               Go module definition
└── README.md            This file
```

## How It Works

1. Validates command-line arguments
2. Checks that both the ZIP file and wordlist exist and are accessible
3. Opens the ZIP archive and locates an encrypted file
4. Counts the total number of passwords in the wordlist
5. Iterates through each password, attempting to extract files
6. Reports success when a valid password is found or failure if none match

## Technical Details

- Implements efficient file reading with buffered I/O
- Uses ZIP format inspection to identify encrypted files automatically
- Provides real-time visual feedback with customizable progress bar
- Handles errors gracefully with informative error messages

## Limitations

- Wordlist-based approach: only effective with passwords in the provided wordlist
- Single-threaded password testing: passwords are tested sequentially
- Requires valid ZIP file format with at least one encrypted file

## Author

GhostGTR666 (Gagaltotal666)

GitHub: [github.com/gagaltotal/zipfile-cracker-tot](https://github.com/gagaltotal/zipfile-cracker-tot)

## Legal Notice

This tool is provided for educational and authorized security testing purposes only. Users are responsible for ensuring they have permission to crack passwords on files they do not own. Unauthorized access to computer systems and files is illegal.
