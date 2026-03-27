package domain

type Parcel struct {
	ID           int64
	ProducerID   int64
	Name         string
	Description  string
	GeometryWKT  string
	AreaHectares float64
	CropType     string
	CreatedAt    int64
}

type CreateParcelRequest struct {
	ProducerID   int64
	Name         string
	Description  string
	GeometryWKT  string
	AreaHectares float64
	CropType     string
}
