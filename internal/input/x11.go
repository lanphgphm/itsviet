package input

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/xgb/proto"
	"github.com/BurntSushi/xgb/xprop"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbxwindow"

	"itsvietbaby/internal/processor"
)

type X11InputHandler struct {
	X				*xgbutil.XUtil
	processor		*processor.VietnameseProcessor
	browserClasses	[]string
	enabled			bool
	activeModifier	bool
	root			xproto.Window
}

func NewX11InputHandler(
	proc *processor.VietnameseProcessor, 
	browserClasses []string) (*X11InputHandler, error) {
		X, err := xgbutil.NewConn()
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to X server, %v", err)
		}

		keybind.Initialize(X)

		root := xproto.Setup(X.Conn()).DefaultScreen(X.Conn()).Root

		return &X11InputHandler {
			X:				X, 
			processor: 		proc, 
			browserClasses: browserClasses, 
			enabled: 		true, 
			activeModifier: false, 
			root: 			root,
		}, nil
}

func (h *X11InputHandler) Start() error {
	if err := keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			h.enabled = !h.enabled
			status := "enabled"
			if !h.enabled {
				status = "disabled"
			}
			log.Printf("Vietnamese typing %s", status)
		},
	)
}
