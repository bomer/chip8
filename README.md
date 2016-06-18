#Golang Chip Emulator!

My first Emulator. Learning the basic of op code decoding and cpu/cycle emumlation. Couldn't have done it without Bisqwit's youtube video's inspiring me that it's possible to actually write one and MultiGesture's tutorial on the basic loop and how to interpret the op code's using bitwise operations

The GPU Rendering I wrote like 3 different times. Finally I settled on GoMobile OpenGL ES2's Image class which creates a texture that gets bound in 2d to the front buffer. Worked really well, much more performant than my first attempts.  

Also, OpenGL ES2 sucks and is not OpenGL at all. But made do eventually.

![ScreenShot](https://raw.githubusercontent.com/bomer/chip8/master/brix.png)

##Run with

go run main.go ROMNAME

If no rom name is presnt a default is loaded.

References:

1-Wikipedia 

https://en.wikipedia.org/wiki/CHIP-8#Virtual_machine_description

2 - Awesome Tutorial

http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/

3 - Bisqwit - The youtube God that inspired me that I might be able to write an emulator.

https://www.youtube.com/user/Bisqwit
