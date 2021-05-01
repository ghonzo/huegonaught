# Huegonaught

A utility to turn off [Hue](https://www.philips-hue.com/en-us) lights when some other light(s) become unreachable. The horribly embarrasing name of the project comes from:
* **Hue** - works with Hue lights
* **Go** - written in Go
* **Naught** - turns lights off

## History
This little project was born out of necessity. We obtained a [Hue HDMI Sync Box](https://www.philips-hue.com/en-us/p/hue-play-hdmi-sync-box-/046677555221) for the game room so we had a mix of:
* Two [Hue Play](https://www.philips-hue.com/en-us/p/hue-white-and-color-ambiance-play-light-bar-double-pack/7820230U7) bars behind the TV
* A Hue [Bloom](https://www.philips-hue.com/en-us/p/hue-white-and-color-ambiance-bloom-table-lamp/046677560188) light off to the side
* A number of Hue bulbs in ceiling can lights connected to a dumb mechanical light switch.

So what I wanted to have happen is: you walk into the room and you flip the (dumb) light switch on the wall, which turns on all the ceiling lights. Then when you turn on the gaming system, it causes the Sync Box to come to life and it controls all the lights in the room (including the two play bars and the Bloom light). So far, so good.

But then when you're done gaming, you turn off the TV and the Hue Sync Box stops syncing. You leave the room and turn off the light switch, which turns off the ceiling lamps, but the other "Hue-controlled" lights stay on. Since I want those other lights to turn off as well, I wrote this little utility program to turn them off once we detect that the ceiling lights have turned off (i.e., they become "unreachable").

That's what this project does.
## Usage
```
Usage of huegonaught:
  -bridge string
        the ip address of the Hue bridge (default "philips-hue")
  -controlled string
        (required) comma-separated list of light ids to turn off when signal bulbs are unreachable
  -off int
        the polling interval (in seconds) when the lights are off (default 10)
  -on int
        the polling interval (in seconds) when the lights are on (default 1)
  -signal string
        (required) comma-separated list of bulb ids to monitor for reachable status
  -user string
        (required) the user for Hue bridge
  -v    verbose
```
This program never exits unless:
* there is a missing command line argument
* the bridge cannot be initially
* the user is initially not valid

To obtain a valid `user`, see [Getting Started](https://developers.meethue.com/develop/get-started-2/) on the Hue Developers website.

Oh, this only works locally by the way. The bridge must be accessible via your local network. I'm sure there's a way you can make it work with remote access.
## Example
```
./huegonaught -bridge 192.168.1.11 -user p0LqF2Cd7spYFeassvebwuV0J26S2hPis73bbdXJ -signal 38,55,54,51,44,52,53 -controlled 46,47,48 -on 10 -off 60
```
## Thanks
Thanks to the [gohue](https://github.com/Collinux/gohue) library, which made an easy project even easier.
## Contributing
This does what I need it to do, which admittedly is a very narrow use case. Feel free to copy or fork if you'd like to extend.
