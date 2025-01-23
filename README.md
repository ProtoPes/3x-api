# 3x-api

API endpoint for managing configurations written in GO

## Plan

![Plan](Plan.odg) 

## Installation

### 

Compile from source code or download binary from releases

### Docker

- Pull built image: `sudo docker pull protopes/awg`
or
- Clone repo and build an image from Dockerfile
```
git clone https://github.com/ProtoPes/awg.git
sudo docker build -t $IMAGE_NAME .
sudo docker run -d $IMAGE_NAME
```

## Usage

Invoke program with flags:

`-c` or `--config`: generate config files from template
`-n` or `--name`: generate ranfom name
`-i` or `--ip`: find unused ip for client
`-g` or `--gen-ip`: generate unused ip adresses file
`-h` or `--help`: show usage


## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

## License

[GPLv3](https://www.gnu.org/licenses/gpl-3.0.html)
