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
<img width="1280" height="1280" alt="frame-0041" src="https://github.com/user-attachments/assets/534b2f48-431d-4d0e-bf92-4ce8e9572046" />
<img width="1280" height="1280" alt="frame-0149" src="https://github.com/user-attachments/assets/64ad3f77-c75d-4d88-87b8-740e9528a1ef" />
