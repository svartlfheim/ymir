# PgAdmin

See here for info about the container, docker hub has very little:
https://www.pgadmin.org/docs/pgadmin4/development/container_deployment.html

## Exporting servers from pgadmin

Firstly once pgadmin is up and running, you can configure connections in the UI. To make life easier (so you don't have to reconfigure the connection between each restart), you can export the servers and put them in the `./servers.json` file. You can view the official docs for this here: https://www.pgadmin.org/docs/pgadmin4/development/import_export_servers.html

For this environment you can run the following commands:

```
# This will put you in the pgadmin container
$ make exec-pgadmin
$ /venv/bin/python3 setup.py --dump-servers /tmp/servers.json --user postgres@ymir.local
```

We need to use the same path to the python executable that is used when pgadmin is launched. I found this path by running `ps` on the server and seeing exactly what command was used to start the pgadmin app. Next we use `--dump-servers` option, to dump servers, and the path `/tmp/servers.json` can be anything. The `--user postgres@ymir.local` needs to be the user you were logged in as when you created the database. This email is defined by an env var: PGADMIN_DEFAULT_EMAIL which you can find the value of in the docker-compose file.

Once you've created this file, simply cat the file with `$ cat /tmp/servers.json` and paste the resulting json into the servers.json file.

The servers.json file is mounted into pgadmin as defined in the dockerfile.