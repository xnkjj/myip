```bash
docker run --name myip -d -p 80:80 -e "SECRET=your_secret" xnkjj/myip
```

```bash
curl ${YOURDOMAIN}/put/$(curl https://ipinfo.io/ip -s) -H "SECRET:your_secret"
```

```bash
ssh root@$(curl ${YOURDOMAIN} -s -H "SECRET:your_secret")
```
