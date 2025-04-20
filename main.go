package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

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

func grabKeyboard(X *xgbutil.XUtil) {
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
}

func setupGlobalKeyCapturing(X *xgbutil.XUtil, vp *VProcessor) {
	// GrabKeyboard makes it such that the application process always 
	// gets a hold of the keyboard, which means the key injection always 
	// injects to itsviet process --> unable to inject to focused window, 
	// despite correctly identifying the focused window. 
	// --> temporarily release (Ungrab) keyboard for pasting in inject()
	grabKeyboard(X)

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


// xtest.FakeInput() requires a keycode (as if presented on the physical
// keyboard) to "fake" the keypress event. Since most Vietnamese characters
// do not have such a key on the standard keyboard layout, this is not 
// suitable for interception --> use X clipboard approach to inject w/o keycode
func inject(X *xgbutil.XUtil, text string, window xproto.Window) {
	fmt.Printf("Injecting: %s\n", text)

	// // backup original clipboard 
	// var originalClipboard string 
	// backupCmd := exec.Command("xclip", "-selection", "clipboard", "-o")
	// backupOutput, err := backupCmd.Output()
	// if err == nil {
	// 	originalClipboard = string(backupOutput)
	// } else {
	// 	panic(fmt.Sprintf("Failed to backup clipboard: %v\n", err))
	// }

	// setup clipboard to hold text before injection 
	cmd := exec.Command("xclip", "-selection", "clipboard", "-i")
	cmd.Stdin = strings.NewReader(text)
	err := cmd.Run() 
	if err != nil {
		panic(fmt.Sprintf("Failed to borrow clipboard: %v\n", err))
	}

	// Release keyboard grab to paste to `window`
	xproto.UngrabKeyboard(X.Conn(), xproto.TimeCurrentTime)
	time.Sleep(23*time.Millisecond) // wait for ungrab keyboard done

	// simulate Ctrl+V to paste injected text 
	ctrlKeycode := keybind.StrToKeycodes(X, "Control_L")[0]
    vKeycode := keybind.StrToKeycodes(X, "v")[0]
    xtest.FakeInput(X.Conn(), xproto.KeyPress, byte(ctrlKeycode), 0, window, 0, 0, 0)
    xtest.FakeInput(X.Conn(), xproto.KeyPress, byte(vKeycode), 0, window, 0, 0, 0)
    xtest.FakeInput(X.Conn(), xproto.KeyRelease, byte(vKeycode), 0, window, 0, 0, 0)
    xtest.FakeInput(X.Conn(), xproto.KeyRelease, byte(ctrlKeycode), 0, window, 0, 0, 0)
    
	time.Sleep(24*time.Millisecond) // wait for pasting to be done 

	// // restore original clipboard if there was one
    // if originalClipboard != "" {
    //     restoreCmd := exec.Command("xclip", "-selection", "clipboard", "-i")
    //     restoreCmd.Stdin = strings.NewReader(originalClipboard)
    //     restoreCmd.Run()
    // }

	grabKeyboard(X)
}
