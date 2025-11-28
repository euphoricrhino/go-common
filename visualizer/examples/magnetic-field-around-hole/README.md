# Program for rendering magnetic field line near planar conductor with hole

## Example commands - generating image frames
```
go run main.go --out-dir /tmp/hole
```
## Example command - generate video from image frames
```
ffmpeg -framerate 20.75 -i /tmp/hole/frame-%04d.png -c:v libx264  -profile:v high -crf 10 -pix_fmt yuv420p -y /tmp/hole-test.mp4
```

## Example images
