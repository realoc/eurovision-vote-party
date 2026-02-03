# Eurovision Vote Party – Server

Local development for the Go backend relies on Docker and the Firebase Firestore emulator. The provided `Makefile` wraps the most common workflows.

## Prerequisites
- Go 1.25 or newer
- Docker with the Compose plugin

## Development Tasks
- `make run` – start the Go server on http://localhost:8080
- `make test` – execute Go unit tests
- `make emulator-start` – launch the Firestore emulator (API on http://localhost:8081, UI on http://localhost:4000). The first run builds a lightweight container with Firebase CLI + JDK 21; subsequent runs reuse the built image.
- `make emulator-stop` – stop and remove the emulator container

Set the environment variable `FIREBASE_PROJECT_ID` when calling `make emulator-start` to override the default local project id (`evp-local`):

```bash
FIREBASE_PROJECT_ID=my-local-project make emulator-start
```

The Firestore emulator configuration lives at the repository root (`firebase.json`, `.firebaserc`, `firestore.rules`, `firestore.indexes.json`). Updates to rules or indexes are picked up automatically the next time the emulator starts.
