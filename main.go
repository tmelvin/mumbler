package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/bmmcginty/go-openal/openal"
	"github.com/cantudo/barnard/gumble/gumble"
	_ "github.com/cantudo/barnard/gumble/opus"
	"github.com/cantudo/barnard/uiterm"
)

func show_devs(name string, args []string) {
	if args == nil {
		fmt.Printf("no items for %s\n", name)
	}
	fmt.Printf("%s\n", name)
	for i := 0; i < len(args); i++ {
		fmt.Printf("%d: %s\n", i, args[i])
	}
}

func do_list_devices() {
	idevs := openal.GetStrings(openal.CaptureDeviceSpecifier)
	show_devs("Inputs:", idevs)
	odevs := openal.GetStrings(openal.AllDevicesSpecifier)
	if odevs != nil && len(odevs) > 0 {
		show_devs("All outputs:", odevs)
	} else {
		odevs = openal.GetStrings(openal.DeviceSpecifier)
		show_devs("All outputs:", odevs)
	}
}

func main() {
	// Command line flags
	server := flag.String("server", "localhost:64738", "the server to connect to")
	username := flag.String("username", "", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	insecure := flag.Bool("insecure", false, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")
	channel := flag.String("channel", "", "channel you would connect")
	inputdevice := flag.Int("inputdevice", 0, "input device to use, see list_devices")
	outputdevice := flag.Int("outputdevice", 0, "output device to use, see list_devices")
	list_devices := flag.Bool("list_devices", false, "do not connect; instead, list available audio devices and exit")
	inmediatestart := flag.Bool("inmediatestart", false, "start transmitting right away")

	flag.Parse()

	if *list_devices {
		do_list_devices()
		os.Exit(0)
	}

	idevs := openal.GetStrings(openal.CaptureDeviceSpecifier)
	odevs := openal.GetStrings(openal.AllDevicesSpecifier)

	// Initialize
	b := Barnard{
		Config:         gumble.NewConfig(),
		Address:        *server,
		Channel:        *channel,
		InputDevice:    idevs[*inputdevice],
		OutputDevice:   odevs[*outputdevice],
		InmediateStart: *inmediatestart,
	}

	b.Config.Username = *username
	b.Config.Password = *password

	if *insecure {
		b.TLSConfig.InsecureSkipVerify = true
	}
	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		b.TLSConfig.Certificates = append(b.TLSConfig.Certificates, cert)
	}

	b.Ui = uiterm.New(&b)
	b.Ui.Run()
}
