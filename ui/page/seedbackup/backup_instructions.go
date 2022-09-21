package seedbackup

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gitlab.com/raedah/cryptopower/app"
	"gitlab.com/raedah/cryptopower/libwallet/wallets/dcr"
	"gitlab.com/raedah/cryptopower/ui/cryptomaterial"
	"gitlab.com/raedah/cryptopower/ui/load"
	"gitlab.com/raedah/cryptopower/ui/modal"
	"gitlab.com/raedah/cryptopower/ui/page/components"
	"gitlab.com/raedah/cryptopower/ui/values"
)

const BackupInstructionsPageID = "backup_instructions"

type (
	C = layout.Context
	D = layout.Dimensions
)

type BackupInstructionsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	wallet *dcr.Wallet

	backButton  cryptomaterial.IconButton
	viewSeedBtn cryptomaterial.Button
	checkBoxes  []cryptomaterial.CheckBoxStyle
	infoList    *layout.List
}

func NewBackupInstructionsPage(l *load.Load, wallet *dcr.Wallet) *BackupInstructionsPage {
	bi := &BackupInstructionsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(BackupInstructionsPageID),
		wallet:           wallet,

		viewSeedBtn: l.Theme.Button("View seed phrase"),
	}

	bi.viewSeedBtn.Font.Weight = text.Medium

	bi.backButton, _ = components.SubpageHeaderButtons(l)
	bi.backButton.Icon = l.Theme.Icons.ContentClear

	bi.checkBoxes = []cryptomaterial.CheckBoxStyle{
		l.Theme.CheckBox(new(widget.Bool), "The 33-word seed word is EXTREMELY IMPORTANT."),
		l.Theme.CheckBox(new(widget.Bool), "seed word is the only way to restore your wallet."),
		l.Theme.CheckBox(new(widget.Bool), "It is recommended to store your seed word in a physical format (e.g. write down on a paper)."),
		l.Theme.CheckBox(new(widget.Bool), "It is highly discouraged to store your seed word in any digital format (e.g. screenshot)."),
		l.Theme.CheckBox(new(widget.Bool), "Anyone with your seed word can steal your funds. DO NOT show it to anyone."),
	}

	for i := range bi.checkBoxes {
		bi.checkBoxes[i].TextSize = values.TextSize16
	}

	bi.infoList = &layout.List{Axis: layout.Vertical}

	return bi
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) OnNavigatedTo() {

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) HandleUserInteractions() {
	for pg.viewSeedBtn.Clicked() {
		if pg.verifyCheckBoxes() {
			// TODO: Will repeat the paint cycle, just queue the next fragment to be displayed
			pg.ParentNavigator().Display(NewSaveSeedPage(pg.Load, pg.wallet))
		}
	}

}

func promptToExit(load *load.Load, pageNavigator app.PageNavigator, window app.WindowNavigator) {
	infoModal := modal.NewCustomModal(load).
		Title("Exit?").
		Body("Are you sure you want to exit the seed backup process?").
		SetNegativeButtonText(values.String(values.StrNo)).
		SetPositiveButtonText(values.String(values.StrYes)).
		SetPositiveButtonCallback(func(_ bool, _ *modal.InfoModal) bool {
			pageNavigator.CloseCurrentPage()
			return true
		})
	window.ShowModal(infoModal)
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *BackupInstructionsPage) Layout(gtx layout.Context) layout.Dimensions {
	sp := components.SubPage{
		Load:       pg.Load,
		Title:      "Keep in mind",
		WalletName: pg.wallet.Name,
		BackButton: pg.backButton,
		Back: func() {
			promptToExit(pg.Load, pg.ParentNavigator(), pg.ParentWindow())
		},
		Body: func(gtx C) D {
			return pg.infoList.Layout(gtx, len(pg.checkBoxes), func(gtx C, i int) D {
				return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, pg.checkBoxes[i].Layout)
			})
		},
	}

	pg.viewSeedBtn.SetEnabled(pg.verifyCheckBoxes())

	layout := func(gtx C) D {
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	isMobile := pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView)
	return container(gtx, isMobile, *pg.Theme, layout, "", pg.viewSeedBtn)
}

func (pg *BackupInstructionsPage) verifyCheckBoxes() bool {
	for _, cb := range pg.checkBoxes {
		if !cb.CheckBox.Value {
			return false
		}
	}
	return true
}

func container(gtx C, isMobile bool, theme cryptomaterial.Theme, body layout.Widget, infoText string, actionBtn cryptomaterial.Button) D {
	bodyLayout := func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return body(gtx)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = gtx.Constraints.Max
				return cryptomaterial.LinearLayout{
					Width:       cryptomaterial.MatchParent,
					Height:      cryptomaterial.WrapContent,
					Orientation: layout.Vertical,
					Direction:   layout.S,
					Alignment:   layout.Baseline,
					Background:  theme.Color.Surface,
					Shadow:      theme.Shadow(),
					Padding:     layout.UniformInset(values.MarginPadding16),
					Margin:      layout.Inset{Left: values.Size0_5},
					Border:      cryptomaterial.Border{Radius: cryptomaterial.Radius(4)},
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if !components.StringNotEmpty(infoText) {
							return D{}
						}
						label := theme.Label(values.TextSize14, infoText)
						label.Color = theme.Color.GrayText1
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, label.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return actionBtn.Layout(gtx)
					}))
			}),
		)
	}
	if isMobile {
		return components.UniformMobile(gtx, false, false, bodyLayout)
	}
	return components.UniformPadding(gtx, bodyLayout)
}
