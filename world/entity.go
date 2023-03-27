package world

import "github.com/Lama06/Herder-Games/option"

type Entity struct {
	Position            option.Option[Position]
	Background          option.Option[BackgroundComponent]
	Image               option.Option[ImageComponent]
	KeyboardController  option.Option[KeyboardControllerComponent]
	RectCollider        option.Option[RectColliderComponent]
	ImageBoundsCollider option.Option[ImageBoundsColliderComponent]
	Collissions         option.Option[CollisionsComponent]
	PreventCollisions   option.Option[PreventCollisionsComponent]
}
