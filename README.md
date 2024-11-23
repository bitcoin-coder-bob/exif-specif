# EXIF-SPECIF

Allows you to edit the keywords tag in the exif data of jpg file. You pass in a directory at program runtime of your images, and it image by image prompts you to type into stdin your desired keywords (space separated). It will not overwrite existing keywords nor add a duplicate keyword if the keyword already exists on the image.

*this gif is out of date in terms of the commands run, but the flow is the same*
![Alt Text](./exifgif.gif)


```
go run main.go ./imgs
```

#### helpful commands:

will show keywords on a file:
```sh
exiftool -Keywords -sep , ./imgs/myimg.jpeg
```
will clear keywords on a file:

```sh
exiftool -overwrite_original -Keywords= ./imgs/myimg.png
```

### coming soon:
- show existing keywords on current image you are editing keywords on
- ability to search for images by keyword


