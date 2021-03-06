# ssh-agent
Support OpenSSH AuthorizedKeysCommand

[![PayPal donate button](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=53QAHAG6GNDMQ)

### Quick Start

required golang compile enviroment

#### agent-server
```
$ git clone https://github.com/colindev/ssh-agent && cd ssh-agent
curl -LO https://storage.googleapis.com/ssh-agent/release/$(curl -s https://storage.googleapis.com/ssh-agent/stable.txt)/linux/amd64/ssh-agent
$ sudo make install -e AGENT=[agent listen on]
```

#### agent-server - golang environment
```
$ git clone https://github.com/colindev/ssh-agent && cd ssh-agent
$ make
$ sudo make install -e AGENT=[agent listen on]
```

then set ssh key and tag
```
$ vim /etc/ssh-agent-server.conf
```

```
user-A,tag1|tag2,user-A-key
user-A,*,user-A-key2
team-user,*,user-A-key
team-user,*,user-B-key
```

tag is the client hostname whitch the machine hostname that you want to login via the agent, and you can use `*` to match part of name

#### agent-client
```
$ git clone https://github.com/colindev/ssh-agent && cd ssh-agent

// on local machines
$ sudo make install-authorization -e AGENT=[your agent host:port]

// on GCP
$ sudo make install-authorization-gcp -e AGENT=[your agent host:port]
```

### Profile

![profile](./profile.png)

### References
- [openssh auth](https://blog.heckel.xyz/2015/05/04/openssh-authorizedkeyscommand-with-fingerprint/)
