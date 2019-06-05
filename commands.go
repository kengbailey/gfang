package main

var commands []string{
	"blue_led on",
	"blue_led off",
	"blue_led status",

	"yellow_led on",
	"yellow_led off",
	"yellow_led status",

	"ir_led on ",
	"ir_led off ",
	"ir_led status",

	"ir_cut on ",
	"ir_cut off ",
	"ir_cut status",
	
	// all motor [direction] commands must have a numerical value attached 
	// e.g. motor up 100 
	"motor up ",
	"motor down ",
	"motor left ",
	"motor right ",
	"motor status",
	"motor status horizontal",

	"http_server on ",
	"http_server off",
	"http_server restart",
	"http_server status",

	"http_password ", // set new http password [http_password new_password]

	"rtsp_h264_server on",
	"rtsp_h264_server off ",
	"rtsp_h264_server status",

	"rtsp_mjpeg_server on",
	"rtsp_mjpeg_server off ",
	"rtsp_mjpeg_server status ",

	"motion_detection on ",
	"motion_detection off ",
	"motion_detection status ",

	"motion_send_mail on",
	"motion_send_mail off ",
	"motion_send_mail status ",

	"motion_send_telegram on",
	"motion_send_telegram off ",
	"motion_send_telegram status", 

	"motion_tracking on",
	"motion_tracking off ",
	"motion_tracking status ",

	"night_mode on ",
	"night_mode off ",
	"night_mode status ",

	"auto_night_mode on",
	"auto_night_mode off ",
	"auto_night_mode status",

	// snapshots saved to /tmp/snapshot.jpeg
	"snapshot",

	"reboot_system",
}