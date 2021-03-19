# Mashiron-GO!

Module-style next gen bot for multiple chat platforms

## Install

### Local

#### Quick start

1. Clone this repo
2. Install make and go
3. Run make
4. `cd bin`
5. Put bot token to mashiron.ini
6. Run binary (Ex: `./discord`)
7. Profit

#### Make sh module work

1. install systemd-nspawn,shellcheck(optional)
2. add systemd-nspawn to sudoers NOPASSWD
3. Profit

#### Deamonize

1. Edit mashiron.service for your env
2. move mashiron.service to /etc/systemd/system
3. `systemctl start mashiron`and`systemctl enable mashiron`(optional)
4. Profit

## Docker (ArchLinux host only)

1. Download mashiron.ini and docker.sh to the host data dir
2. `sudo bash docker.sh`
3. Profit

## Bugs

You tell me

//WIP

