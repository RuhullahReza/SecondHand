package repository

import (
	"fmt"
	"log"
	"strings"

	"github.com/RuhullahReza/SecondHand/helper"
	"github.com/RuhullahReza/SecondHand/model/entity"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

type ImageRepository interface {
	Upload(ctx context.Context, input interface{}, folder string) (string, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, image *entity.Image) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	GetByProductId(ctx context.Context, id uuid.UUID) ([]entity.ProductImage, error) 
	GetPathById(ctx context.Context, id uuid.UUID) (*entity.Image, error)
	GetPublicId(url string) string
}

type ImageRepositoryImpl struct {
	CLD *cloudinary.Cloudinary
	DB *sqlx.DB
}

func NewImageRepository(cld *cloudinary.Cloudinary, db *sqlx.DB) ImageRepository {
	return &ImageRepositoryImpl{
		CLD: cld,
		DB: db,
	}
}

func (r *ImageRepositoryImpl) Upload(ctx context.Context, input interface{}, folder string) (string, error) {

	uploadParam, err := r.CLD.Upload.Upload(ctx, input, uploader.UploadParams{Folder: folder})
	if err != nil {
		log.Printf("failed to upload image, err : %v\n", err)
		return "", helper.NewInternal()
	}
	
	return uploadParam.SecureURL, nil
}

func (r *ImageRepositoryImpl) Delete(ctx context.Context, id string) error {

	publicId := fmt.Sprintf("secondHand-go/%s", r.GetPublicId(id))
	_, err := r.CLD.Upload.Destroy(ctx, uploader.DestroyParams {PublicID: publicId})
	if err != nil {
		log.Printf("failed to delete image, err : %v\n", err)
		return helper.NewInternal()
	}
	
	return nil
}

func (r *ImageRepositoryImpl) GetPublicId(url string) string {

	urlSplitted := strings.Split(url,"/")
	folder := urlSplitted[len(urlSplitted)-2]
	id := strings.Split(urlSplitted[len(urlSplitted)-1], ".")[0]

	return folder + "/" + id
}

func (r *ImageRepositoryImpl) Create(ctx context.Context, image *entity.Image) error {

	query := `
	INSERT INTO 
		images 
		(product_id, image_url) 
	VALUES 
		($1, $2)
	`
	_, err := r.DB.ExecContext(ctx, query, image.ProductId, image.Url)

	if err != nil {
		log.Printf("failed to query create image, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil
}

func (r *ImageRepositoryImpl) DeleteById(ctx context.Context, id uuid.UUID) error {

	query := `
	DELETE FROM
		images
	WHERE 
		id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, id)

	if err != nil {
		log.Printf("failed to query delete image by id, err : %v\n", err)
		return helper.NewInternal()
	}

	return nil 
}

func (r *ImageRepositoryImpl) GetPathById(ctx context.Context, id uuid.UUID) (*entity.Image, error) {
	
	image := &entity.Image{}

	query := `
	SELECT 
		image_url
	FROM 
		images 
	WHERE 
		id=$1
	LIMIT 1
	`

	if err := r.DB.GetContext(ctx, image, query, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return image, helper.NewNotFound("id", id.String())
		}
		
		log.Printf("failed to query get image by Id, err : %v\n", err)
		return image, helper.NewInternal()
	}

	return image, nil
}


func (r *ImageRepositoryImpl) GetByProductId(ctx context.Context, id uuid.UUID) ([]entity.ProductImage, error) {

	images := []entity.ProductImage{}

	query := `
		SELECT 
			id, image_url
		FROM
			images
		WHERE 
			product_id = $1
	`
	rows, err := r.DB.QueryContext(ctx,query, id)
	if err != nil {
		log.Printf("failed to query get all product, err : %v\n", err)
		return images, helper.NewInternal()
	}

	for rows.Next(){
		image := entity.ProductImage{}
		err := rows.Scan(&image.Id, &image.Url)
		if err != nil {
			log.Printf("failed to scanning image, err : %v\n", err)
			return images, helper.NewInternal()
		}

		images = append(images, image)
	}
	
	return images, nil
}