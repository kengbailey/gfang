package main

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
var customCommandsMap = map[string]string{
	"snap":      "get_snapshot",    // retrieves snapshot from cam and places in current directory
	"calibrate": "motor calibrate", // calls command: motor reset_pos_count
	"ftpon":     "enable ftp",      // starts bftpd service on cam
}

// executeCustomCommand executes custom commands
// TODO: implement snapshot/snapfetch + video recording custom commands
func executeCustomCommand(selectedCam, username, password, command string) (err error) {

	// check for mapped custom command
	for k := range customCommandsMap {
		if k == command {
			commands := []string{customCommandsMap[command], "exit"}
			err = executeSSH(selectedCam, username, password, commands)
			if err != nil {
				return err
			}
			return
		}
	}

	if command == "snapshot" {

	}
	return
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
