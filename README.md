# computeblade-agent
> :warning: this is still a beta-release & configuration&APIs might see breaking changes! It's not 100% feature complete yet but works

The `computeblade-agent` is an OS agent interfacing with the [ComputeBlade](http://computeblade.com) hardware.
It controls fan speed, LEDs and handles common events e.g.  to _identify_/find an individual blade in a server rack.
In addition, it exposes hardware- and agent-related metrics on a [Prometheus](http://prometheus.io) endpoint.

**TL;DR, I just want it running on my blade script**:
```bash
curl -L -o /tmp/computeblade-agent-installer.sh https://raw.githubusercontent.com/Uptime-Lab/computeblade-agent/main/hack/autoinstall.sh
chmod +x /tmp/computeblade-agent-installer.sh
/tmp/computeblade-agent-installer.sh
```

## Components

### computeblade-agent
The agent is an event-loop handler that's reacting on system events such as button presses and temperature changes.
It also exposes a prometheus endpoint allowing monitoring of core-metrics such as PoE status.

By default, the computeblade agent runs in _normal_ operation mode; the LEDs are static and fanspeed is set based on the configuration.
In case the SoC temperature raises above a predefined level, the _critical_ mode is active and sets the fan-speed to 100% alongside changing the LED color (Red by default)

Aside from the above mentioned normal and critical modes, the _identify_ action (independend of the mode), which lets the edge LED blink.
This can be toggled using `bladectl` on the blade (`bladectl identify`) or by pressing the edge button.


### bladectl - interacting with the agent
The bladectl interacts with the blade-local API exposed by the computeblade-agent.
You can e.g. identify the blade in a rack using `bladectl identify --wait`, which will block & make the edge-LED blink until the button is pressed.


## Install Options

The agent and bladectl are provided as package for Debian, RPM and ArchLinux or as OCI image to run within docker/Kubernetes.
Packages ship with a systemd unit which can be enabled using `systemd enable computeblade-agent.service --now`.

`bladectl` is available globally, but has to be executed as root since the socket (default `/tmp/computeblade-agent.sock`) does not have a user/group accessed due to privileged access on critical resources.

**Kubernetes deployment**:
A kustomize environment can be found in `hack/deploy`. A `kubectl -k hack/deploy` does the trick - or use a GitOps tool such as FluxCD.


## Configuration
The configuration is driven by a config file or environment variables. Linux packages ship with the default configuration placed in `/etc/computeblade-agent/config.yaml`.
Alternatively (specifically for running within Kubernetes), all parameters in the YAML configuration can be overwritten using environment variables, prefixed with `BLADE_`:

Changing the metric address defined in YAML like this:
```yaml
# Listen configuration
listen:
  metrics: ":9666"
```
is driven by the environment variable `BLADE_LISTEN_METRICS=":1234"`.

Some useful parameters:
- `BLADE_STEALTH_MODE=false` Enables/disables stealth mode
- `BLADE_FAN_SPEED_PERCENT=80` Sets static fan-speed (by default, there's a linear fan-curve of 40-80%
- `BLADE_CRITICAL_TEMPERATURE_THRESHOLD=60` Configures critical temperature threshold of the agent
- `BLADE_HAL_BCM2711_DISABLE_FANSPEED_MEASUREMENT=false` enables/disables fan speed measnurement (disabling it reduces CPU load of the agent)

