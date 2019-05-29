package bricklinkstore

import "fmt"

// GetColors retrieves a list of the colors defined within BrickLink catalog.
func (c *Client) GetColors() ([]Color, error) {
	url := "/colors"
	var colors colorsResponse
	if err := c.doGet(url, &colors); err != nil {
		return nil, err
	}
	return colors.Data, checkMeta(colors.Meta)
}

type colorsResponse struct {
	Meta meta    `json:"meta"`
	Data []Color `json:"data"`
}

// GetColor retrieves information about a specific color.
func (c *Client) GetColor(id int) (*Color, error) {
	url := fmt.Sprintf("/colors/%d", id)
	var color colorResponse
	if err := c.doGet(url, &color); err != nil {
		return nil, err
	}
	return &color.Data, checkMeta(color.Meta)
}

type colorResponse struct {
	Meta meta  `json:"meta"`
	Data Color `json:"data"`
}

// Color contains information about a color in the BrickLink catalog.
type Color struct {
	ColorID   int       `json:"color_id"`   // ID of the color
	ColorName string    `json:"color_name"` // The name of the color
	ColorCode string    `json:"color_code"` // HTML color code of this color
	ColorType ColorType `json:"color_type"` // The name of the color group to which this color belongs
}

type ColorType string

const (
	ColorTypeBrickArms   ColorType = "BrickArms"
	ColorTypeChrome      ColorType = "Chrome"
	ColorTypeGlitter     ColorType = "Glitter"
	ColorTypeMetallic    ColorType = "Metallic"
	ColorTypeMilky       ColorType = "Milky"
	ColorTypeModulex     ColorType = "Modulex"
	ColorTypePearl       ColorType = "Pearl"
	ColorTypeSolid       ColorType = "Solid"
	ColorTypeSpeckle     ColorType = "Speckle"
	ColorTypeTransparent ColorType = "Transparent"
)
