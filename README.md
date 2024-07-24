# Image Processor

## How to run
* After cloning, from the project root, run the command **docker-compose up --build**
* With the pods running, it`s possible then to call the API´s.

# API`s
## Generate Upload Link

URL: **localhost:8000/upload-link**

Method: **POST**

Request Body: **{
	"secret_token":"{token}",
	"expiration":100
}**

Response Body: **{
	"upload_link": "/upload?token=
	{token}"
}**

## Upload Image

URL: **localhost:8000/upload?token={token}**

Method: **POST**

Request Header: **Content-Type:multipart/form-data**

Request Body: **{
	"images":"{file}"
}**

Response Body: **{
	"image_ids": [
		4
]}**

## Get Image

URL: **localhost:8000/image/{id}**

Method: **GET**

## Get Statistics

URL: **localhost:8000/statistics**

Method: **GET**

Header: Authorization / secret-token

Response Body: **{
	"most_popular_format": ".jpeg",
	"top_camera_models": [
		""
	],
	"upload_frequency_per_day": {
		"2024-07-23T00:00:00Z": 3,
		"2024-07-24T00:00:00Z": 2
	}
}**


Next steps

- Remove the logic from the handlers and move them to a use cases layer
- Add tests to repositories and use cases layers
- Improve access token logic in the application

# Architecture

* refer to **/architecture** folder