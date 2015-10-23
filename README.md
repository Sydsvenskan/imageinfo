# Image info

Small utility for generating basic image metadata from a source file. This utility detects the format of an image (jpeg/png/gif) reads its dimensions, and calculates a sha1 hash of its contents to produce a unique URI. A UUIDv5 is then created from the URI and added to the result.

## Command usage

Either use the file argument `imageinfo -file test-images/jelly.jpeg` or pipe the image data to the application `cat cat test-images/jelly.jpeg | imageinfo`, the result will be the same:

```json
{
  "uri": "cmbr://image/oT1bFPH3qvMkXTv6ewod9CoM9AA.jpeg",
  "uuid":"8481e3d7-b5ef-5e74-a18c-8551fa09ba41",
  "width":3000,
  "height":2000
}
```

## Server mode

Start imageinfo with `imageinfo -port 1080` and send images with a multipart mime request (feel free to use a boring naming scheme for the fields like: 'img1', 'img2', 'imgN'...):

`curl -X POST localhost:1080 -F "analog=@test-images/hipster.jpeg" -F "nature=@test-images/hiker.jpeg" -F "jellybelly=@test-images/jelly.jpeg"`

imageinfo then responds with the image information:

```json
{
   "nature" : {
      "width" : 3000,
      "height" : 2000,
      "uri" : "cmbr://image/OOZgjPEBdVQd7b-HQ3hIpA0yuX4.jpeg",
      "uuid" : "9a822379-489c-5979-9c1c-239844525362"
   },
   "jellybelly" : {
      "uri" : "cmbr://image/oT1bFPH3qvMkXTv6ewod9CoM9AA.jpeg",
      "uuid" : "8481e3d7-b5ef-5e74-a18c-8551fa09ba41",
      "width" : 3000,
      "height" : 2000
   },
   "analog" : {
      "uuid" : "0788077a-c719-5ec8-9da1-9954f83eac35",
      "uri" : "cmbr://image/VcHoDTqrGvd1fL8FOsX2TqoRiQM.jpeg",
      "height" : 3264,
      "width" : 4896
   }
}
```
