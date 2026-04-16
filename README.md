# Cato Linux Status Monitor

A lightweight, Pure Go system tray monitor for the Cato SDP Linux client.
It Provides a traffic-light indicator for your Cato status and integrates with the native Linux desktop notifications via D-Bus.


## Features 
* **Status Icons:** Instantly know if you are connected (Green), connected but need to update (Blue), need to authenticate (yellow), or are disconnected (Red).
* **Native Notifications:** Alerts you of state changes using the native Linux D-Bus notifications- no external CLI tolls required.
* **Resource Efficient:** Built in Go, runs quietly in the background without draining system resources.

## Prerequisites
This application requuires the standard Linux app indicator runtime libraries. 
On Debian/Ubuntu-base systems, you can ensure you have them by running the following with sudo or as root:
`apt install libayatana-appindicator3-1`

## Installation

1. Head over to the [Releases](../../releases) page and download the latest `cato-monitor-linux-amd64.tar.gz` .
2. Extract the archive into a permanent folder in your home directory (e.g., `~/cato-monitor`).
    ```bash
    mkdir -p ~/cato-monitor
    tar -xzf cato-monitor-linux-amd64.tar.gz -C ~/cato-monitor
    ```

## Configure Autostart
1. Create a `.desktop` file in your user's autostart directory
    ```bash
    nano ~/.config/autostart/cato-monitor.desktop
    ```

2. Paste the following configuration. Important: Replace `/home/YOUR_USERNAME/cato-monitor/` with the actute path to where you extracted the files. (Do not use `~` in the path here).
    ```bash
    [Desktop Entry]
    Type=Application
    Exec=/home/YOUR_USERNAME/cato-monitor/cato-monitor
    Path=/home/YOUR_USERNAME/cato-monitor/
    Hidden=false
    NoDisplay=false
    X-GNOME-Autostart-enabled=true
    Name=Cato Monitor
    Comment=System tray monitor for Cato SDP
    ```

3. Save the file. Logout and back in, and the monitor will start automatically.


## Customisation (Custom Icons)
By default, this application uses standard traffic-light icons that are embeded directly into the single executable file.
However, you can easily override these with your own custom icoms!

Recomened size is 256x256.

To use custom icons:
1. Create a folder named `icons` in the exact same directory where your `cato-monitor` executable lives.
2. Add your custom transparent PNG files into that folder. They must be named exactly:
    * `green.png` (Connected)
    * `blue.png` (update Available)
    * `yellow.png` (Authentication Needed)
    * `red.png` (Disconnected)

When the monitor starts, it checks this local `icons` folder first. if it finds your custom images, it will use them (perfect for custom system themes or Tux penguins!). if a file is missing or the folder dosen't exist, it seamlessly falls back to the embedded default icons.

any custom Icons I have created or that have been submitted will be put in the folder called `customIcons`