# lavacli

A golang implementation of the offical [lavacli](https://pypi.org/project/lavacli/)
python application.

Written against the 2019.12 LAVA spec. Older or future version might be incompatible!

## License

BSD-3-Clause

## Implemented functions

The following functions are implemented:
* identities list
* identities show
* identities add
* identities delete
* devices list
* devices show
* devices tags list
* devices tags delete
* devices tags add
* device-types list
* device-types template set
* device-types template get
* device-types health-check set
* device-types health-check get
* jobs list
* jobs logs
* jobs show
* jobs definition
* jobs validate
* jobs submit
* jobs cancel
* jobs fail
* results (testjob only)

## Building the cli

```
cd cmd/lavacli
go build .
```

## Using the API

1. Import the package
```
	import "github.com/siro20/lavacli/pkg/lava"
```

2. Connect to the server using identity "default" as defined in ~/.config/lavacli.yaml:

```
c, err := lava.LavaConnectByConfigID("default")
if err != nil {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
```

3. List jobs:

```
	ret, err := c.LavaJobsList("", "", 0, 25)
	if err != nil {
		return err
	}
	fmt.Printf("jobs:\n")
	for _, v := range ret {
		fmt.Printf("* %d %s,%s [%s] (%s) - %s\n", v.ID, v.State, v.Health, v.Submitter, v.Description, v.DeviceType)
	}
```
