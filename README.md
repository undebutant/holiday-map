# holiday-map

A basic holiday map website to pin places you went and add photos to it, with a Go backend and a vanilla JS frontend
leveraging [Leaflet](https://leafletjs.com) for the map part.

## Required files

On the backend side two files need to be created in the `/data` folder for the website to work.

An `auth.json` file, used to limit the access to the website with Basic Auth:
```json
{
    "users": [
        {
            "username": "veryDoge",
            "password": "muchSecretManySafe"
        }
    ]
}
```

A `data.json` file, used as a makeshift database for an easy setup:
```json
{
    "markers": [],
    "photoCount": 0
}
```
