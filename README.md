# myeelight

golang command for controlling Mi LED Desk Lamp - but probably works for other Yeelight lights

# NOTE - this is very first drafty but it works

See the code for more but..


Set `-host` to the address of your desk lamp.

```
$ myeelight -host 192.168.178.1 on
```

Each command line arg performs an operation in sequence and there's a 200 millisecond pause on each:

* on, off, warm, cold, low, high, toggle
* 1--100 sets the brightness
* numbers over 1000 set the colour temperature, useful values are 2700..6500
* comma will pause for one second
* listen will connect and listen for status updates from the lamp forever

```
$ myeelight off
192.168.178.31:55443 <- {"id": 1, "method": "set_power",  "params":["off", "smooth", 3000]}
{"id":1, "result":["ok"]}
{"method":"props","params":{"power":"off"}}


$ myeelight on
192.168.178.31:55443 <- {"id": 1, "method": "set_power",  "params":["on", "smooth", 3000]}
{"id":1, "result":["ok"]}
{"method":"props","params":{"power":"on"}}


$ myeelight on warm 40
192.168.178.31:55443 <- {"id": 1, "method": "set_power",  "params":["on", "smooth", 3000]}
{"id":1, "result":["ok"]}
192.168.178.31:55443 <- {"id": 1, "method": "set_ct_abx", "params":[2700, "sudden", 0]}
{"id":1, "result":["ok"]}
{"method":"props","params":{"ct":2700}}
192.168.178.31:55443 <- {"id": 1, "method": "set_bright", "params":[40, "sudden", 0]}
{"id":1, "result":["ok"]}


Commas will give 1 second pauses

$ myeelight warm , cold , warm , cold , warm
192.168.178.31:55443 <- {"id": 1, "method": "set_ct_abx", "params":[2700, "sudden", 0]}
{"id":1, "result":["ok"]}
192.168.178.31:55443 <- {"id": 1, "method": "set_ct_abx", "params":[6500, "sudden", 0]}
{"id":1, "result":["ok"]}
{"method":"props","params":{"ct":6500}}
192.168.178.31:55443 <- {"id": 1, "method": "set_ct_abx", "params":[2700, "sudden", 0]}
{"id":1, "result":["ok"]}
{"method":"props","params":{"ct":2700}}
192.168.178.31:55443 <- {"id": 1, "method": "set_ct_abx", "params":[6500, "sudden", 0]}
{"id":1, "result":["ok"]}
{"method":"props","params":{"ct":6500}}
192.168.178.31:55443 <- {"id": 1, "method": "set_ct_abx", "params":[2700, "sudden", 0]}
{"id":1, "result":["ok"]}
{"method":"props","params":{"ct":2700}}


```

# References

* https://tasmota.github.io/docs/devices/Xiaomi-Mi-Desk-Lamp/
* https://www.mi.com/global/smartlamp
* https://github.com/hphde/yeelight-shell-scripts
* https://www.yeelight.com/en_US/developer
* https://www.yeelight.com/download/Yeelight_Inter-Operation_Spec.pdf
* https://yeelight.readthedocs.io/en/latest/
* https://gitlab.com/stavros/python-yeelight

# TODO

Need to be able to add smooth transition time to each arg, ex: `myeelight off:5000`

