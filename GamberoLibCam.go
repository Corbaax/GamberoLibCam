package GamberoLibCam

/**
*	@Author: Luca Corbetta
*	@Ver: 01.A
*	@Date: 17.12.2021
*	@Description: The Object will help you displaying bigger image than the screen size by rendering a portion of the image. it's possible to surfing the image.
**/
import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

//#Region CamController
type Ground struct {
	Background, Cam              *ebiten.Image
	CamX, CamY, CamDimX, CamDimY int
	BackBoundX, BackBoundY       int
	Translation                  int
	ErrorOnCam                   string
}

type MovingOption struct {
	Up, Down, Right, Left ebiten.Key
}

// Populate the Moving option with default value WASD
func (NewObj *MovingOption) NewDefaultOption() {
	NewObj.down = ebiten.KeyS
	NewObj.up = ebiten.KeyW
	NewObj.right = ebiten.KeyD
	NewObj.left = ebiten.KeyA
}

// Populate the Moving option with the value in the slice, taking the first 4 value
// assing as: 0 = down, 1 = up, 2 = right; 3 = left
func (NewObj *MovingOption) NewOptionFromSlice(keys []ebiten.Key) {
	NewObj.down = keys[0]
	NewObj.up = keys[1]
	NewObj.right = keys[2]
	NewObj.left = keys[3]
}

// Populate the Moving option with the value given int the func
func (NewObj *MovingOption) NewOptionFromDirect(Down, Up, Right, Left ebiten.Key) {
	NewObj.down = Down
	NewObj.up = Up
	NewObj.right = Right
	NewObj.left = Left
}

// The func will rendere the display dimension if it's inside the Ground Bound
func (Display *Ground) RenderDisplay() {
	Display.cam = Display.background.SubImage(image.Rect(Display.camX, Display.camY, Display.camX+Display.camDimX, Display.camY+Display.camDimY)).(*ebiten.Image)
}

// The func will check if the Display will overflow from the border: owf = flase, inbound = true
func isInBound(x, y, sizeCamX, sizeCamY, BackX, BackY int) bool {
	if x+sizeCamX <= BackX && y+sizeCamY <= BackY && x >= 0 && y >= 0 {
		return true
	} else {
		return false
	}
}

// Set the new Ground From an Image, default value of cam 400x400 @ 10 Unit
func (Display *Ground) NewDefaultGroundFromImage(BackGround *ebiten.Image) {
	Display.background = BackGround
	Display.BackBoundX = Display.background.Bounds().Dx()
	Display.BackBoundY = Display.background.Bounds().Dy()
	Display.camDimX = 400
	Display.camDimY = 400
	Display.camX = 0
	Display.camY = 0
	Display.Translation = 10
	Display.RenderDisplay()
}

// Set the new Ground From a Path, default value of cam 400x400 @ 10 Unit
func (Display *Ground) NewDefaultGroundFromPath(Path string) {
	Display.background = LoadImage(Path)
	Display.BackBoundX = Display.background.Bounds().Dx()
	Display.BackBoundY = Display.background.Bounds().Dy()
	Display.camDimX = 400
	Display.camDimY = 400
	Display.camX = 0
	Display.camY = 0
	Display.Translation = 10
	Display.RenderDisplay()
}

// Set the new Ground From an Image and set the dimension of the cam as 'DimX'x'DimY' @ 10 Unit
func (Display *Ground) NewGroundWithDim(BackGround *ebiten.Image, dimX, dimY int) {
	Display.background = BackGround
	Display.BackBoundX = Display.background.Bounds().Dx()
	Display.BackBoundY = Display.background.Bounds().Dy()
	Display.camDimX = dimX
	Display.camDimY = dimY
	Display.camX = 0
	Display.camY = 0
	Display.Translation = 10
	Display.RenderDisplay()
}

// Set the Cam speed as 'speedUnit'
func (Display *Ground) SetCamSpeed(speedUnit int) {
	if speedUnit > 0 {
		Display.Translation = speedUnit
	} else {
		Display.ErrorOnCam = "Speed_Can_Not_Be_Negative_Or_0"
	}
}

func (Display *Ground) SetNewCoordinate(x, y int) {
	if x >= 0 && y >= 0 {
		Display.camX = x
		Display.camY = y
	} else {
		Display.ErrorOnCam = "x_and_y_must_Be_Positive"
	}
}

// Resize the Cam Size as 'dimX' x 'dimY'
func (Display *Ground) ResizeCam(dimX, dimY int) {
	Display.SetNewCoordinate(0, 0)
	Display.camDimX = dimX
	Display.camDimY = dimY
}

// Return the Center of the Cam as x , y int
func (Display *Ground) GetCenterPosition() (int, int) {
	return Display.camDimX / 2, Display.camDimY / 2
}

// if the cam is touching one the border will return true.
// Bool Order Up - Down - Right - Left
func (Display *Ground) IsTouchingBorder() (bool, bool, bool, bool) {
	var up, down, right, left bool = false, false, false, false
	if Display.camX+Display.camDimX >= Display.BackBoundX {
		right = true
	}
	if Display.camY+Display.camDimY >= Display.BackBoundY {
		up = true
	}
	if Display.camX <= 0 {
		left = true
	}
	if Display.camY <= 0 {
		down = true
	}
	return up, down, right, left
}

// Check the size and move the cam of 10
func (Display *Ground) MovingCam(move MovingOption) {
	if ebiten.IsKeyPressed(move.up) && isInBound(Display.camX, Display.camY-Display.Translation, Display.camDimX, Display.camDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.camY -= Display.Translation
	}
	if ebiten.IsKeyPressed(move.left) && isInBound(Display.camX-Display.Translation, Display.camY, Display.camDimX, Display.camDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.camX -= Display.Translation
	}
	if ebiten.IsKeyPressed(move.down) && isInBound(Display.camX, Display.camY+Display.Translation, Display.camDimX, Display.camDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.camY += Display.Translation
	}
	if ebiten.IsKeyPressed(move.right) && isInBound(Display.camX+Display.Translation, Display.camY, Display.camDimX, Display.camDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.camX += Display.Translation
	}

	Display.RenderDisplay()
}

// LoadImage will load an image from a string Path
func LoadImage(Path string) *ebiten.Image {

	infile, err := os.Open(Path)
	if err != nil {
		panic(err)
	}
	defer infile.Close()
	src, _, err := image.Decode(infile)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(src)
}

//#EndRegion
