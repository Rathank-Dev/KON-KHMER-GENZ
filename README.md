# KON-KHMER-GENZ

## Features (Safe Overview)

* Supports different testing modes:

  * **TCP Connection Load Test** (Mode 1)
  * **UDP High-Throughput Test** (Mode 2)
  * **HTTP Load Simulation** (Mode 3)
* Customizable parameters such as connections, duration, timeout, and packet size.
* Optional proxy file usage for authorized HTTP load simulations.

### Preview
![image](1.png)
![image](2.png)
<!-- Embed MP4 video -->
<video width="640" height="360" controls>
  <source src="https://github.com/Rathank-Dev/KON-KHMER-GENZ/blob/main/2025-11-15%2018-28-46.mp4?raw=true" type="video/mp4">
  Your browser does not support the video tag.
</video>

## Usage

This tool accepts the following arguments:

```
<host> <port> <mode> <connections> <seconds> <timeout(second)> <packetsize(bytes)> [proxyfile.txt]
```

If the argument count is incorrect, the tool displays:

```
Usage: <program> <host> <port> <mode> <connections> <seconds> <timeout(second)> <packetsize(bytes)> [proxyfile.txt]

```
Usagec Example : .\kmeng-GEnZ.exe example.com 443 3 1000 60 5 65000 proxy.txt
---

Use it at your own risk.

