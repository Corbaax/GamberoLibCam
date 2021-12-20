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
	CamX, CamY, camDimX, CamDimY int
	BackBoundX, BackBoundY       int
	Translation                  int
	ErrorOnCam                   string
}

type MovingOption struct {
	Up, Down, Right, Left ebiten.Key
}

// Populate the Moving option with default value WASD
func (NewObj *MovingOption) NewDefaultOption() {
	NewObj.Down = ebiten.KeyS
	NewObj.Up = ebiten.KeyW
	NewObj.Right = ebiten.KeyD
	NewObj.Left = ebiten.KeyA
}

// Populate the Moving option with the value in the slice, taking the first 4 value
// assing as: 0 = down, 1 = up, 2 = Right; 3 = left
func (NewObj *MovingOption) NewOptionFromSlice(keys []ebiten.Key) {
	NewObj.Down = keys[0]
	NewObj.Up = keys[1]
	NewObj.Right = keys[2]
	NewObj.Left = keys[3]
}

// Populate the Moving option with the value given int the func
func (NewObj *MovingOption) NewOptionFromDirect(Down, Up, Right, Left ebiten.Key) {
	NewObj.Down = Down
	NewObj.Up = Up
	NewObj.Right = Right
	NewObj.Left = Left
}

// The func will rendere the display dimension if it's inside the Ground Bound
func (Display *Ground) RenderDisplay() {
	Display.Cam = Display.Background.SubImage(image.Rect(Display.CamX, Display.CamY, Display.CamX+Display.camDimX, Display.CamY+Display.CamDimY)).(*ebiten.Image)
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
	Display.Background = BackGround
	Display.BackBoundX = Display.Background.Bounds().Dx()
	Display.BackBoundY = Display.Background.Bounds().Dy()
	Display.camDimX = 400
	Display.CamDimY = 400
	Display.CamX = 0
	Display.CamY = 0
	Display.Translation = 10
	Display.RenderDisplay()
}

// Set the new Ground From a Path, default value of cam 400x400 @ 10 Unit
func (Display *Ground) NewDefaultGroundFromPath(Path string) {
	Display.Background = LoadImage(Path)
	Display.BackBoundX = Display.Background.Bounds().Dx()
	Display.BackBoundY = Display.Background.Bounds().Dy()
	Display.camDimX = 400
	Display.CamDimY = 400
	Display.CamX = 0
	Display.CamY = 0
	Display.Translation = 10
	Display.RenderDisplay()
}

// Set the new Ground From an Image and set the dimension of the cam as 'DimX'x'DimY' @ 10 Unit
func (Display *Ground) NewGroundWithDim(BackGround *ebiten.Image, dimX, dimY int) {
	Display.Background = BackGround
	Display.BackBoundX = Display.Background.Bounds().Dx()
	Display.BackBoundY = Display.Background.Bounds().Dy()
	Display.camDimX = dimX
	Display.CamDimY = dimY
	Display.CamX = 0
	Display.CamY = 0
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
		Display.CamX = x
		Display.CamY = y
	} else {
		Display.ErrorOnCam = "x_and_y_must_Be_Positive"
	}
}

// Resize the Cam Size as 'dimX' x 'dimY'
func (Display *Ground) ResizeCam(dimX, dimY int) {
	Display.SetNewCoordinate(0, 0)
	Display.camDimX = dimX
	Display.CamDimY = dimY
}

// Return the Center of the Cam as x , y int
func (Display *Ground) GetCenterPosition() (int, int) {
	return Display.camDimX / 2, Display.CamDimY / 2
}

// if the cam is touching one the border will return true.
// Bool Order Up - Down - Right - Left
func (Display *Ground) isTouchingBorder() (bool, bool, bool, bool) {
	var up, down, Right, left bool = false, false, false, false
	if Display.CamX+Display.camDimX >= Display.BackBoundX {
		Right = true
	}
	if Display.CamY+Display.CamDimY >= Display.BackBoundY {
		up = true
	}
	if Display.CamX <= 0 {
		left = true
	}
	if Display.CamY <= 0 {
		down = true
	}
	return up, down, Right, left
}

// Check the size and move the cam of 10
func (Display *Ground) MovingCam(move MovingOption) {
	if ebiten.IsKeyPressed(move.Up) && isInBound(Display.CamX, Display.CamY-Display.Translation, Display.camDimX, Display.CamDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.CamY -= Display.Translation
	}
	if ebiten.IsKeyPressed(move.Left) && isInBound(Display.CamX-Display.Translation, Display.CamY, Display.camDimX, Display.CamDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.CamX -= Display.Translation
	}
	if ebiten.IsKeyPressed(move.Down) && isInBound(Display.CamX, Display.CamY+Display.Translation, Display.camDimX, Display.CamDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.CamY += Display.Translation
	}
	if ebiten.IsKeyPressed(move.Right) && isInBound(Display.CamX+Display.Translation, Display.CamY, Display.camDimX, Display.CamDimY, Display.BackBoundX, Display.BackBoundY) {
		Display.CamX += Display.Translation
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
