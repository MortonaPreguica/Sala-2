package rest

import (
	"Sala-2/internal/environment"
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/teste", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/upload", uploadImageHandler)
}

func uploadImageHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao ler o corpo da solicitação")
	}

	imageData, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao decodificar a imagem")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(environment.Env.Region),
	})
	if err != nil {
		panic(err)
	}

	svc := rekognition.New(sess)

	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: imageData,
		},
	}

	result, err := svc.DetectText(input)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao chamar o Amazon Rekognition")
	}

	// Process the result and extract the detected text
	var extractedText []string
	for _, item := range result.TextDetections {
		extractedText = append(extractedText, *item.DetectedText)
	}
	var confidence []float64
	for _, item := range result.TextDetections {
		confidence = append(confidence, *item.Confidence)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"extracted_text": extractedText,
		"confidence":     confidence,
	})
}
