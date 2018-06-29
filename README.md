# RKN monitor


## Config

```yaml
http-addr: ":8090"

git-repo: "https://github.com/zapret-info/z-i"
git-repo-filename: dump.csv
git-interval: 10s

debug: true

smtp:
  host: "smtp.gmail.com:465"
  user: "email@gmail.com"
  pass: "password"

subscribers:
  - email: "email@gmail.com"
    ip: [127.0.0.1]
```

If SMTP host not specify subscribers alert not sending and output it in to stdout

## API

All requests is method `GET`, response `text`

### Get blocked status

 - url path: `/`
 - query: `?ip=127.0.0.1` 

Return current status of this IP address, blocked or not

### Subscription

 - url path: `/subscribe/`
 - query: `?email=name@host.com&ip=127.0.0.1&ip=255.255.255.255`

Add email address to subscription list, if any ip address is blocked, on this email address sending alert message

### Unsubscription

 - url path: `/unsubscribe/`
 - query: `?email=name@host.com`

Remove email address of subscription list
