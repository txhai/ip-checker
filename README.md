# ip-checker
A simple utility to check if any IP belongs to a given list of CIDR notation range.

---

Environment variables
```
HOST=127.0.0.1
PORT=8000
CIDR_LIST_FILE=/path/to/example-data/list.txt
```

---
Query
```shell
curl -o - http://127.0.0.1:8000/check?ip=172.20.0.10
```