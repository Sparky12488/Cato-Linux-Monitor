package main

import (
	"embed"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus/v5"
)

//go:embed icons/*.png
var iconFS embed.FS

type CatoState int

const (
	StateUnknown CatoState = iota
	StateConnected
	StateConnectedUpdate
	StateDisconnected
	StateAuthNeeded
)

var currentState = StateUnknown

var (
	iconRed    []byte
	iconYellow []byte
	iconGreen  []byte
	iconBlue   []byte
)

var mStatus *systray.MenuItem

func main() {
	// Load icons with our new fallback function
	iconRed = loadIcon("icons/red.png")
	iconYellow = loadIcon("icons/yellow.png")
	iconGreen = loadIcon("icons/green.png")
	iconBlue = loadIcon("icons/blue.png")
	systray.Run(onReady, onExit)
}

func loadIcon(path string) []byte {
	data, err := os.ReadFile(path)
	if err == nil {
		log.Printf("Loaded custom icon from filesystem: %s\n", path)
		return data
	}

	data, err = iconFS.ReadFile(path)
	if err != nil {
		log.Printf("Notice: Embedded icon %s not found\n", path)
	}
	return data
}

func onReady() {
	if len(iconRed) > 0 {
		systray.SetIcon(iconYellow)
	}
	systray.SetTitle("Cato")
	systray.SetTooltip("Initializing Cato Monitor...") // Keeps Windows/Mac compatibility

	// --- ADDED THIS BLOCK ---
	// Create a non-clickable menu item to act as our "Hover" tooltip
	mStatus = systray.AddMenuItem("Initializing Cato Monitor...", "")
	mStatus.Disable()
	systray.AddSeparator()
	// ------------------------

	mQuit := systray.AddMenuItem("Quit", "Close the monitor")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go monitorCato()
}

func onExit() {

}

func monitorCato() {
	// Run infinitely
	for {
		out, err := exec.Command("cato-sdp", "status").Output()

		output := string(out)
		outputLower := strings.ToLower(output)

		if err != nil {
			log.Printf("Command failed with error: %v. Output was: %s\n", err, output)
			updateState(StateDisconnected)
		} else {

			isConnected := strings.Contains(output, "STATE_AUTHENTICATED")
			needsUpdate := strings.Contains(output, "New client version is available")

			if isConnected {
				if needsUpdate {
					updateState(StateConnectedUpdate)
				} else {
					updateState(StateConnected)
				}
			} else if strings.Contains(outputLower, "auth") || strings.Contains(outputLower, "login") || strings.Contains(outputLower, "code") {
				updateState(StateAuthNeeded)
			} else {
				updateState(StateDisconnected)
			}
		}

		// Sleep for 5 seconds BEFORE looping again
		time.Sleep(5 * time.Second)
	}
}

func updateState(newState CatoState) {
	if newState == currentState {
		return
	}

	currentState = newState

	switch newState {
	case StateConnected:
		if len(iconGreen) > 0 {
			systray.SetIcon(iconGreen)
		}
		systray.SetTooltip("Cato: Connected")
		mStatus.SetTitle("Status: Connected") // Updates the menu item
		sendNotification("Cato VPN", "Connected successfully", "network-transmit-receive")

	case StateConnectedUpdate:
		if len(iconBlue) > 0 {
			systray.SetIcon(iconBlue)
		}
		systray.SetTooltip("Cato: Update Available")
		mStatus.SetTitle("Status: Update Available") // Updates the menu item
		sendNotification("Cato Update Available", "A new version of Cato is ready.", "software-update-available")

	case StateDisconnected:
		if len(iconRed) > 0 {
			systray.SetIcon(iconRed)
		}
		systray.SetTooltip("Cato: Disconnected")
		mStatus.SetTitle("Status: Disconnected") // Updates the menu item
		sendNotification("Cato Alert", "VPN has disconnected", "network-error")

	case StateAuthNeeded:
		if len(iconYellow) > 0 {
			systray.SetIcon(iconYellow)
		}
		systray.SetTooltip("Cato: Auth Needed")
		mStatus.SetTitle("Status: Auth Needed") // Updates the menu item
		sendNotification("Cato Action Required", "Authentication is needed", "dialog-warning")
	}
}

func sendNotification(title, message, icon string) {
	// Connect to the session bus natively
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Printf("Failed to connect to D-Bus: %v\n", err)
		return
	}
	// We don't defer conn.Close() here because dbus.SessionBus() uses a shared connection

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	// The standard FreeDesktop Notification D-Bus signature
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"Cato Monitor",            // App Name
		uint32(0),                 // Replaces ID
		icon,                      // App Icon
		title,                     // Summary
		message,                   // Body
		[]string{},                // Actions
		map[string]dbus.Variant{}, // Hints
		int32(5000),               // Timeout in milliseconds (5 seconds)
	)

	if call.Err != nil {
		log.Printf("Failed to send D-Bus notification: %v\n", call.Err)
	}
}
