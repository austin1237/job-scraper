# Local

This is where scraping scripts that can only be ran locally and not in a remote
enviorment go.

## Deno

These scripts are written to be used with the [Deno](https://deno.land/) runtime using the following command

```bash
deno task run
```

## Set up (OSX)

For these scripts to run sucessfully it requires an already running instance of
chrome in debug mode, as well a tab logged into the target site to bypass
auth/recaptcha.

Run the following command with no other chrome process running to get Chrome
started in debug mode:

```bash
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --remote-debugging-port=9222
```
