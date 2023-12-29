
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
./NastySecrets -h
```

Encryption without renaming files
```bash
./NastySecrets -e -p /folder/to/encrypt -o /config/output/path
```

Encryption with renaming files
```bash
./NastySecrets -e -n -p /folder/to/encrypt
```

Encryption and re-using an old key
```bash
./NastySecrets -e -n -p folder/to/encrypt -k /config/file/path
```

Decryption
```bash
./NastySecrets -d -k config/file/path -p /folder/to/decrypt
```



## Installation

Clone the repository

```bash
git clone https://github.com/Mega-Kranus/NastySecrets.git
```

Either use the binary provided
```bash
cd NastySecrets/bin
```
```bash
./NastySecrets -h
```

Or build it yourself (recommended)
```bash
cd NastySecrets
```
```bash
go build -o NastySecrets ./cmd/NastySecrets/main.go
```
```bash
./NastySecrets -h
```
## Documentation

[Overview](https://github.com/Mega-Kranus/NastySecrets/tree/main/docs)


## Authors

- [@Mega-Kranus](https://www.github.com/Mega-Kranus)


## Feedback

If you have any feedback, please reach out to me at Mega-Kranus@proton.me


## License

[GLP V 3.0](https://github.com/Mega-Kranus/NastySecrets/blob/main/LICENSE)