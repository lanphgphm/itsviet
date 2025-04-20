package main

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgb/xtest"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

const (
	KEY_SPACE = " "
)

func main () {
	// initialize new connection to X11 server 
	X, err := xgbutil.NewConn()
	if err != nil {
		panic(fmt.Sprintf("Cannot establish connection with X11 server: %v", err))
	}
	defer X.Conn().Close()

	err = xtest.Init(X.Conn())
	if err != nil {
		panic(fmt.Sprintf("Cannot initialize XTEST extension: %v", err))
	}

	// X connection success --> initialize keybind listener 
	keybind.Initialize(X)
	vp := NewVProcessor()
	// captures key globally 
	setupGlobalKeyCapturing(X, vp)

	xevent.Main(X)
}

func getFocusedWindow(X *xgbutil.XUtil) xproto.Window {
	reply, err := xproto.GetInputFocus(X.Conn()).Reply()
	if err != nil {
		fmt.Printf("Failed to get focused window: %v\nReturning root window of process.", err)
		return X.RootWin()
	}
	return reply.Focus
}

func setupGlobalKeyCapturing(X *xgbutil.XUtil, vp *VProcessor){
	cookie := xproto.GrabKeyboard(
		X.Conn(), 
		true, 
		X.RootWin(),
		xproto.TimeCurrentTime,
		xproto.GrabModeAsync, 
		xproto.GrabModeAsync,
	)

	reply, err := cookie.Reply()
	if err != nil {
		panic(fmt.Sprintf("Grab keyboard failed: %v", err))
	}
	if reply.Status != xproto.GrabStatusSuccess {
		panic(fmt.Sprintf("Grab keyboard failed: %v", err))
	}

	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			focusedWindow := getFocusedWindow(X)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			modStr := keybind.ModifierString(e.State)

			fmt.Printf("Key pressed: %s (with modifiers: %s)\n", keyStr, modStr)
			
			if e.State&xproto.ModMask4 != 0 && keyStr == KEY_SPACE {
				// Super+Space to toggle vietnamese typing
				vp.Toggle()
				if vp.enabled {
					fmt.Println("Vietnamese typing enabled. itsvietbaby!")
				} else {
					fmt.Println("Vietnamese typing disabled.")
				}
				return 
			}

			intercepted, transformedText := vp.Process(keyStr, byte(e.Detail), e.State)
			if intercepted {
				inject(X, transformedText, focusedWindow)

				// use ReplayKeyboard to prevent sending original key (effectively intercepting)
				xproto.AllowEvents(X.Conn(), xproto.AllowReplayKeyboard, xproto.TimeCurrentTime)
			} else {
				xproto.AllowEvents(X.Conn(), xproto.AllowSyncKeyboard, xproto.TimeCurrentTime)
			}
		}).Connect(X, X.RootWin())
}

var unicodeToKeysym = map[rune]string {
	'â': "acircumflex",
} 

// xtest.FakeInput() requires a keycode (as if presented on the physical
// keyboard) to "fake" the keypress event. Since most Vietnamese characters
// do not have such a key on the standard keyboard layout, this is not 
// suitable for interception
func inject(X *xgbutil.XUtil, text string, window xproto.Window) {
	fmt.Printf("Injecting: %s\n", text)

	for _, char := range text {
		symName, ok := unicodeToKeysym[char]

		if ok {
			keycodes := keybind.StrToKeycodes(X, symName)
			if len(keycodes) > 0 {
				keycode := keycodes[0]
				xtest.FakeInput(X.Conn(), xproto.KeyPress, byte(keycode), 0, window, 0, 0, 0)
                xtest.FakeInput(X.Conn(), xproto.KeyRelease, byte(keycode), 0, window, 0, 0, 0)
            } else {
                fmt.Printf("No keycode found for: %s\n", symName)
            }
		} else {
			keycodes := keybind.StrToKeycodes(X, string(char))
            if len(keycodes) > 0 {
                keycode := keycodes[0]
                xtest.FakeInput(X.Conn(), xproto.KeyPress, byte(keycode), 0, window, 0, 0, 0)
                xtest.FakeInput(X.Conn(), xproto.KeyRelease, byte(keycode), 0, window, 0, 0, 0)
            } else {
                fmt.Printf("No keycode found for character: %s (Unicode: %U)\n", string(char), char)
            }
		}
	}
}
