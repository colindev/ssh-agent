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

if [ -f /etc/selinux/config ]
then
    sed -i 's/SELINUX=enforcing/SELINUX=permissive/' /etc/selinux/config
    setenforce 0
fi

for user in `curl http://{{selfLink}}/users?project=$project 2>/dev/null`
do
    echo add user [$user]
    if uname -a | grep debian
    then
        adduser --disabled-password --gecos '' %s
    else
        adduser %s
    fi
done

