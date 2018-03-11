# ldp2lpd

Track your [OVH Log Data Platform streams](https://www.ovh.com/fr/data-platforms/logs/) on [Novation Lauchpad Mini](https://global.novationmusic.com/launch/launchpad-mini)

This tool is currently only working with Launchpad S (Green-Red Launchpads).

## Installation

```
$ go get github.com/fsamin/ldp2lpd
```

Portmidi is required to use this package.

```
$ apt-get install libportmidi-dev
# or
$ brew install portmidi
```

### Usage
```
Usage of ./ldp2lpd (Version dev):
      --address string     URI of the websocket
      --config string      Configuration file
      --time-factor int    time factor of interval (default 30)
      --time-unit string   time unit of interval (default "second")
      --verbose            verbose mode
```

### Example
```
$ ldp2lpd --address "wss://your.logs.ovh.com>/tail/?tk=h12366f3-25aa-754b-abcd-zrd678e1345678" --time-factor 5 --time-unit minute
```