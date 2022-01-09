package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"time"

	"net/http"

	"github.com/bmmcginty/go-openal/openal"
	"github.com/cantudo/barnard/gumble/gumble"
	_ "github.com/cantudo/barnard/gumble/opus"

	// "github.com/cantudo/barnard/uiterm"
	"github.com/gin-gonic/gin"
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
	if odevs != nil {
		show_devs("All outputs:", odevs)
	} else {
		odevs = openal.GetStrings(openal.DeviceSpecifier)
		show_devs("All outputs:", odevs)
	}
}

func connect(username string, password string, server string, channel string,
	inputdevice int, outputdevice int, insecure bool, certificate string,
	inmediatestart bool) (*Barnard, error) {

	idevs := openal.GetStrings(openal.CaptureDeviceSpecifier)
	odevs := openal.GetStrings(openal.AllDevicesSpecifier)

	if inputdevice >= len(idevs) {
		return nil, fmt.Errorf("invalid input device")
	}
	if outputdevice >= len(odevs) {
		return nil, fmt.Errorf("invalid output device")
	}

	// Initialize
	b := Barnard{
		Config:         gumble.NewConfig(),
		Address:        server,
		Channel:        channel,
		InputDevice:    idevs[inputdevice],
		OutputDevice:   odevs[outputdevice],
		InmediateStart: inmediatestart,
	}

	b.Config.Username = username
	b.Config.Password = password

	if insecure {
		b.TLSConfig.InsecureSkipVerify = true
	}
	if certificate != "" {
		cert, err := tls.LoadX509KeyPair(certificate, certificate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return &b, err
		}
		b.TLSConfig.Certificates = append(b.TLSConfig.Certificates, cert)
	}

	fmt.Printf("Initializing...\n")

	// b.Ui = uiterm.New(&b)
	// b.Ui.Run()
	// Time to initialize
	c1 := make(chan error, 1)
	go func() {
		c1 <- b.start()
	}()

	select {
	case err := <-c1:
		return &b, err
	case <-time.After(3 * time.Second):
		return &b, fmt.Errorf("connection timeout")
	}
}

type client struct {
	Server       string `json:"server"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Channel      string `json:"channel"`
	InputDevice  int    `json:"inputdevice"`
	OutputDevice int    `json:"outputdevice"`
}

type devices struct {
	Inputs  []string `json:"inputs"`
	Outputs []string `json:"outputs"`
}

var clients = []client{}
var clientBarnardMap map[string]*Barnard = make(map[string]*Barnard)

func getClients(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, clients)
}

func getDevices(c *gin.Context) {
	var d devices
	d.Inputs = openal.GetStrings(openal.CaptureDeviceSpecifier)
	d.Outputs = openal.GetStrings(openal.AllDevicesSpecifier)

	c.IndentedJSON(http.StatusOK, d)
}

// postClients adds a client from JSON received in the request body.
func postClients(c *gin.Context) {
	var newClient client

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newClient); err != nil {
		return
	}
	id := fmt.Sprintf("%s@%s", newClient.Username, newClient.Server)
	// Add the new album to the slice.
	if b_test := clientBarnardMap[id]; b_test != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"status": "already connected"})
		return
	}
	b, err := connect(
		newClient.Username, newClient.Password, newClient.Server,
		newClient.Channel, newClient.InputDevice, newClient.OutputDevice,
		true, "", true,
	)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error connecting to server: %s", err))
	} else {
		clients = append(clients, newClient)

		clientBarnardMap[id] = b
		c.IndentedJSON(http.StatusCreated, newClient)
	}
}

func disconnectClient(c *gin.Context) {
	id := c.Param("id")

	b := clientBarnardMap[id]
	if b != nil {
		b.Stream.Destroy()
		b.Client.Disconnect()
		delete(clientBarnardMap, id)
	} else {
		c.String(http.StatusNotFound, "Client not found")
	}
	c.String(http.StatusOK, fmt.Sprintf("Disconnected %s", id))
}

func api(port int) {
	router := gin.Default()
	router.GET("/devices", getDevices)
	router.GET("/clients", getClients)
	router.POST("/clients", postClients)
	router.POST("/clients/disconnect/:id", disconnectClient)

	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
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
	servermode := flag.Bool("servermode", false, "start in server mode")
	serverport := flag.Int("serverport", 8001, "port to listen on for server mode")

	flag.Parse()

	if !*servermode {
		if *list_devices {
			do_list_devices()
			os.Exit(0)
		}
		_, err := connect(
			*username, *password, *server, *channel,
			*inputdevice, *outputdevice, *insecure, *certificate,
			*inmediatestart)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		api(*serverport)
	}
}
