# Backup

### DB백업

```bash
$ mongodump -o /dbdump/20200816
```

## DB복원

```bash
$ mongorestore --drop /dbdump/20200816
```