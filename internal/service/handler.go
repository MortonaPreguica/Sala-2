package service

import (
	"Sala-2/internal/utils"
	"bytes"
	"image"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/labstack/echo/v4"
)

func UploadImageHandler(c echo.Context) error {
	ctx := c.Request().Context()
	imageFile, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "Erro ao obter o arquivo da imagem")
	}

	src, err := imageFile.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao abrir o arquivo da imagem")
	}
	defer src.Close()

	imge, _, err := image.Decode(src)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao decodificar a imagem")
	}

	imageData, err := utils.ImageToBytesJPEG(imge)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erro ao ler o arquivo da imagem")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := textract.NewFromConfig(cfg)

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return c.String(500, "Erro ao decodificar a imagem")
	}

	result, err := client.AnalyzeDocument(ctx, &textract.AnalyzeDocumentInput{
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
			img, err := utils.CropImage(img, block)
			if err != nil {
				return c.String(500, err.Error())
			}

			imgBaW := FilterSignature(img)
			imgNoBg := utils.MakeWhiteTransparent(imgBaW)

			imgBytes, err := utils.ImageToBytes(imgNoBg)
			if err != nil {
				return c.String(500, err.Error())
			}

			c.Response().Header().Set(echo.HeaderContentType, "image/png")
			return c.Blob(200, "image/png", imgBytes)
		}
	}

	return c.String(400, "Campo de assinatura não encontrado na imagem")
}
