FROM archlinux:latest
WORKDIR /opt/Mashiron-go/bin
COPY bin .
COPY mashiron.service /etc/systemd/system
RUN pacman -Syy --noconfirm sudo && systemctl enable mashiron && systemctl mask systemd-firstboot
ENTRYPOINT /sbin/init
