# Program for rendering static electric field lines formed by two positive charges and two negative charges placed at the vertex of a tetrahedron.

## Example commands - generating image frames
```
go run main.go --out-dir /tmp/tetra
```
## Example command - generate video from image frames
```
ffmpeg -framerate 20.75 -i /tmp/tetra/frame-%04d.png -c:v libx264  -profile:v high -crf 10 -pix_fmt yuv420p -y /tmp/tetra.mp4
```

## Example images

<img width="1280" height="1280" alt="frame-0014" src="https://github.com/user-attachments/assets/6faae15e-8385-4ad1-915f-ff7342bea122" />
<img width="1280" height="1280" alt="frame-0117" src="https://github.com/user-attachments/assets/6f896ea2-8ae4-4aba-943b-7783adfad813" />
