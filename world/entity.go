package world

import "github.com/Lama06/Herder-Games/option"

type Entity struct {
	Level    Level
	Position option.Option[Coordinates]

	KeyboardControllerComponent option.Option[KeyboardControllerComponent]

	BackgroundComponent option.Option[BackgroundComponent]
	ImageComponent      option.Option[ImageComponent]
	RectComponent       option.Option[RectComponent]

	RectColliderComponent        option.Option[RectColliderComponent]
	ImageBoundsColliderComponent option.Option[ImageBoundsColliderComponent]

	PreventCollisionsComponent option.Option[PreventCollisionsComponent]

	PortalComponent option.Option[PortalComponent]
}
