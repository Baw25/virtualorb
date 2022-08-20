# virtualorb

Simulator for orb signups and status report (battery, cpu usage, cpu temp, disk space) requests for the purpose of testing. 
The orb can: 
- Simulate periodic report status requests reporting cpu usage, cpu temp, disk space of current machine 
- Simulate signups and submit png images to a mock API with generated associated ids

The idea is that can be used as part of a CI pipeline or deployed as a workload.

## Setup

## Example
```
   python simulator.py 
```

### Example Output:
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