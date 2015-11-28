### AC41001 - Assignment 2
#### Yago Carballo

##### Description

The OpenGL app shows a terrain with transparent moving water and fish as well as some loaded Objects which include a Gingerbread House, a Car, a Wall, a Dragon, two types of fish and a blue gopher.

The terrain and the water are generated using the example from the labs adapted to go and instead of perling noise is using OpenSimpleX Noise. The terrain has a fixed seed and the water has a changing seed (which is not random, is incremental).

All the objects are loaded from Wavefront's `.obj` and `.mtl` files. The parser was initially adapted from the VU project (which handle's some simple .obj cases) and then extended so it could work with a few more `.obj` files as I was finding them on the internet. (.obj files might have more than one object as for example the car has 52 objects and 43 materials)

##### Shaders

- **colorMaterial:** The Gopher and the Car don't have textures so they load the colors given in the material file referenced in the `.obj` file, if the material is transparent, it enables transparencies to the object.
- **textureMaterial:** The Dragon has a texture so is loading the texture given in the material file and using it (colors from the material are loaded and given to the shader but the shader is not using them).
- **bumpMapMaterial:** The rest of the loaded models have Bump Maping and they load the bump map texture from the material file.
- **terrain:** The colors of the terrain are generated with the terrain in black and white (to to give a more realistic tone), then on the shader a single tone is applied to give color to the terrain.
- **basic:** Implementation of the phong lighting, used only on the light sphere.


##### Problems

Almost every .obj I found was different, so I tried to follow the specification for some things. When trying to load objects with more than 3 faces, some models had holes on them, so only triangulated objects work well. (but if is triangulated and exported from blender it will probably work)

Almost at the end of the assignment I found some tree models that had a different material for different groups of faces. Those models work, but they use the last material found for the whole object.

<br />
##### References
- Most of the 3D objects come from [http://www.reinerstilesets.de/3d-grafiken/](http://www.reinerstilesets.de/3d-grafiken/)
- The dragon I used is from an old project I had [https://warofsides.com](https://warofsides.com)
- The Gopher model [Gopher 3D (https://github.com/golang-samples/gopher-3d)](https://github.com/golang-samples/gopher-3d)
- The [VU Engine (https://github.com/gazed/vu)](https://github.com/gazed/vu) was used as reference to see how they're doing openGL in Go.

##### Instructions
Instructions on how to run the app are inside the `readme.md` file (method is a bit different than the first assignment, (as this time I used the "officially recommended" approach to organize the code which is suposed to be easier))
> Source code will be in `src/github.com/yagocarballo/Go-GL-Assignment-2`

##### Controls
To control the different objects, you can press one of the numbers and that will select the object (selected object's name will be displayed in the title of the window), then that object can be moved, rotated, scaled or changed the drawing mode (dots, lines or 3D).

The following instructions will be printed to the terminal when opening the App:
```
Keyboard Instructions
---------------------
Select A Model with a number
- The name of the selected model is displayed on the title of the window.
- Once a model is selected, it can be moved using the Keys Below
- Move and Rotate work like a cross: ( Example )

| W   \ S 	-->  Up  \ Down	 |
| A   \ D 	--> Left \ Right |
| Q   \ E 	--> Back \ Front |
| Tab \ R 	-->   -W \ +W	 |

 |-----||-----||-----||-----||-----|
 |  1  ||  2  ||  3  ||  4  ||  5  |
 |-----||-----||-----||-----||-----|
 |-----||-----||-----||-----||-----|
 |  6  ||  7  ||  8  ||  9  ||  0  |
 |-----||-----||-----||-----||-----|

1. Gopher				| 6.
2. Gingerbread House	| 7.
3. Dragon				| 8. Terrain
4. Car					| 9. Camera / View
5. Wall					| 0. Light Point

		  Position (Move)
----------------------------

|-------||-----||-----||-----||-----|
|  Tab  ||  Q  ||  W  ||  E  ||  R  |
|-------||-----||-----||-----||-----|
		 |-----||-----||-----|
		 |  A  ||  S  ||  D  |
		 |-----||-----||-----|

		 Rotation (Rotate)
----------------------------
|-----||-----||-----||-----||-----|
|  Y  ||  U  ||  I  ||  O  ||  P  |
|-----||-----||-----||-----||-----|
 	   |-----||-----||-----|
 	   |  J  ||  K  ||  L  |
 	   |-----||-----||-----|


 Zoom (-/+)  | Speed Up/Down
----------------------------
|-----||-----||-----||-----|
|  Z  ||  X  ||  Z  ||  X  |
|-----||-----||-----||-----|



Instructions / Draw Mode / Color
	---------------------
	|-----||-----||-----|
	|  B  ||  N  ||  M  |
	|-----||-----||-----|

----------------------------

DEBUG Options

- The Enter Key will reload the Shaders.
- The Space Key will regenerate the water

```
