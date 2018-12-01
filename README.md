# rescreen [![Go Report Card](https://goreportcard.com/badge/github.com/vcraescu/rescreen)](https://goreportcard.com/report/github.com/vcraescu/rescreen) [![Build Status](https://travis-ci.com/vcraescu/rescreen.svg?branch=master)](https://travis-ci.com/vcraescu/rescreen) [![Coverage Status](https://coveralls.io/repos/github/vcraescu/rescreen/badge.svg?branch=master)](https://coveralls.io/github/vcraescu/rescreen?branch=master)

The project started because of GNOME lack of support for multiple monitors with different pixel densities. At this moment, 
in GNOME you can't have multiple monitors with different pixel densities. On one screen you will see everything too small 
or too large. Therefore you have to set the scaling factor in GNOME until everything is too large and then adjust each 
monitor with `xrandr`. `rescreen` makes that easier allowing you to define a config file. 

## Overview

#### Config file

**config.json**
```
{
  "layout": [
    "HDMI-0", "eDP-1-1", "", "",
    "", "", "", ""
  ],
  "monitors": {
    "HDMI-0": {
      "scale": 1.333, // everything is too large so we zoom out by 33%
      "primary": true // set this monitor as primary, here's where you will see the top bar
    }
  }
}
```

**HDMI-0**, **eDP-1-1** - are the output names you get when you run `xrandr`:

e.g:

```
$ xrandr
Screen 0: minimum 8 x 8, current 8959 x 2880, maximum 32767 x 32767                                                                                                                                                       
HDMI-0 connected primary 5119x2879+0+1 (normal left inverted right x axis y axis) 597mm x 336mm panning 5119x2880+0+0 tracking 8959x2880+0+0 border 0/0/0/0                                                               
   3840x2160     30.00*+  29.97    25.00    23.98                                                                                                                                                                         
   1920x1200     59.88                                                                                                                                                                                                    
   1920x1080     60.00    59.94    50.00    23.98    60.05    60.00    50.04                                                                                                                                              
   1680x1050     59.95                                                                                                                                                                                                    
   ...
eDP-1-1 connected 3840x2160+5119+0 (normal left inverted right x axis y axis) 346mm x 194mm
   3840x2160     60.00*+  59.98    48.02    59.97
   3200x1800     59.96    59.94
   2880x1620     59.96    59.97
   2560x1600     59.99    59.97
   ...
```

Layout is actually a matrix with 4 columns per row even though it is given as an array. It defines your arrangement 
when you have multiple monitors.

#### Run

`rescreen /path/to/config.json`

If no config file path is given, it will look up for a config.json inside current folder.

#### Dry run

`rescreen --dry-run /path/to/config.json`

Just display the commands which will run under the hood.

### Supported platforms

It works only on Linux and it depends on `xrandr`.
