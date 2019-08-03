# Notes
## Chip8

### References

[http://devernay.free.fr/hacks/chip8/C8TECH10.HTM]

[https://en.wikipedia.org/wiki/CHIP-8]

## Graphics

### SDL2

Set surface pixels.

Use rects like this: [https://github.com/berserkingyadis/chippyJ/blob/master/src/main/java/app/Display.java#L38] 

Use a scale factor for w and h.

FillRect docs:
 - [https://godoc.org/github.com/veandco/go-sdl2/sdl#Surface.FillRect]
 - [https://godoc.org/github.com/veandco/go-sdl2/sdl#Rect]

#### Links:

[https://wiki.libsdl.org/SDL_Surface]

[https://wiki.libsdl.org/SDL_PixelFormat]

[https://stackoverflow.com/questions/20070155/how-to-set-a-pixel-in-a-sdl-surface]

[https://stackoverflow.com/questions/48734398/how-can-i-read-write-pixels-in-a-sdl-surface-in-sdl2]

[https://gamedev.stackexchange.com/questions/38186/trying-to-figure-out-sdl-pixel-manipulation]

[https://stackoverflow.com/questions/6852055/how-can-i-modify-pixels-using-sdl]

## TODO
- [ ] Implement timer module.
    - [https://golang.org/pkg/time/#Ticker]
- [ ] Finish implementing instructions.
- [ ] Implement interfaces for graphics and input modules, will make it easier to test.
    - Reference [https://gist.github.com/jorygeerts/e887856cc15b64cb9681639cd83c4a37]
