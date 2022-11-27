# go-object-store

A basic object store for self hosted apps.

## Build and Run

In the base directory of this repository run this command to build the service
```
go build
```

Then run the following command to run the service. This will create a new `./objectData`
directory where objects will be stored. And the service will be running at `http://localhost:3000`.
```
./go_object_store
```

## Endpoints

```
GET /alive - check if service is up

PUT /object/<object-key> - create an object where the body is the file data

GET /object/<object-key> - fetch an object's data

DELETE /object/<object-key> - remove an object
```

