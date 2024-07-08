package models

import (
	productErrors "github.com/AleksK1NG/products-microservice/pkg/product_errors"
	"github.com/AleksK1NG/products-microservice/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Product models
type Product struct {
	ProductID   primitive.ObjectID `json:"productId" bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID `json:"categoryId,omitempty" bson:"categoryId,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	ImageURL    *string            `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	Photos      []string           `json:"photos,omitempty" bson:"photos,omitempty"`
	Quantity    int64              `json:"quantity,omitempty" bson:"quantity,omitempty" validate:"required"`
	Rating      int                `json:"rating,omitempty" bson:"rating,omitempty" validate:"required,min=0,max=10"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type ProductsList struct {
	TotalCount int64      `json:"totalCount"`
	TotalPages int64      `json:"totalPages"`
	Page       int64      `json:"page"`
	Size       int64      `json:"size"`
	HasMore    bool       `json:"hasMore"`
	Products   []*Product `json:"products"`
}

// Create new product
func (p *postgresRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresRepo.Create")
	defer span.Finish()

	collection := p.dbpool.Database(productsDB).Collection(productsCollection)

	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()

	result, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "InsertOne")
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.Wrap(productErrors.ErrObjectIDTypeConversion, "result.InsertedID")
	}

	product.ProductID = objectID

	return product, nil
}

// Update Single product
func (p *postgresRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresRepo.Update")
	defer span.Finish()

	collection := p.dbpool.Database(productsDB).Collection(productsCollection)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var prod models.Product
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": product.ProductID}, bson.M{"$set": product}, ops).Decode(&prod); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &prod, nil
}

// GetByID Get single product by id
func (p *postgresRepo) GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresRepo.GetByID")
	defer span.Finish()

	collection := p.dbpool.Database(productsDB).Collection(productsCollection)

	var prod models.Product
	if err := collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&prod); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	return &prod, nil
}

// Search Search product
func (p *postgresRepo) Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresRepo.Search")
	defer span.Finish()

	collection := p.dbpool.Database(productsDB).Collection(productsCollection)

	f := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{
				Pattern: search,
				Options: "gi",
			}}},
			bson.D{{Key: "description", Value: primitive.Regex{
				Pattern: search,
				Options: "gi",
			}}},
		}},
	}

	count, err := collection.CountDocuments(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &models.ProductsList{
			TotalCount: 0,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Products:   make([]*models.Product, 0),
		}, nil
	}

	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, f, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx)

	products := make([]*models.Product, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var prod models.Product
		if err := cursor.Decode(&prod); err != nil {
			return nil, errors.Wrap(err, "Find")
		}
		products = append(products, &prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return &models.ProductsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Products:   products,
	}, nil
}
