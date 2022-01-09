package main

import (
	"fmt"
	"net"
	"os"

	"github.com/cantudo/barnard/gumble/gumble"
	"github.com/cantudo/barnard/gumble/gumbleopenal"
	"github.com/cantudo/barnard/gumble/gumbleutil"
)

func (b *Barnard) start() error {
	b.Config.Attach(gumbleutil.AutoBitrate)
	b.Config.Attach(b)

	// Audio
	if os.Getenv("ALSOFT_LOGLEVEL") == "" {
		os.Setenv("ALSOFT_LOGLEVEL", "0")
	}
	var err error
	_, err = gumble.DialWithDialer(new(net.Dialer), b.Address, b.Config, &b.TLSConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		// os.Exit(1)
		return err
	}
	if stream, err := gumbleopenal.New(b.Client, b.InputDevice, b.OutputDevice); err != nil {
		fmt.Printf("Error starting stream\n")
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	} else {
		b.Stream = stream
		if b.InmediateStart {
			b.Stream.StartSource()
		}
	}

	return nil
}

func (b *Barnard) OnConnect(e *gumble.ConnectEvent) {
	b.Client = e.Client

	// If there is a channel on arguments move it
	if b.Channel != "" {
		target := e.Client.Self.Channel.Find(b.Channel)
		if target != nil {
			e.Client.Self.Move(e.Client.Self.Channel.Find(b.Channel))
		} else if b.Ui != nil {
			b.AddOutputLine(fmt.Sprintf("Could not connect to %s, moving to "+
				"default channel %s", b.Channel, e.Client.Self.Channel.Name))
		}
	}

	if b.Ui != nil {
		b.Ui.SetActive(uiViewInput)
		b.UiTree.Rebuild()
		b.Ui.Refresh()
		b.UpdateInputStatus(fmt.Sprintf("To: %s", e.Client.Self.Channel.Name))
		b.AddOutputLine(fmt.Sprintf("Connected to %s", b.Client.Conn.RemoteAddr()))
		if e.WelcomeMessage != nil {
			b.AddOutputLine(fmt.Sprintf("Welcome message: %s", esc(*e.WelcomeMessage)))
		}
	}
}

func (b *Barnard) OnDisconnect(e *gumble.DisconnectEvent) {
	var reason string
	switch e.Type {
	case gumble.DisconnectError:
		reason = "connection error"
	}
	if b.Ui != nil {
		if reason == "" {
			b.AddOutputLine("Disconnected")
		} else {
			b.AddOutputLine("Disconnected: " + reason)
		}
		b.UiTree.Rebuild()
		b.Ui.Refresh()
	}
}

func (b *Barnard) OnTextMessage(e *gumble.TextMessageEvent) {
	b.AddOutputMessage(e.Sender, e.Message)
}

func (b *Barnard) OnUserChange(e *gumble.UserChangeEvent) {
	if e.Type.Has(gumble.UserChangeChannel) && e.User == b.Client.Self {
		b.UpdateInputStatus(fmt.Sprintf("To: %s", e.User.Channel.Name))
	}
	b.UiTree.Rebuild()
	b.Ui.Refresh()
}

func (b *Barnard) OnChannelChange(e *gumble.ChannelChangeEvent) {
	b.UiTree.Rebuild()
	b.Ui.Refresh()
}

func (b *Barnard) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	var info string
	switch e.Type {
	case gumble.PermissionDeniedOther:
		info = e.String
	case gumble.PermissionDeniedPermission:
		info = "insufficient permissions"
	case gumble.PermissionDeniedSuperUser:
		info = "cannot modify SuperUser"
	case gumble.PermissionDeniedInvalidChannelName:
		info = "invalid channel name"
	case gumble.PermissionDeniedTextTooLong:
		info = "text too long"
	case gumble.PermissionDeniedTemporaryChannel:
		info = "temporary channel"
	case gumble.PermissionDeniedMissingCertificate:
		info = "missing certificate"
	case gumble.PermissionDeniedInvalidUserName:
		info = "invalid user name"
	case gumble.PermissionDeniedChannelFull:
		info = "channel full"
	case gumble.PermissionDeniedNestingLimit:
		info = "nesting limit"
	}
	b.AddOutputLine(fmt.Sprintf("Permission denied: %s", info))
}

func (b *Barnard) OnUserList(e *gumble.UserListEvent) {
}

func (b *Barnard) OnACL(e *gumble.ACLEvent) {
}

func (b *Barnard) OnBanList(e *gumble.BanListEvent) {
}

func (b *Barnard) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (b *Barnard) OnServerConfig(e *gumble.ServerConfigEvent) {
}
