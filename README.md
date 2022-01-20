# GoMota

GoMota is a mass upgrade tool for Tasmota devices. For each execution, multiple
devices on the specified network can be upgraded one version, according to the
specified - or default - upgrade path.

This is a tool built out of need, by someone still learning the language.

Latest version in the default upgrade path is `10.1.0`.

## Usage

```
$ ./bin/gomota -h
Usage of ./bin/gomota:
      --debug                 Enable debugging
      --network string        Network to scan in network/prefix notation (default "10.69.1.0/24")
      --password string       Tasmota password
      --upgrade-path string   Firmware upgrade path (default "1.0.11,3.9.22,4.2.0,5.14.0,6.7.1,7.2.0,8.5.1,9.1.0,10.1.0")
      --username string       Tasmota username
```

```
$ ./bin/gomota --network 10.69.1.0/24 --username admin --password MySecretPassword
```
