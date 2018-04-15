# Mobile API Description

The mobile APIs will manager wallet and communicate with skycoin-class daemon. 

## Build

Use the following cmd to build an ios framework file; for android just replace the target `ios` with `android`.

```bash
$ gomobile bind -target=ios github.com/spolabs/wallet-api/src/api/mobile
```

## Build android lib in osx

Install android sdk and ndk

```bash
brew tap caskroom/cask
brew cask install android-sdk
brew cask install android-ndk
export ANDROID_HOME=/usr/local/cask/android-sdk/$version
export ANDROID_NDK=/usr/local/cask/android-ndk/$version
export PATH=$ANDROID_HOME/tools:$PATH
export PATH=$ANDROID_HOME/platform-tools:$PATH

ln -s $ANDROID_NDK $ANDROID_HOME/ndk-bundle
```

Download platforms which contains android.jar file into $ANDROID_NDK/platforms/

```bash
cd $ANDROID_NDK/platforms
git clone https://github.com/Sable/android-platforms
// then copy the android.jar files from the android-platforms/android-* to corresponding folders
```

Init the gomobile

```bash
gomobile init
```

## APIs

### Initialization

```go
func Init(walletDir, passwd string)
```

We use `walletDir` to init the API env, Wallet dir is the place for persisting the wallet files;

also, load wallet already exists 

Params:

* walletDir: walelt directory 
* passwd: password for load wallet 

Return:

* first: error info

### Register a new coin into wallet 

must register coin after Init

```go
func RegisterNewCoin(coinType, serverAddr string) error 
```

Params:

* coinType: can be `skycoin` `spo` `suncoin` and so on
* serverAddr: the server address is consisted of ip and port, eg: 127.0.0.1:6420

Return:

* first: error info

### Create wallet

Create wallet base on coin type, lable and seed. if seed is "", then will auto generate a random seed

```go
func NewWallet(coinType, lable, seed, passwd string) (string, error)
```

Params:

* coinType: can be `skycoin` `spo` `suncoin` and so on
* lable: identified wallet 
* seed: wallet seed, can be any string
* passwd: password for wallet

Return:

* frist: wallet id
* second: error info

Example:
   coinType: skycoin
   lable: lableandseed
   return:
   skycoin_lableandseed

### Create address

This api is used to create addresses in specific wallet.

```go
func NewAddress(walletID string, num int, passwd string) (string, error)
```

Params:

* walletID: wallet id, which was generated by `NewWallet` api
* num: the address number you need to generate
* passwd : password for wallet

Return:

* first: address entries json, eg:

```json
{
    "addresses": [
        {
            "address": "QNpH7Y2spJtSAbdufM4qwchWvg71mAsbNx",
            "pubkey": "02e1eaa54233495faed0d50ecbfdc3e2e9fcac829b3d406a4d7bde43ff4452a0f7",
            "seckey": "58824dc46b14e28ecd0a4e93835a61129c770dda58073e8a7bd042d6b5f32a17"
        },
        {
            "address": "2Wm8wyZPh6HtFUBMAEewA2ZHxXbAvX4n5En",
            "pubkey": "0263c39ab3dba2c0bc8c08d4fdab297a21ff4505679cb7b8d832af27e4db7a0344",
            "seckey": "81c0fb22b4a23570f0bc7d30bade4bfca47f6d7f9e0a59613a54eccc333197ed"
        }
    ]
}
```

* second: error info

### Get addresses in wallet

This api is used to get all generated addresses in specific wallet

```go
func GetAddresses(walletID string) (string, error)
```

Params:

* walletID: wallet id

Return:

* first: address list json, eg:

```json
{
    "addresses": [
        "QNpH7Y2spJtSAbdufM4qwchWvg71mAsbNx",
        "2Wm8wyZPh6HtFUBMAEewA2ZHxXbAvX4n5En"
    ]
}
```

### Get pubkey and seckey pair of address

This api is used to get keypair of specific address.

```go
func GetKeyPairOfAddr(walletID string, addr string) (string, error)
```

Param:

* walletID: id of the wallet you are going to query
* addr: coin address

Return:

* first: keypair json, eg:

```json
{
    "pubkey": "02e1eaa54233495faed0d50ecbfdc3e2e9fcac829b3d406a4d7bde43ff4452a0f7",
    "seckey": "KzBm5cmRgGEgPESM3izfvesdz8caf14c6pZ3spG5eaJNE9abZmuc"
}
```

* second: error info

### Get balance

This api is used to query the balance of specific address.

```go
func GetBalance(coinType string, address string) (string, error)
```

Params:

* coinType: the coin type, can be `skycoin` or `spo`
* address: coin address

Return:

* frist: balance json, eg:

```json
{
    "balance":"40.000000",
    "hours":"32001"
}
```

the balance unit of skycoin is `drop`, spo is `drop`.

### Get wallet balance

This api is used to query the balance of specific wallet.

```go
func GetWalletBalance(coinType string, wltID string) (string, error)
```

Params:

* coinType: the coin type, can be `skycoin` or `spo`
* wltID: wallet id

Return:

* frist: balance json, eg:

```json
{
    "balance":"40.000000",
    "hours":"32001"
}
```

the balance unit of skycoin is `drop`, spo is `drop`.

### Send skycoin

This api can be used to send skycoin to one recipient address.

```go
func Send(coinType, walletID, toAddr, amount string) (string, error)
```

Params:
*coinType : coin type
* walletID: wallet id
* toAddr: recipient address
* amount: the coins you will send, it's value must be the multiple of 0.001.

Return:

* first: txid json, eg:

```json
{
    "txid":"05d52650917f4233795d12e76f7f228409863ce8b304b0d0dfc778f2b023112a"
}
```

* second: error info

### Get transaction

```go
func GetTransactionByID(coinType, txid string) (string, error)
```

Params:

* coinType: the coin type, can be `skycoin` or `spo`
* txid: transaction id

Return:

* first: detailed transaction info, eg:

```json
{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 89,
        "unknown": false
    },
    "txn": {
        "length": 183,
        "type": 0,
        "txid": "b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc0",
        "inner_hash": "7583587d02bedbeb3c15dde9e13baac36b0eb2b7ba7b2063c323a226d0784619",
        "sigs": [
            "67565680295b8758e07d0ee67f4f07b711e1c711da6af025dd4e2277de6e54941e35e5123f3d45eaa9bca131240eeb2067274199109eba17e5f8b1ee5aeef62301"
        ],
        "inputs": [
            "a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c5"
        ],
        "outputs": [
            {
                "uxid": "f9e39908677cae43832e1ead2514e01eaae48c9a3614a97970f381187ee6c4b1",
                "dst": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
                "coins": "1",
                "hours": 100
            }
        ]
    }
}
```

* second: error info

### Check transaction confirmed

```go
func IsTransactionConfirmed(coinType, txid string) (bool, error)
```

Params:

* coinType: the coin type, can be `skycoin` or `spo`
* txid: transaction id

Return:

* first: true if transaction confrimed, else false

* second: error info

### Get wallet seed 

```go
func GetSeed(walletID string) (string, error)
```

Params:

* walletID: wallet id

Return:

* first: wallet seed 

* second: error info

### Create wallet seed 

```go
func NewSeed() string
```

Params:

  None

Return:

* first: wallet seed 

### Validate address  

```go
func ValidateAddress(coinType, addr string) (bool, error) {
```

Params:

* coinType: the coin type, can be `skycoin` or `spo`
* addr : address

Return:

* first: true if address is valid, else false 

* second: error info

### Remove wallet 

```go
func Remove(walletID string) error 
```

Params:

* walletID: wallet id

Return:

* first: error info

### Wallet exists or not

```go
func IsExist(walletID string) bool 
```

Params:

* walletID: wallet id

Return:

* first: true if wallet exists, else false 

### Wallet contains addresses or not

```go
func IsContain(walletID string, addrs string) (bool, error) 
```

Params:

* walletID: wallet id

* addrs : address , notice many addresses should join by ",", such as a1,a2,a3

Return:

* first: true if contains, else false 

* second: error info
