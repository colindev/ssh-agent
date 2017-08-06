#!/bin/bash

# only for GCP
AUTH_SCRIPT=/usr/local/bin/ssh-authorization.sh

project=`curl "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google" --connect-timeout 3 2>/dev/null`
name=`curl "http://metadata.google.internal/computeMetadata/v1/instance/name" -H "Metadata-Flavor: Google" --connect-timeout 3 2>/dev/null`
cat <<EOF > $AUTH_SCRIPT 
#!/bin/bash

curl http://{{selfLink}}/users/\${1}/keys?fingerprint=${name}.${project}

EOF

chmod u+x $AUTH_SCRIPT
chown root:root $AUTH_SCRIPT

cp /etc/ssh/sshd_config /etc/ssh/sshd_config.bak
sed -i "s@#\?AuthorizedKeysCommand\s\+[^#]\+@AuthorizedKeysCommand ${AUTH_SCRIPT}@" /etc/ssh/sshd_config
sed -i 's/#\?AuthorizedKeysCommandUser\s\+[^#]\+/AuthorizedKeysCommandUser root/' /etc/ssh/sshd_config
systemctl reload sshd

for user in `curl http://{{selfLink}}/users?project=$project 2>/dev/null`
do
adduser $user
done

