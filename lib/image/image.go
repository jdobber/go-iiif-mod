package image

type Image interface {
	Identifier() string
	Transform(*Transformation) error // http://iiif.io/api/image/2.1/#order-of-implementation
	//Format() string
	Height() int
	Width() int
}
