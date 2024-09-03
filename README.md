<a href="https://github.com/Mega-Kranus/NastySecrets/releases/latest">![GitHub release (with filter)](https://img.shields.io/github/v/release/mega-kranus/NastySecrets?style=flat-square&color=%23018c08)</a>
<a href="https://github.com/Mega-Kranus/NastySecrets/blob/main/go.mod">![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/mega-kranus/NastySecrets?style=flat-square)</a>
<a href="https://github.com/Mega-Kranus/NastySecrets/">![GitHub repo size](https://img.shields.io/github/repo-size/mega-kranus/NastySecrets?style=flat-square)</a>
<a href="https://github.com/Mega-Kranus/NastySecrets/blob/main/LICENSE">![GitHub License](https://img.shields.io/github/license/mega-kranus/NastySecrets?style=flat-square)</a>

# NastySecrets

NastySecrets, developed in GO, is a security/privacy tool. It encrypts files and can optionally rename them for added privacy. The tool automatically generates a configuration file with the encryption key and encrypted file names for easy future decryption.

## Features

- Encrypt and decrypt files swiftly
- Rename files for added privacy and protection
- Encrypt/decrypt upto 25 files simultaneously
- Simple to use and very straight forward

## Installation

Clone the repository

```bash
git clone https://github.com/Mega-Kranus/NastySecrets.git
```

If you only want to build
```bash
cd NastySecrets
```
```bash
make build
```
```bash
bin/nastysecrets -h
```

If you want to build and put the binary in "$HOME/.local/bin" (_run from anywhere_)
```bash
cd NastySecrets
```
```bash
make all
```
```bash
nastysecrets -h
```
## Removal

In order to uninstall NastySecrets you can use the make file
```bash
cd NastySecrets
```

```bash
make uninstall
```

## Usage/Examples

Use man page after installation for detailed information
```bash
man nastysecrets
```

Help menu
```bash
nastysecrets -h
```

Encryption without renaming files
```bash
nastysecrets -e -p /folder/to/encrypt -o /config/output/path
```

Encryption with renaming files
```bash
nastysecrets -e -n -p /folder/to/encrypt
```

Encryption and re-using an old key
```bash
nastysecrets -e -n -p folder/to/encrypt -c /config/file/path
```

Decryption
```bash
nastysecrets -d -c config/file/path -p /folder/to/decrypt
```

## Authors

[@Mega-Kranus](https://www.github.com/Mega-Kranus)

## Documentation

[Overview](https://github.com/Mega-Kranus/NastySecrets/tree/main/docs)

## License

[GLP V 3.0](https://github.com/Mega-Kranus/NastySecrets/blob/main/LICENSE)

## Feedback

If you have any feedback, please reach out to me at m3gakr4nus@proton.me
