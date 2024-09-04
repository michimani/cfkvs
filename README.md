cfkvs
===

This is a simple command line tool for CloudFront Key Value Store.

[![codecov](https://codecov.io/gh/michimani/cfkvs/graph/badge.svg?token=PWKPWONA8G)](https://codecov.io/gh/michimani/cfkvs)

## Features

- KeyValueStore
  - list
  - create
- Item (Key-Value pair)
  - list
  - get
  - put
  - delete
  - sync

### Comparison with AWS CLI commands

| AWS CLI | cfkvs |
| --- | --- |
| `cloudfront create-key-value-store` | `cfkvs kvs create` |
| `cloudfront delete-key-value-store` | - |
| `cloudfront describe-key-value-store` | - |
| `cloudfront list-key-value-stores` | `cfkvs kvs list` |
| `cloudfront update-key-value-store` | - |
| `cloudfront-keyvaluestore delete-key` | `cfkvs item delete` |
| `cloudfront-keyvaluestore describe-key-value-store` | - |
| `cloudfront-keyvaluestore get-key` | `cfkvs item get` |
| `cloudfront-keyvaluestore list-keys` | `cfkvs item list` |
| `cloudfront-keyvaluestore put-key` | `cfkvs item put` |
| `cloudfront-keyvaluestore update-keys` | - |
| - | `cfkvs item sync` |
| - | `cfkvs kvs info` |

## Installation

### Homebrew

```bash
brew install michimani/tap/cfkvs
```

## Usage

```
$ cfkvs -h

Usage: cfkvs <command> [flags]

A simple cli tool to manage CloudFront Key Value Stores.

Flags:
  -h, --help              Show context-sensitive help.
  -D, --debug             Enable debug mode.
      --output="table"    Output format. One of: json, table.
      --version           Print version information and quit

Commands:
  kvs list       List key value stores in your account.
  kvs create     Create a key value store.
  item list      List items in the key value store.
  item get       Get an item in the key value store.
  item put       Put an item in the key value store.
  item delete    Delete an item in the key value store.
  item sync      Sync items in the key value store with S3 object.
```

Run `cfkvs <command> --help` for more information on a command.

## Examples

### Sync items in the key value store with S3 object

Assume that the following key-value pairs exist in the KeyValueStore `cf-kvs-sample`.

```
+-------+---------+
| KEY   | VALUE   |
+-------+---------+
| key-1 | value-1 |
| key-2 | value-2 |
| key-3 | value-3 |
+-------+---------+
```

Consider synchronizing this KeyValueStore with the data of the following S3 object named `data.json`.

```json
{
  "data": [
    {"key": "key-1", "value": "v 1"},
    {"key": "key-2", "value": "value-2"},
    {"key": "key-4", "value": "v 4"}
  ]
}
```

You can check the changes in the KeyValueStore by synchronizing it with `data.json` using the following command.

```bash
$ cfkvs item sync \
--kvs-name='cf-kvs-sample' \
--bucket="${YOUR_BUCKET_NAME}" \
--key='data.json'
```

```
[ADDED] Following items will be added.
+---+-------+-------+
| # | KEY   | VALUE |
+---+-------+-------+
| 1 | key-4 | v 4   |
+---+-------+-------+

[UPDATED] Following items will be updated.
+---+-------+--------------+-------------+
| # | KEY   | BEFORE VALUE | AFTER VALUE |
+---+-------+--------------+-------------+
| 1 | key-1 | value-1      | v 1         |
+---+-------+--------------+-------------+

[DELETED] No items will be deleted.
```

If you want to delete keys that do not exist in `data.json`, add the `--delete` flag.

```bash
$ cfkvs item sync \
--kvs-name='cf-kvs-sample' \
--bucket="${YOUR_BUCKET_NAME}" \
--key='data.json' \
--delete
```

```
[ADDED] Following items will be added.
+---+-------+-------+
| # | KEY   | VALUE |
+---+-------+-------+
| 1 | key-4 | v 4   |
+---+-------+-------+

[UPDATED] Following items will be updated.
+---+-------+--------------+-------------+
| # | KEY   | BEFORE VALUE | AFTER VALUE |
+---+-------+--------------+-------------+
| 1 | key-1 | value-1      | v 1         |
+---+-------+--------------+-------------+

[DELETED] Following items will be deleted.
+---+-------+---------+
| # | KEY   | VALUE   |
+---+-------+---------+
| 1 | key-3 | value-3 |
+---+-------+---------+
```

If you add the `--yes` or `-y` flag, the synchronization will actually be executed.

```bash
$ cfkvs item sync \
--kvs-name='cf-kvs-sample' \
--bucket="${YOUR_BUCKET_NAME}" \
--key='data.json' \
--delete \
--yes
```

### Describe a key value store

The Describe action for CloudFront Key Value Store has two actions: **CloudFront:DescribeKeyValueStore** and **CloudFrontKeyValueStore:DescribeKeyValueStore**. The `cfkvs kvs info` command can get the merged information of these actions.

```bash
$ cfkvs kvs info --kvs-name='cf-kvs-sample' --output=json
```

Output:

```json
{
    "id": "xxxxxxxx-0000-0000-0000-xxxxxxxxxxxx",
    "arn": "arn:aws:cloudfront::000000000000:key-value-store/xxxxxxxx-0000-0000-0000-xxxxxxxxxxxx",
    "name": "cf-kvs-sample",
    "comment": "sample kvs",
    "status": "READY",
    "itemCount": 2,
    "totalSizeInBytes": 8,
    "created": "2024-08-30T14:12:09.334Z",
    "lastModified": "2024-08-30T14:12:09.334Z",
    "failureReason": "",
    "eTag": "E3XXXXXXXXXXXX"
}
```


## License

[MIT](https://github.com/michimani/cfkvs/blob/main/LICENSE)

## Author

[michimani](https://github.com/michimani)

