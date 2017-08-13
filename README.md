# Hemtjänst

Hemtjänst is:
  * A Swedish word for home care or home service
  * A specification of sorts on how devices should register with an MQTT broker
  * An MQTT to HomeKit bridge

This project was started so that a collection of $platform-to-MQTT bridges
could be created or adapted that announce the devices they have in a similar
enough fashion that we can automatically generate HomeKit accessories from it.
Though the aim for now is HomeKit the metadata published should be enough to
also be able to create bridges for other platforms, like SmartThings.

## Installation

The binary is fully self-contained and can be installed by simply issuing a:

```
go install github.com/hemtjanst/hemtjanst/cmd/hemtjanst
```

## Usage

Once you've `go install`ed the project a binary will be in your `$GOPATH/bin`.
If you do not have explicitly set a `$GOPATH` environment variable it defaults
to `$HOME/go`. Ensure to add `$GOPATH/bin` to your `$PATH` or use the full path
to the binary.

By default it will connect to an MQTT broker on `localhost:1883` and expose a
HomeKit bridge on port `12345` with pairing pin-code `01020304`.

Pass a `--help` for all available options.

## Web UI

There's currently a very experimental web UI included. However, the templates
aren't bundled when building the binary just yet so you'll need to put those
in a `web/templates` directory relative to the path of the binary for it to
work.

Once the UI starts to do something useful aside from just listing devices and
their features the templates will be embedded in the binary to avoid this
problem entirely.

## Specification

### Discovery

In order to not have to continuously scan for new devices, or to subscribe
to every event, a discovery mechanism has to be implemented.

A `discover` topic with retain set must exist. Any device that joins the
network must subscribe to `discover`. Since the retain bit is set they will
receive the discovery request and must now announce themselves to the rest
of the network. If someone then wants to initiate a full discovery all they
need to do is publish again to discover (with the retain bit set).

An announce is done by publishing to the `announce` topic. The body of the
message must be the "root" topic of entity that joined, for example
`light/ground_floor/kitchen/stove`, for the kitchen light near the stove.
The topic layout is entirely arbitrary and so is the naming of the device.
Anyone interested in knowing about devices now simply subscribes to `announce`
and gets to know any device that joins.

**Please note**, the name of the device (the last bit of the "root" topic) is
used to generate a unique identifier for this device, so if you change it it'll
be like you removed the existing device and added a new one.

Whenever a device leaves it publishes its "root" topic to the `leave` topic.
Similarly to announce this allows other clients to cancel their subscriptions
if they were explicitly wathcing that device or take any other ation.

#### Last Will and Testament

MQTT includes a Last Will And Testament feature that will cause the broker to
publish a message on a specific topic when the client is gone, gracefully or
not.

If every device sets up its own MQTT client it can hence specify that as a
last will and testament the broker should publish to `leave` with its original
topic. This ensures that we can always properly clean up, even if the device
falls of a cliff.

However, for bridged clients, lets say IKEA Trådfri lights that are published
through a trådfri-to-mqtt broker, the bridge is the only MQTT client. Since we
can't specify multiple actions in a last will and testament the bridge can
still tell the broker to publish to a message to the topic, the content being
a string representing a "leave ID". Every device that the bridge used to bridge
must specify that same "leave ID" in its metadata so that a mapping can be
maintained between "leave ID"s and devices for any cleanup purposes, such as no
longer announcing this accessory to HomeKit.

## Metadata

In order to know what type of device we're dealing with every device is
expected to have a `meta` topic underneath its "root" topic. This topic can be
written to by any existing bridge or by any entity watching `announce` that
knows how to map this device to the metadata specification.

This `meta` topic is a JSON object serialised as a string that contains any
additional information needed for Hemtjänst to do its thing and generate
HomeKit accessories.

As it currently stands the `meta` topic types and features follow the exact
names in the HomeKit specification. To support another platform like
SmartThings a similar mapping would need to be created.

The `meta` document contains a number of required and optional entries. The
required ones are: `name`, `type`, `feature`. The rest is optional.

Optional keys are: `lastWillID`.

The naming of the keys follows [Google's JSON style guide][json-style] and as
such are in *camelCase*. However, `ID` is always fully uppercase and any
chemical formula expressed in chemical symbols follows its relevant casing, so
`CO2`, not `Co2` for carbon dioxide.

### `name`

A human friendly device name. Can be the same thing as the topic name or
something else entirely. This will be the name of the accessory as HomeKit sees
it so do pick something that makes sense and allows you to relatively easily
identify the accessory.

### `type`

The type of device, for example `light` or `CO2Sensor`. These map directly onto
HomeKit services and are considered the "primary" service. There is currently no
support for hidden, secondary or linked services.

You can find the supported devices [here][types] and how they map to HomeKit
services.

### `feature`

Every device needs at least one feature to be defined for it, the usual
required characteristic (it can be more than one). Similarly to device the
[characteristics][characteristics] mapping contains a list of which feature
names map to what HomeKit characteristics.

A feature is an object itself, which can be empty, in which case the defaults
apply for the `min`, `max` and `step` value as defined in the HomeKit spec. In
order to override them you can specify `min`, `max` and `step` keys and set the
appropriate value. Do note that you cannot change the type of the value, so if
the HomeKit specifies something as a `uint8` it will be deserialised as a
`uint8`.

Similarly, we expect that in order to get and set the value of a feature a
`"root topic"/<feature>/get` and `set` topic exist that we can use. If those
topics are named differently you have to specify a `getTopic` and a `setTopic`
key that have the full path to a topic (so not necessarily nested under the
"root" topic) that should be used instead.

### `lastWillID`

The `lastWillID` only has to be set for bridged devices, so in cases where each
device doesn't itself maintain a connection to the MQTT broker. When that is
the case the Last Will And Testament should be set instead to publish the
"root" topic name to the `leave` topic.

The `lastWillID` can be anything but needs to be unique. As such it's recommended
to use a UUIDv4 for this.

### Examples

The `meta` topic for a light that can just be turned on and off looks like
this:

```json
{
  "name": "kitchen stove light",
  "type": "light",
  "feature": {
    "on": {}
  }
}
```

Most smart lights however can also dim and usually set the colour temperature
too, so it's more likely you'll want to publish something like this:

```json
{
  "name": "living room light",
  "type": "light",
  "feature": {
    "on": {},
    "brightness": {},
    "colorTemperature": {
      "getTopic": "light/ground_floor/living_room/warmth/get",
      "setTopic": "light/ground_floor/living_room/warmth/set",
    },
  }
}
```

A contact sensor, something you might put on a window to detect if it is open
or closed (as part of a security system for example) can be defined like this:

```json
{
  "name": "bathroom window",
  "type": "contactSensor",
  "feature": {
   "contactSensorState": {}
  },
  "lastWillID": "f56ad37c-aa0f-45f4-8e92-f9a6dba39d84"
}
```

[json-style]: https://google.github.io/styleguide/jsoncstyleguide.xml
[types]: homekit/util/service.go
[characteristics]: homekit/util/characteristic.go
