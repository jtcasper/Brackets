#!/bin/bash

mkdir -p /usr/share/brackets
groupadd brackets
useradd -g brackets -d /usr/share/brackets -s $(which nologin) brackets
chown brackets:brackets /usr/share/brackets
chmod 700 /usr/share/brackets
sqlite3 /usr/share/brackets/brackets.sqlite "$(cat ../backend/migrations/*)"
chown brackets:brackets /usr/share/brackets/brackets.sqlite
cp brackets.service /etc/systemd/system/
