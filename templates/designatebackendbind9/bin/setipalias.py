#!/usr/bin/python3
import os
import sys
import ipaddress
import netifaces
from pyroute2 import IPRoute

podname = os.environ.get("POD_NAME").strip()
interface_name = os.environ.get("NAD_NAME", "designate").strip()
if not len(podname):
    print(f"{sys.argv[0]} requires POD_NAME in environment variables")
    sys.exit(1)

namepieces = podname.split('-')
pod_index = namepieces[-1]

print(f"working with {interface_name} {pod_index}")

nodefile = f"bind_address_{pod_index}"
print(f"working with address file {nodefile}")
filename = os.path.join('/var/lib/predictableips', nodefile)
if not os.path.exists(filename):
    print(f"Required alias address file {filename} does not exist")
    sys.exit(1)

ip = IPRoute()
designateinterface = ip.link_lookup(ifname=interface_name)

if not len(designateinterface):
    print(f"{interface_name} attachment not present")
    sys.exit(1)


ipfile = open(filename, "r")
ipaddr = ipfile.read()
ipfile.close()
print(f"Setting {ipaddr} on {interface_name}")
if ipaddr:
    # Get our current addresses so we can avoid trying to set the
    # same address again.
    version = ipaddress.ip_address(ipaddr).version
    ifaceinfo = netifaces.ifaddresses(interface_name)[
        netifaces.AF_INET if version == 4 else netifaces.AF_INET6]
    current_addresses = [x['addr'] for x in ifaceinfo]
    if ipaddr not in current_addresses:
        mask_value = 32
        if version == 6:
            mask_value = 128
        ip.addr('add', index = designateinterface[0], address=ipaddr, mask=mask_value)
ip.close()
