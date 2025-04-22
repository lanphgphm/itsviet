# NamKey*: Vietnamese typing interceptor 

I can't get ibus-bamboo to work on my machine\*\*, and i thought to myself "How hard can it be?"\*\*\*. So this is a Vietnamese typing interceptor I wrote for myself. 

The idea is to start a daemon process that interferes with my keyboard input and send the processed Vietnamese word (with diacritics) to whatever application I'm typing to. 

## Core component breakdown && their interfaces
1. X11 input interceptor 
- Input: raw keystrokes 
- Output: intercepted (==converted to Vietnamese) keys to send to applications

2. Vietnamese processor 
- Input: sequence of raw keystrokes 
- Output: Vietnamese character with diacritics 

## Flow 
Keyboard 
-> X11 Server (key capturing)
-> itsviet interceptor 
-> Vietnamese processor 
-> X11 Server (key injecting)
-> Application

## Relevant functions from `xgbutils` Go module
- starting point: 
    - `Main()`: Main X event loop 
    - `NewConn()`: open new connection to X server 
- xevent package: 
    - `KeyPressFun()`: inject function to handle keypress event 
    - `KeyReleaseFun()`: inject function to handle key release event 
- keybind package: to grab key events 
    - `LookupString()`
    - `ModifierString()` 
- ~~xtest package: to simulate keypress event with processed Vietnamese characters~~
    - ~~`FakeKeyEvent()`~~
- borrow clipboard from "xclip -sel c -i" to buffer and inject interepted sequences

## Helpful resources 
- Wikipedia on X11: [X11](https://en.wikipedia.org/wiki/X_Window_System) -- especially the Software Architecture part 
- Before this, get familiar with X11 functions through the tutorials in [xgbutil repository](https://github.com/BurntSushi/xgbutil)

**Footnote**:   
\* The name is a play off of [VietKey](https://vi.wikipedia.org/wiki/Vietkey)    
\*\* The maintainer of ibus-bamboo also uses NixOS, so I'm pretty sure I'm 2 configurations away from making it work. But there is more character development in writing my own, so I'm doing it for the plot.    
\*\*\* Stay tuned to see my answer to this question :) 
