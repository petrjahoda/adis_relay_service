[![developed_using](https://img.shields.io/badge/developed%20using-Jetbrains%20Goland-lightgrey)](https://www.jetbrains.com/go/)
<br/>
![GitHub](https://img.shields.io/github/license/petrjahoda/adis_relay_service)
[![GitHub last commit](https://img.shields.io/github/last-commit/petrjahoda/adis_relay_service)](https://github.com/petrjahoda/adis_relay_service/commits/master)
[![GitHub issues](https://img.shields.io/github/issues/petrjahoda/adis_relay_service)](https://github.com/petrjahoda/adis_relay_service/issues)
<br/>
![GitHub language count](https://img.shields.io/github/languages/count/petrjahoda/adis_relay_service)
![GitHub top language](https://img.shields.io/github/languages/top/petrjahoda/adis_relay_service)
![GitHub repo size](https://img.shields.io/github/repo-size/petrjahoda/adis_relay_service)
<br/>
[![Docker Pulls](https://img.shields.io/docker/pulls/petrjahoda/adis_relay_service)](https://hub.docker.com/r/petrjahoda/adis_relay_service)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/petrjahoda/adis_relay_service?sort=date)](https://hub.docker.com/r/petrjahoda/adis_relay_service/tags)
<br/>
[![developed_using](https://img.shields.io/badge/database-MariaDB-red)](https://www.mariadb.org) [![developed_using](https://img.shields.io/badge/runtime-Docker-red)](https://www.docker.com)

# Adis Relay Service

## Description
Go service that enables Zapsi Relay1, if there exists an open terminal_input_order and Relay1 is not enabled

## Installation Information
Windows: install using sc.exe
Linux: install using systemd
Docker: install using image from https://hub.docker.com/r/petrjahoda/adis_relay_service

## Developer Information
For every workplace one go routine is running in a 10-second loop.
1. At the beginning of a loop:
    - program checks for assigned terminal deviceid
    - program checks for assigned zapsi deviceid
    - program checks for any open terminal_input_order (DTE == null) for assigned terminal device_id, two times with 3 second pause
2. If any open terminal_input_order found, program checks assigned zapsi for open Relay
    - If Relay1 is closed, program opens Relay1
<br><br><br><br>

#### Example of opened relay
```
[Inputs 1-8_]
0-0-0-0-0-0-0-0

[Output 1-8]
1-0-0-0-0-0-1-0
```

#### Example of closed relay
```
[Inputs 1-8_]
0-0-0-0-0-0-0-0

[Output 1-8]
0-0-0-0-0-0-1-0
```

    
© 2021 Petr Jahoda