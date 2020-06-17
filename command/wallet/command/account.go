package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/keystore"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	"github.com/Futuremine-chain/futuremine/tools/crypto/mnemonic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	accountCmds := []*cobra.Command{
		CreateAccountCmd,
	}

	RootCmd.AddCommand(accountCmds...)
	RootSubCmdGroups["account"] = accountCmds
}

var CreateAccountCmd = &cobra.Command{
	Use:     "CreateAccount {password}",
	Short:   "CreateAccount {password}; Create account;",
	Aliases: []string{"createaccount", "CA", "ca"},
	Example: `
	CreateAccount  
		OR
	CreateAccount 123456
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  CreateAccount,
}

func CreateAccount(cmd *cobra.Command, args []string) {
	var passWd []byte
	var err error
	if len(args) == 1 && args[0] != "" {
		passWd = []byte(args[0])
	} else {
		fmt.Println("please set account password, cannot exceed 32 bytes：")
		passWd, err = readPassWd()
		if err != nil {
			log.Error(cmd.Use+" err: ", fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	if len(passWd) > 32 {
		log.Error(cmd.Use+" err: ", fmt.Errorf("password too long! "))
		return
	}
	entropy, err := kit.Entropy()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	mnemonicStr, err := kit.Mnemonic(entropy)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	key, err := kit.MnemonicToEc(mnemonicStr)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("generate secp256k1 key failed! %s", err.Error()))
		return
	}
	p2pId, err := kit.GenerateP2PID(key)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("generate p2p id failed! %s", err.Error()))
	}
	if j, err := keystore.GenerateKeyJson(Net, Cfg.KeystoreDir, key, mnemonicStr, passWd); err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("generate key failed! %s", err.Error()))
	} else {
		j.P2pId = p2pId.String()
		bytes, _ := json.Marshal(j)
		output(string(bytes))
	}
}

func readPassWd() ([]byte, error) {
	var passWd [33]byte

	n, err := os.Stdin.Read(passWd[:])
	if err != nil {
		return nil, err
	}
	if n <= 1 {
		return nil, errors.New("not read")
	}
	return passWd[:n-1], nil
}

var ShowAccountCmd = &cobra.Command{
	Use:     "ShowAccounts",
	Short:   "ShowAccounts; Show all account of the wallet;",
	Aliases: []string{"showaccounts", "sa", "SA"},
	Example: `
	ShowAccounts
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  ShowAccount,
}

func ShowAccount(cmd *cobra.Command, args []string) {
	if addrList, err := keystore.ReadAllAccount(Cfg.KeystoreDir); err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("read account failed! %s", err.Error()))
	} else {
		bytes, _ := json.Marshal(addrList)
		output(string(bytes))
	}
}

var DecryptAccountCmd = &cobra.Command{
	Use:     "DecryptAccount {address} {password} {keyfile}；Decrypting account json file generates the private key and mnemonic;；",
	Short:   "DecryptAccount {address} {password} {keyfile}; Decrypting account json file generates the private key and mnemonic;",
	Aliases: []string{"decryptaccount", "DA", "da"},

	Example: `
	DecryptAccount 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf
		OR
	DecryptAccount 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf 123456
		OR
	DecryptAccount 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf 123456 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf.json
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  DecryptAccount,
}

func DecryptAccount(cmd *cobra.Command, args []string) {
	var passWd []byte
	var keyFile string
	var err error
	if len(args) >= 2 && args[1] != "" {
		passWd = []byte(args[1])
	} else {
		fmt.Println("please input password：")
		passWd, err = readPassWd()
		if err != nil {
			log.Error(cmd.Use+" err: ", fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	if len(args) == 3 && args[2] != "" {
		keyFile = args[2]
	} else {
		keyFile = getAddJsonPath(args[0])
	}

	privKey, err := ReadAddrPrivate(keyFile, passWd)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("wrong password"))
		return
	}

	bytes, _ := json.Marshal(privKey)
	output(string(bytes))
}

var MnemonicToAccountCmd = &cobra.Command{
	Use:     "MnemonicToAccount {mnemonic} {password}；Restore address by mnemonic and set new password;",
	Short:   "MnemonicToAccount {mnemonic} {password}; Restore address by mnemonic and set new password;",
	Aliases: []string{"mnemonictoaccount", "MTA", "mta"},
	Example: `
	MnemonicToAccount "sadness ladder sister camp suspect sting height diagram confirm program twist ostrich blush bronze pass gasp resist random nothing recycle husband install business turtle"
		OR
	MnemonicToAccount "sadness ladder sister camp suspect sting height diagram confirm program twist ostrich blush bronze pass gasp resist random nothing recycle husband install business turtle" 123456
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  MnemonicToAccount,
}

func MnemonicToAccount(cmd *cobra.Command, args []string) {
	var passWd []byte
	var err error
	priv, err := mnemonic.MnemonicToEc(args[0])
	if err != nil {
		log.Error(cmd.Use+" err: ", errors.New("[mnemonic] wrong"))
		return
	}
	if len(args) == 2 && args[1] != "" {
		passWd = []byte(args[1])
	} else {
		fmt.Println("please set address password, cannot exceed 32 bytes：")
		passWd, err = readPassWd()
		if err != nil {
			log.Error(cmd.Use+" err: ", fmt.Errorf("read pass word failed! %s", err.Error()))
			return
		}
	}
	if len(passWd) > 32 {
		log.Error(cmd.Use+" err: ", fmt.Errorf("password too long! "))
		return
	}
	p2pId, err := kit.GenerateP2PID(priv)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("generate p2p id failed! %s", err.Error()))
	}
	if j, err := keystore.GenerateKeyJson(Net, Cfg.KeystoreDir, priv, args[0], passWd); err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("generate key failed! %s", err.Error()))
	} else {
		j.P2pId = p2pId.String()
		bytes, _ := json.Marshal(j)
		output(string(bytes))
	}
}

func getAddJsonPath(addr string) string {
	return Cfg.KeystoreDir + "/" + addr + ".json"
}

func ReadAddrPrivate(jsonFile string, password []byte) (*keystore.Private, error) {
	j, err := keystore.ReadJson(jsonFile)
	if err != nil {
		return nil, err
	}
	return keystore.DecryptPrivate(password, j)
}
