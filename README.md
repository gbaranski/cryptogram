# Cryptogram
Fully decentralized P2P messaging app built on top of [LibP2P](https://libp2p.io/). Currently only CLI version is available.

![CLI Showcase](docs/cli-showcase-1.gif)

## Quick start with CLI app
##### Install Go
###### macOS
```bash
brew install go
```
###### Arch Linux
```bash
sudo pacman -S go
```

#### Download repository

```bash
git clone https://github.com/gbaranski/cryptogram # Clone this repo
cd cryptogram/cli
```

#### Run directly
```bash
# Run with MDNS peer discovery and nickname "Charlie"
go run main.go -mdns --nick Charlie 

# Run with DHT peer discovery with default bootstrap settings and nickname "Charlie"
go run main.go -dht --nick Charlie 
```

#### Compile to single executable
```bash
# Compile packages and dependencies into single executable
go build
# Run with DHT peer discovery with default bootstrap settings and nickname "Charlie"
./cli -dht --nick Charlie
```

