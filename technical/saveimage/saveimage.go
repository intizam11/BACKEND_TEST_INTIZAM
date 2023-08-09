package saveimage

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveImage(imageFile multipart.File, imageName string) error {
	destinationPath := `C:\Users\LENOVO\Desktop\technical\image\` + imageName
	outputImage, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer outputImage.Close()

	_, err = io.Copy(outputImage, imageFile)
	if err != nil {
		return err
	}

	return nil
}

// destinationPatch := `C:\Users\LENOVO\Desktop\technical\image` + imageName
