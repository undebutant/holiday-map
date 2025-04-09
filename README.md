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

## Additional config

You will have to configure these fields to run this application from your server:

```go
// main.go
const SERVER_FOLDER = "/path/to/the/folder/with/app/and/data"
const DOMAIN_NAME = "my.awesome.website"
```

```html
<!-- resources/main.html customize FontAwesome kit -->

<!-- FontAwesome -->
<script src="https://kit.fontawesome.com/4988160acf.js" crossorigin="anonymous"></script>
```

```javascript
// resources/main.html
const API_PATH = 'https://my.awesome.website'
```
