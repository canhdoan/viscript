package wallet

import (
	"fmt"
	"time"

	"encoding/hex"

	"github.com/skycoin/skycoin/src/cipher"
	//"math/rand"

	logging "github.com/op/go-logging"
)

var (
	logger = logging.MustGetLogger("wallet")
)

// WalletExt  wallet file extension
const WalletExt = "wlt"

// WalletTimestampFormat  wallet timestamp layout
const WalletTimestampFormat = "2006_01_02"

// NewWalletFilename check for collisions and retry if failure
func NewWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	//should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}
