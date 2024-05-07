<p align="center">
  <img src="https://steven-steven.github.io/electroninvoice/readme_assets/dpLogo.jpg" width="50" />
 </p>

Backend API that powers the DP Invoice APP. See the [native app](https://github.com/steven-steven/electroninvoice) built in Electron.

This app is created with:

- [Go Kit](http://gokit.io/)
- [jeremyschlatter's Firebase wrapper](https://godoc.org/github.com/jeremyschlatter/firebase) to connect to Firebase Realtime DB and provide mock interface for testing.
- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP router

## Starting Development

1. import `./config/config.json` file for secrets (firebase configurations)
2. Install dependencies

   ```bash
   go get
   ```
3. Run the app

   ```bash
   go run .
   ```
## Run Unit Tests
2. Run test for invoice API

   ```bash
   go test ./invoice
   ```
2. Run test for inventory item API

   ```bash
   go test ./item
   ```
