# KON-KHMER-GENZ

## Features (Safe Overview)

* Supports different testing modes:

  * **TCP Connection Load Test** (Mode 1)
  * **UDP High-Throughput Test** (Mode 2)
  * **HTTP Load Simulation** (Mode 3)
* Customizable parameters such as connections, duration, timeout, and packet size.
* Optional proxy file usage for authorized HTTP load simulations.

## Usage (Safe & Official)

This tool accepts the following arguments:

```
<host> <port> <mode> <connections> <seconds> <timeout(second)> <packetsize(bytes)> [proxyfile.txt]
```

If the argument count is incorrect, the tool displays:

```
Usage: <program> <host> <port> <mode> <connections> <seconds> <timeout(second)> <packetsize(bytes)> [proxyfile.txt]
---

