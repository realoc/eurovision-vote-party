# S2: Docker & Firestore Emulator Setup

## Status
- [x] Completed

## Dependencies
- S1 (Project Setup)

## Tasks
- [x] Create `docker-compose.yml` with Firestore emulator
- [x] Create Firestore config files (`firebase.json`, `.firebaserc`, `firestore.rules`, `firestore.indexes.json`)
- [x] Add scripts to start/stop emulator
- [x] Document local development setup

## Details

### docker-compose.yml
Use `google/cloud-sdk` or `firebase/firebase-tools` image with Firestore emulator.

### Firebase Config Files

**firebase.json**
```json
{
  "emulators": {
    "firestore": {
      "port": 8081,
      "host": "0.0.0.0"
    },
    "ui": {
      "enabled": true,
      "port": 4000
    }
  }
}
```

**firestore.rules** (permissive for dev)
```
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /{document=**} {
      allow read, write: if true;
    }
  }
}
```

### Makefile Targets
- `make run` - Start the server
- `make test` - Run tests
- `make emulator-start` - Start Firestore emulator
- `make emulator-stop` - Stop Firestore emulator

## Verification
- `docker-compose up` starts Firestore emulator
- Emulator UI accessible at localhost:4000
- Firestore accessible at localhost:8081
