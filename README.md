# Playwright Go Invoker Sample

## Set Up

[Install Playwright](https://playwright.dev/docs/intro) in a specific directory:

```sh
cd /Users/code/development/
npm init playwright@latest
# Continue the instructions from the official docs, etc.
```

And run the included sample tests from this directory:

```
export PLAYWRIGHT_DIR=/Users/code/development/playwright/
cd src/playwright-go-invoker/
go run .
```

## Configuration

The global configuration resides in `playwright.config.js` in
the directory where Playwright was installed. Observe many (all?) configuration
can be overriden via command line parameters.

By defult, the test files will be picked from the directory where Playwright was
installed, e.g. `/Users/code/development/playwright/tests`, but this can be overriden
by telling Playwright where the file tests reside. In the case of this Invoker,
simply specify a file or a directory via TEST_DIR, e.g. `export TEST_DIR=/Users/code/test-microservice.ts`.

## Details

We invoke `npx` directly, instruct it to use JSON reporting, which we directly capture,
parse and interpret. See an [alternative](#previous-approach) for a alternative. A file to unmarshal JSON
is generated from the observed output, as the JSON Schema is not explicit
(see https://github.com/microsoft/playwright/blob/main/packages/playwright/types/testReporter.d.ts
for further digging). Alternatively, we could include a custom javascript reporter and define
our own output.

We are invoking `playwright` from scratch everytime we call it, which could be alleviated if
`playwright` could be used as a long running service instead (most of the load time seems to come from starting
the headless chromium/firefox/webkit instances).

Finally, similarly to `playwright-go`, we could decide to add sugar so users do not have
to manually install `playwright` themselves first.

### Previous approach

It was considered using [playwright-go](https://github.com/playwright-community/playwright-go) directly, i.e.

```go
	opts := &playwright.RunOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	driver, err := playwright.NewDriver(opts)
	if err != nil {
		logger.Println(fmt.Sprintf("Could not create playwright driver: %w", err))
		return
	}
	err = driver.Install()
	if err != nil {
		logger.Println(fmt.Sprintf("Could not install playwright: %v", err))
		return
	}
	c := driver.Command("test", "/Users/code/development/testbed/playwright-go-invoker/src/playwright-go-invoker/tests/")

	bytes, err := c.Output()
	if err != nil {
		logger.Println(fmt.Sprintf("Could not fetch output: %v", err))
		return
	}

	handleJsonInput(bytes)
}
```

However, this cannot work as the driver used by the playwright-go is a stripped
down version of playwright itself, and hence we cannot run actual javascript or typescript
tests againt it.
