<a href="https://github.com/Mega-Kranus/NastySecrets/releases/latest">![GitHub release (with filter)](https://img.shields.io/github/v/release/mega-kranus/NastySecrets?style=flat-square&color=%23018c08)</a>
<a href="https://github.com/Mega-Kranus/NastySecrets/blob/main/go.mod">![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/mega-kranus/NastySecrets?style=flat-square)</a>
<a href="https://github.com/Mega-Kranus/NastySecrets/">![GitHub repo size](https://img.shields.io/github/repo-size/mega-kranus/NastySecrets?style=flat-square)</a>
<a href="https://github.com/Mega-Kranus/NastySecrets/blob/main/LICENSE">![GitHub License](https://img.shields.io/github/license/mega-kranus/NastySecrets?style=flat-square)</a>

# NastySecrets

NastySecrets, developed in GO, is a security/privacy tool. It encrypts files and can optionally rename them for added privacy. The tool automatically generates a configuration file with the encryption key and encrypted file names for easy future decryption.


## Features

- Simple
- Small
- Secure
- Swift




## Usage/Examples

Help menu
```bash
NastySecrets -h
```

Encryption without renaming files
```bash
NastySecrets -e -p /folder/to/encrypt -o /config/output/path
```

Encryption with renaming files
```bash
NastySecrets -e -n -p /folder/to/encrypt
```

Encryption and re-using an old key
```bash
NastySecrets -e -n -p folder/to/encrypt -k /config/file/path
```

Decryption
```bash
NastySecrets -d -k config/file/path -p /folder/to/decrypt
```



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
bin/NastySecrets -h
```

If you want to build and put the binary in "$HOME/.local/bin" (_run from anywhere_)
```bash
cd NastySecrets
```
```bash
make all
```
```bash
NastySecrets -h
```

## Authors

[@Mega-Kranus](https://www.github.com/Mega-Kranus)


## Documentation

[Overview](https://github.com/Mega-Kranus/NastySecrets/tree/main/docs)


## License

[GLP V 3.0](https://github.com/Mega-Kranus/NastySecrets/blob/main/LICENSE)


## Feedback

If you have any feedback, please reach out to me at Mega-Kranus@proton.me
