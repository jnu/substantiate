Substantiate
===

Tool to materialize environment variables as files on disk in a specified directory.

If the vars appear to be base64-encoded, `substantiate` will automatically decode them.

## Usage

```
substantiate -vars FOO,BAR -directory /secrets
```

This will cause the env vars `FOO` and `BAR` to be written to `/secrets/FOO` and `/secrets/BAR`.
If either is base64-decodable, the decoded content will be written to the file.
If the `/secrets` directory does not exist, it will be created.

As an alternate to the `-vars` flag,
the environment variable `SUBSTANTIATE` can be passed in the same format.

## Motivation

AWS ECS does not allow mounting secrets directly as files.
This tool allows secrets to be passed into a container as environment variables and then materialized as files for compatibility with processes that require files (configs, certificates, etc).