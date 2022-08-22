# virtualorb

Simulator for orb signups and status report (battery, cpu usage, cpu temp, disk space) requests for the purpose of testing. 
The orb can: 
- Simulate periodic report status requests reporting cpu usage, cpu temp, disk space of current machine 
- Simulate signups and submit png images to a mock API with generated associated ids

The idea is that can be used as part of a CI pipeline or deployed as a workload.

## Setup

### Dependencies
Make sure to have the following:
- python >= 3.6.0
- go >= 1.16
- gin web framework = 1.8.1
- [`osx-cpu-temp project`](https://github.com/lavoiesl/osx-cpu-temp)
- lm-sensors (if you have linux machine)

[`Install go (if needed)`](https://go.dev/doc/install)
[`Install python 3 (if needed)`](https://www.python.org/downloads/)
[`golang gin`](https://github.com/gin-gonic/gin#installation)

### Getting Started
Install requirements:
```
pip install -r requirements.txt
```


Build project and osx-cpu-temp subproject
```
make build
```
Run tests:
```
make test
```
Run the virtualorb server
```
make run
```
Open additional terminal and run simulator.py to simulate report and signup submissions

## Examples
Run 1 test signup and 1 test report
```
   python simulator.py 
```
Or run any combination of report or signup submissions over intervals of time. 

Run 5 test signups with 2 second intervals in between
```
   python simulator.py iterations=5 interval=2
```
Run 10 test report submissions with 2 second intervals in between
```
   python simulator.py iterations=10 interval=2 report
```

### Example Output for 1 signup and 1 report:
```
Running simulator
Sending signup request...
Signup generated:
{'action_id': '40df4e77-4b3f-4d11-aecb-cfeb5c335429',
 'message': 'Signup successfully submitted!'}
Sending status report request...
Report sent:
{'battery_level_percent': '75.000000',
 'battery_voltage': '11.200000V',
 'device_cpu_percent': '31.663366',
 'device_cpu_temp': '43°C',
 'device_disk_space_available_percent': '92.929450',
 'device_disk_space_used_percent': '7.070550'}
```

## Improvements and next steps