package world

import "github.com/Lama06/Herder-Games/option"

type Entity struct {
	Position           option.Option[WorldPosition]
	BackgroundRenderer option.Option[BackgroundRendererComponent]
	Renderer           option.Option[RendererComponent]
	RectRenderer       option.Option[RectRendererComponent]
	ImageRenderer      option.Option[ImageRendererComponent]
	KeyboardController option.Option[KeyboardControllerComponent]
	RectCollider       option.Option[RectColliderComponent]
	Collissions        option.Option[CollisionsComponent]
	PreventCollisions  option.Option[PreventCollisionsComponent]
}
