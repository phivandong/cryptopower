// util contains functions that don't contain layout code. They could be considered helpers that aren't particularly
// bounded to a page.

package page

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/planetdecred/godcr/ui/load"

	"gioui.org/gesture"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

func translateErr(err error) string {
	switch err.Error() {
	case dcrlibwallet.ErrInvalidPassphrase:
		return values.String(values.StrInvalidPassphrase)
	}

	return err.Error()
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func EditorsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			return false
		}
	}
	return true
}

func GenerateRandomNumber() int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

// getLockWallet returns a list of locked wallets
func getLockedWallets(wallets []*dcrlibwallet.Wallet) []*dcrlibwallet.Wallet {
	var walletsLocked []*dcrlibwallet.Wallet
	for _, wl := range wallets {
		if !wl.HasDiscoveredAccounts && wl.IsLocked() {
			walletsLocked = append(walletsLocked, wl)
		}
	}

	return walletsLocked
}

// createClickGestures returns a slice of click gestures
func createClickGestures(count int) []*gesture.Click {
	var gestures = make([]*gesture.Click, count)
	for i := 0; i < count; i++ {
		gestures[i] = &gesture.Click{}
	}
	return gestures
}

// showBadge loops through a slice of recent transactions and checks if there are transaction from different wallets.
// It returns true if transactions from different wallets exists and false if they don't
func showLabel(recentTransactions []wallet.Transaction) bool {
	var name string
	for _, t := range recentTransactions {
		if name != "" && name != t.WalletName {
			return true
		}
		name = t.WalletName
	}
	return false
}

func goToURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Error(err)
	}
}

func computePasswordStrength(pb *decredmaterial.ProgressBarStyle, th *decredmaterial.Theme, editors ...*widget.Editor) {
	password := editors[0]
	strength := dcrlibwallet.ShannonEntropy(password.Text()) / 4.0
	pb.Progress = float32(strength * 100)
	pb.Color = th.Color.Success
}

func ticketStatusIcon(th *decredmaterial.Theme, ic load.Icons, ticketStatus string) *struct {
	icon       *widget.Image
	color      color.NRGBA
	background color.NRGBA
} {
	m := map[string]struct {
		icon       *widget.Image
		color      color.NRGBA
		background color.NRGBA
	}{
		"UNMINED": {
			ic.TicketUnminedIcon,
			th.Color.DeepBlue,
			th.Color.LightBlue,
		},
		"IMMATURE": {
			ic.TicketImmatureIcon,
			th.Color.DeepBlue,
			th.Color.LightBlue,
		},
		"LIVE": {
			ic.TicketLiveIcon,
			th.Color.Primary,
			th.Color.LightBlue,
		},
		"VOTED": {
			ic.TicketVotedIcon,
			th.Color.Success,
			th.Color.Success2,
		},
		"MISSED": {
			ic.TicketMissedIcon,
			th.Color.Gray,
			th.Color.LightGray,
		},
		"EXPIRED": {
			ic.TicketExpiredIcon,
			th.Color.Gray,
			th.Color.LightGray,
		},
		"REVOKED": {
			ic.TicketRevokedIcon,
			th.Color.Orange,
			th.Color.Orange2,
		},
	}
	st, ok := m[ticketStatus]
	if !ok {
		return nil
	}
	return &st
}

func HandleSubmitEvent(editors ...*widget.Editor) bool {
	var submit bool
	for _, editor := range editors {
		for _, e := range editor.Events() {
			if _, ok := e.(widget.SubmitEvent); ok {
				submit = true
			}
		}
	}
	return submit
}

func GetAbsolutePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %s", err.Error())
	}

	exSym, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", fmt.Errorf("error getting filepath after evaluating sym links")
	}

	return path.Dir(exSym), nil
}

func handleSubmitEvent(editors ...*widget.Editor) bool {
	var submit bool
	for _, editor := range editors {
		for _, e := range editor.Events() {
			if _, ok := e.(widget.SubmitEvent); ok {
				submit = true
			}
		}
	}
	return submit
}
