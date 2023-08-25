package service

import (
	"Sala-2/internal/utils"
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/labstack/echo/v4"
)

func UploadImageHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao ler o corpo da solicitação")
	}

	imageData, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao decodificar a imagem")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client := textract.NewFromConfig(cfg)

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return c.String(500, "Erro ao decodificar a imagem")
	}

	result, err := client.AnalyzeDocument(context.TODO(), &textract.AnalyzeDocumentInput{
		Document: &types.Document{
			Bytes: imageData,
		},
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeSignatures,
		},
	})
	if err != nil {
		return c.String(500, err.Error())
	}

	blocks := result.Blocks

	for _, block := range blocks {
		if block.BlockType == "SIGNATURE" {
			err := utils.CropImage(img, block)
			if err != nil {
				return c.String(500, err.Error())
			}

			return c.File("cropped_signature_field.jpg")
		}
	}

	return c.String(400, "Campo de assinatura não encontrado na imagem")
}
