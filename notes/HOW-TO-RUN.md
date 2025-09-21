## How To Run

Steps to run a simple crazyhttp server. The following directives are genric
but have been tested only on `arch`. You should ideally use an equivalent command
in your distro and it should work.

### Setting up the server
1. Install prerequisite tools
    + [go 1.24.4](https://go.dev) (should work with any go version > 1.20, but this project has go v1.24)
        `sudo pacman -Syu go`
    + mkcert, to generate a self signed SSL certificate.
        `sudo pacman -S mkcert`
        `sudo pacman -S openssh`, might also require openssh
    + might need these libs too for debugging issues, or diving further.
        `sudo pacman -S tcpdump lsof`
2. Go over to `examples/` directory, and use any example to run. Say we do a naive server implementation. `cd naive/`
3. Populate your config.yaml, for TLS certificate, populate either path or raw. In raw, you can simply put
   file contents of your script as-is after escaping new line with `\n`. Remember spaces are consider as 
   different characters in certificate and should be removed.
4. Run `go mod run main.go`. If you see any error in port allocation make sure to grant privelege for lower numbered ports.
   We are using port 443 for all our examples here.

A server directory should simply look like 
```
go-project/
├── your-files-and-folders
| 
|   ## only 3 files are needed
|   ...
├── config.yaml // must be in root directory
├── cert.pem // not needed if you populate certificate.raw in config.yaml; path is configurable, can be anywhere
├── key.pem // not needed if you populate key.raw in config.yaml; path is configurable, can be anywhere
|   ...
| 
|
├── go.mod
└── go.sum

```

### Setting up the client
Trust the Root CA of your certificate, you should find it in `/etc/ssl/certs/ca-certificates.crt` or 
all the below examples will have a way to pass certificates explicitely. Omit them if you have trusted already.

On production you would skip the certificate trust process, and use your existing Domain's certificate.

1. Try using curl
```sh
curl https://172.27.192.203:443/test --cacert /etc/ssl/certs/ca-certificates.crt --show-error -v --http3-only
```

Or run a simple client via chrome
1. Launch Chrome using these flags [for arch, find windows below]
```zsh
google-chrome-stable \
--enable-quic \
--quic-version=80 \
--origin-to-force-quic-on=172.27.192.203:443 // replace this with your local server i\
--log-net-log="path/to/chrome-net-export-log-1.json" // replace this with your net-log file (helps in debugging) \
--net-log-capture-mode=IncludeSensitive \
--ignore-certificate-errors-spki-list=<base64_hash_of_certificate> // find steps to generate this below \
--user-data-dir="/Temp/ChrTemp" // some random directory for temp usage \
--enable-features=NetworkService,NetworkServiceInProcess \
--disable-features=ChromeWhatsNewUI  \
--disable-dev-shm-usage \
--disable-extensions \
--disable-software-rasterizer \
--test-type // disable that annoying top bar
```

or for windows
```sh
chrome.exe ^
--enable-quic ^
--quic-version=80 ^
--origin-to-force-quic-on=172.27.192.203:443 ^ // replace this with your local server ip
--log-net-log="path\to\chrome-net-export-log-1.json" ^ // replace this with your net-log file (helps in debugging)
--net-log-capture-mode=IncludeSensitive ^
--ignore-certificate-errors-spki-list=<base64_hash_of_certificate> ^ // find steps to generate this below
--user-data-dir="\Temp\ChrTemp" ^ // some random directory for temp usage
--enable-features=NetworkService,NetworkServiceInProcess ^
--disable-features=ChromeWhatsNewUI  ^
--disable-dev-shm-usage ^
--disable-extensions ^
--disable-software-rasterizer ^
--test-type // disable that annoying top bar
```

2. How to generate the base64 of certificate hash. Here cert.pem is essentially localhost.pem if generated using mkcert.
```zsh
openssl x509 -in cert.pem -noout -pubkey | openssl pkey -pubin -outform DER | openssl dgst -sha256 -binary | openssl base64
```

### Other Notes
1. When you generate a public-private key (SSL Certificate) using mkcert, it generates two files
`localhost.pem`, and `localhost-key.pem`. We call `localhost.pem` -> `cert.pem`, and `localhost-key.pem` -> key.pem
in this lib everywhere.
