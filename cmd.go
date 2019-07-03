package main

import "fmt"

// Command Mappings
var commandMap = map[string]string{
	"blue":   "blue_led",
	"yellow": "yellow_led",
	"ir":     "ir_led",
	"ircut":  "ir_cut",
	"mtr":    "motor",
	"http":   "http_server",
	"rtsp":   "rtsp_h264_server",
	"motion": "motion_detection",
	"reboot": "reboot_system",
}

// Custom Commands
var customCommands = []string{
	"snapshot",
	"snap",
}

// executeCustomCommand executes custom commands
// TODO: implement video recording custom command using ffmpeg
// TODO: implement motor calibrate custom command
// TODO: implement ftp status checking custom command[on off status]
// TODO: implement status custom command; gets status of of various settings
func executeCustomCommand(selectedCam, username, password, command string) (err error) {

	if command == "snapshot" || command == "snap" {
		// execute snapshot command
		err = executeSSH(selectedCam, username, password, []string{"snapshot", "exit"})
		if err != nil {
			return err
		}
		// fetch snapshot
		err = getFile(username, password, selectedCam, "/tmp/snapshot.jpg")
		if err != nil {
			return err
		}

		return
	}

	return fmt.Errorf("custom command(%s) not found", command)
}

// Official Commands
var commands = []string{
	"blue_led on",
	"blue_led off",
	"blue_led status",

	"yellow_led on",
	"yellow_led off",
	"yellow_led status",

	"ir_led on",
	"ir_led off",
	"ir_led status",

	"ir_cut on",
	"ir_cut off",
	"ir_cut status",

	// Motor [direction] commands may have a numerical value attached
	// If no numerical values given, default is 100
	// e.g. motor up 100; motor left 30; motor down
	"motor up",
	"motor down",
	"motor left",
	"motor right",
	"motor status", // vertical
	"motor status horizontal",

	"http_server on",
	"http_server off",
	"http_server restart",
	"http_server status",

	"http_password", // set new http password [http_password new_password]

	"rtsp_h264_server on",
	"rtsp_h264_server off",
	"rtsp_h264_server status",

	"rtsp_mjpeg_server on",
	"rtsp_mjpeg_server off",
	"rtsp_mjpeg_server status",

	"motion_detection on",
	"motion_detection off",
	"motion_detection status",

	"motion_send_mail on",
	"motion_send_mail off",
	"motion_send_mail status",

	"motion_send_telegram on",
	"motion_send_telegram off",
	"motion_send_telegram status",

	"motion_tracking on",
	"motion_tracking off",
	"motion_tracking status",

	"night_mode on",
	"night_mode off",
	"night_mode status",

	"auto_night_mode on",
	"auto_night_mode off",
	"auto_night_mode status",

	// snapshots saved to /tmp/snapshot.jpeg
	"snapshot", // makes call to get_snapshot custom command after taking snapshot

	"reboot_system",
}

// enableFTP enables FTP on cam
func enableFTP(username, password, ip string) error {
	commands := []string{
		"/system/sdcard/controlscripts/ftp_server start", // start bftpd daemon
		"exit",
	}
	err := executeSSH(ip, username, password, commands)
	if err != nil {
		return err
	}
	return nil
}
