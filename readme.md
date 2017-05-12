# hksensibo

Sensibo accessory for [hc](https://github.com/brutella/hc).

## Usage

With single pod accessory

```golang
package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/llun/hksensibo"
)

func main() {
	pods := hksensibo.Lookup("user-api-key")
	if (len(pods) > 0) {
		firstPod := pods[0]
		t, err := hc.NewIPTransport(hc.Config{
			Pin:  "32191123",
		}, firstPod)
		if err != nil {
			log.Fatal(err)
		}

		hc.OnTermination(func() {
			t.Stop()
		})

		t.Start()
	}
}
```

or using with bridge

```golang
package main

import (
	"github.com/llun/hkbridge"
)

func main() {
	hkbridge.Start()
}
```

and with [hkbridge](https://github.com/llun/hkbridge) configuration

```golang
{
  "name": "Bridge",
  "manufacturer": "AwesomeMe",
  "serial": "141592653",
  "model": "WPQ864",
  "pin": "12345678",
  "accessories": [
    {
      "type": "github.com/llun/hksensibo",
      "option": {
        "key": "sensibo-key"
      }
    }
  ]
}

```

## License

MIT
