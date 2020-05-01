package image

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"image/png"
	"math"
	"os"
	"regexp"
	"strings"
)

const width = 680
const height = 480
const fontTextPath = "./LilitaOne-Regular.ttf"
const fontAuthorPath = "./Lora-Bold.ttf"

type Generator struct {
	context *gg.Context
}

func NewGenerator() IGenerator {
	return &Generator{}
}
func (g *Generator) GenerateBaseImage() error {
	context := gg.NewContext(width, height)
	img, err := gg.LoadImage("./background.jpg")
	if err != nil {
		return errors.New("failed to load image on generator")
	}
	context.DrawImage(img, 0, 0)
	context.SetColor(color.RGBA{0, 0, 0, 100})
	context.DrawRectangle(0, 0, width, height)
	context.Fill()
	g.context = context
	return nil
}

func (g *Generator) SetImage(imagePath string) error {
	img, err := gg.LoadImage(imagePath)
	if err != nil {
		return errors.New("failed to load image on generator")
	}
	g.context.DrawImage(img, width-100, height-100)
	return nil
}

func (g *Generator) SetString(text, fontPath string, size, positionX, positionY, lineSpacing float64, color color.Color) error {
	if err := g.context.LoadFontFace(fontPath, size); err != nil {
		return errors.New("failed to load font")
	}
	g.context.SetColor(color)
	g.context.DrawStringWrapped(text, positionX, positionY, 0, 0, 450, lineSpacing, gg.AlignLeft)
	return nil
}

func (g *Generator) Generate(text, name string) (string, error) {
	err := g.GenerateBaseImage()
	if err != nil {
		return "", err
	}
	nick := `@` + strings.ToUpper(name)
	err = g.SetString(nick, fontAuthorPath, 25, width/2-(math.Ceil(float64(len(nick))/2)*15), height-100, 3, color.Black)
	if err != nil {
		return "", err
	}
	textWithoutLink := removeLink(text)
	if text == "" {
		return "", errors.New("invalid text")
	}
	err = g.SetString(`“`+textWithoutLink+`”`, fontTextPath, 45, 150, 80, 1.5, color.RGBA{244, 228, 33, 255})
	if err != nil {
		return "", err
	}
	err = g.SetString("By: @GetTweetImage", fontAuthorPath, 12, width-150, height-40, 1.5, color.RGBA{0, 0, 0, 200})
	if err != nil {
		return "", err
	}
	fileName := "./tmp/" + getRandomFileName(5) + "_generated.png"
	err = g.SavePNG(fileName)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (g *Generator) SavePNG(filename string) error {
	if err := createTempDir(); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, g.context.Image()); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func createTempDir() error {
	if _, err := os.Stat("./tmp"); err != nil && !os.IsNotExist(err) {
		err := os.MkdirAll("./tmp", os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRandomFileName(n int) (str string) {
	b := make([]byte, n)
	rand.Read(b)
	str = fmt.Sprintf("%x", b)
	return
}

func removeLink(text string) string {
	valid := regexp.MustCompile(`https?://(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)`)
	if valid.MatchString(text) {
		newText := valid.ReplaceAllString(text, "$1W")
		if len(newText) == 0 {
			return newText
		}
		return newText[:len(newText)-1]
	}
	return text
}
