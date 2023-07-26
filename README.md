# computeblade-agent
> :warning: this is WIP & just had its first alpha release. APIs/configuration/... is subject to change!

The `computeblade-agent` is an OS agent interfacing with the [ComputeBlade](http://computeblade.com) hardware.
It controls fan speed, LEDs and handles common events e.g.  to _identify_/find an individual blade in a server rack.

In addition, it exposes hardware- and agent-related metrics on a [Prometheus](http://prometheus.io) endpoint (hardcoded on port 9666 right now).


## Components

### computeblade-agent
The agent is an event-loop handler that's reacting on system events such as button presses, API calls or temperature changes (eventually).
It also exposes a prometheus endpoint.


### bladectl - interacting with the agent
The bladectl interacts with the node-local API exposed by the computeblade-agent.
You can e.g. identify the blade in a rack using `bladectl identify --wait`, which will block & make the edge-LED blink until the button is pressed.

Or change the fan-speed to 80% by invoking `bladectl fan set-percent 80`.

## Install Instructions
At this point, there are no easy-to-use install instructions but the goal is to integrate this as a default component on any computeblade.

Debian, RPM and archlinux packages are released and can be installed after downloading, e.g. using `dpkg -i <whateveristhelatestversion>.deb`

The computeblade-agent also ships with a systemd unit which can be enabled using `systemd enable computeblade-agent.service --now`.
`bladectl` is available within the PATH, but has to be executed as sudo since the socket (hardcoded on `/tmp/computeblade-agent.sock`) has no permissions configured yet.
