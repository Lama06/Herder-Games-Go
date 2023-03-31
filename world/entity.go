package world

import "github.com/Lama06/Herder-Games/option"

type Entity struct {
	Level    Level
	Position option.Option[Coordinates]

	BackgroundComponent option.Option[BackgroundComponent]
	ImageComponent      option.Option[ImageComponent]
	RectComponent       option.Option[RectComponent]

	VelocityComponent option.Option[VelocityComponent]

	KeyboardControllerComponent option.Option[KeyboardControllerComponent]

	RectColliderComponent        option.Option[RectColliderComponent]
	ImageBoundsColliderComponent option.Option[ImageBoundsColliderComponent]

	PortalComponent option.Option[PortalComponent]
}
