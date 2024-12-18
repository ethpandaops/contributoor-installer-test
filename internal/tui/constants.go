package tui

import "github.com/gdamore/tcell/v2"

// Logo is used as a bit of branding in the installer UI.
const Logo string = `
   ______            __       _ __          __                  
  / ____/___  ____  / /______(_) /_  __  __/ /_____  ____  _____
 / /   / __ \/ __ \/ __/ ___/ / __ \/ / / / __/ __ \/ __ \/ ___/
/ /___/ /_/ / / / / /_/ /  / / /_/ / /_/ / /_/ /_/ / /_/ / /    
\____/\____/_/ /_/\__/_/  /_/_.___/\__,_/\__/\____/\____/_/    
`

// NetworkOption is used to represent a network option, eg: Ethereum Mainnet, Sepolia Testnet, etc.
type NetworkOption struct {
	Label       string
	Value       string
	Description string
}

// AvailableNetworks is a list of available networks.
var AvailableNetworks = []NetworkOption{
	{
		Label:       "Ethereum Mainnet",
		Value:       "mainnet",
		Description: "This is the real Ethereum main network.",
	},
	{
		Label:       "Holesky Testnet",
		Value:       "holesky",
		Description: "The Holesky test network.",
	},
	{
		Label:       "Sepolia Testnet",
		Value:       "sepolia",
		Description: "The Sepolia test network.",
	},
}

// OutputServerOption is used to represent an output server option, eg: ethPandaOps Production, ethPandaOps Staging, etc.
type OutputServerOption struct {
	Label       string
	Value       string
	Description string
}

// AvailableOutputServers is a list of available output servers.
var AvailableOutputServers = []OutputServerOption{
	{
		Label:       "ethPandaOps Production",
		Value:       "https://xatu.primary.production.platform.ethpandaops.io",
		Description: "The production server provided by ethPandaOps.",
	},
	{
		Label:       "ethPandaOps Staging",
		Value:       "https://xatu.primary.staging.platform.ethpandaops.io",
		Description: "The staging server provided by ethPandaOps.",
	},
	{
		Label:       "Custom",
		Value:       "custom",
		Description: "Use your own custom output server.",
	},
}

// Colors used throughout the UI.
var (
	ColorBackground      = tcell.ColorDarkSlateGray
	ColorFormBackground  = tcell.ColorLightSlateGray
	ColorBorder          = tcell.ColorWhite
	ColorButtonActivated = tcell.ColorYellow
	ColorButtonText      = tcell.ColorBlack
	ColorError           = tcell.ColorRed
	ColorSuccess         = tcell.ColorGreen
	ColorHeading         = tcell.ColorYellow
)

// Common strings used in the various UI screens.
const (
	ButtonSaveSettings = "Save Settings"
	ButtonClose        = "Close"
	ButtonNext         = "Next"
	ButtonTryAgain     = "Try Again"
	TitleDescription   = "Description"
	TitleSettings      = "Settings"
)