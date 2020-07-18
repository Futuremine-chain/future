# future mine chain

### Start the testing

- Test with the command line wallet.
- Or you can test it by downloading the wallet APP, [DOWNLOAD](https://wallet.futuremine.io).
- Test COINS can be applied through APP wallet.

### How to build

####  Prerequisites

- Update Go to version at least 1.13 (required >= **1.13**)

Check your golang version

```bash
~ go version
go version go1.13 darwin/amd64
```

```bash
cd futuremine/cmd/futuremine
go build

cd futuremine/cmd/wallet
go build
```

#### How to use


##### Copy configuration file for reconfiguration

```bash
 cp config.toml.example config.toml
```

##### Modify configuration file

* set RpcPass
* set ExternalIp

##### Start the futuremine

```bash

./futuremine --config config.toml
```

##### Copy wallet configuration file for reconfiguration

```
cd cmd/wallet
cp wallet.toml.example wallet.toml
```

##### Modify wallet configuration file

* set RpcIp
* set RpcPass
* If the node has the RpcTLS switch turned on, you need to configure the node's server.pem path to RpcCert and set RpcTLS in wallet.config to true

##### Use wallet

```bash
./wallet --help
```
##### Create an account or set password at create

```bash
./wallet Create 

./wallet Create 123456
```
##### Send transaction

SendTransaction {from} {to} {contract} {amount} {fee} [password] [nonce]

```bash
./wallet SendTransaction xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 FMC 10 0.1

./wallet SendTransaction xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 FMC 10 0.1 123456
```

##### Create token

SendCreateToken {from} {to} {name} {shorthand} {amount} {fees} [password] [nonce]

```bash
./wallet SendCreateToken xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ "M token" MT 1000 0.1

./wallet SendCreateToken xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ "M token" MT 1000 0.1 123456
```

##### Get account balance
Account {address}
```bash
./wallet Account xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ
```

##### Get token record
Token {token address}
```bash
./wallet Token TfR8dgAAesNZum1PWDS9fFE6iRuTKWPrFCq
```

##### Become a super node to generate blocks

If you want to become a super node, you need to be a candidate and wait for the next round of elections after being voted.

* Configure the candidate account on the node config

    ```
    KeyFile = "cmd/wallet/keystore/xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ.json"
    KeyPass = "123456"
    ```
* Become a candidate

    SendCandidate {address} {fee} [password]

    ```bash
    ./wallet SendCandidate xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001
    ```
* Cancel candidate

    SendCancel {address} {fee} [password]

    ```bash
    ./wallet SendCancel xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001 123456
    ```
* Vote for candidates

    SendVote {address} {candidate} {fee} [password]

    ```bash
    ./wallet SendVote xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001
    ```
  
##### View the super nodes in an election cycle
    
```bash
./wallet CycleSupers 18450
```
    
##### About elections

* Up to 9 super nodes loop out blocks.
* Block every 5s.
* Re-election after every 24 hours .
* After more than two-thirds of the super nodes are confirmed, the block is fully confirmed.
* The number of votes is based on the total number of voters' current account balances, and the top 9 votes become the super node.
* If the super node produces less than one-third of the blocks that should have been produced in the current cycle, the candidate will be kicked out.