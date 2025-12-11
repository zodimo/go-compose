package card

import (
	"go-compose-dev/internal/layoutnode"

	"git.sr.ht/~schnwalter/gio-mw/widget/card"
)

const Material3CardNodeID = "Material3Card"
const Material3CardImageSlotID = "Material3CardImage"

func Elevated(contents CardContentContainer, options ...CardOption) Composable {
	return cardComposable(cardElevated, contents, options...)
}

func Filled(contents CardContentContainer, options ...CardOption) Composable {
	return cardComposable(cardFilled, contents, options...)
}

func Outlined(contents CardContentContainer, options ...CardOption) Composable {
	return cardComposable(cardOutlined, contents, options...)
}

func cardComposable(kind cardKind, contents CardContentContainer, options ...CardOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultCardOptions()
		for _, option := range options {
			option(&opts)
		}

		args := CardConstructorArgs{
			Kind:     kind,
			Contents: contents,
		}

		c.StartBlock(Material3CardNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})

		images := []*GioImage{}
		for _, child := range contents.children {
			if child.contentType == CardContentImage {
				images = append(images, child.image)
				continue
			}
			c.WithComposable(child.composable)
		}

		c.EmitSlot(Material3CardImageSlotID, images)

		c.SetWidgetConstructor(cardWidgetConstructor(args))

		return c.EndBlock()
	}
}

type CardConstructorArgs struct {
	Kind     cardKind
	Contents CardContentContainer
}

func getCardImages(node layoutnode.LayoutNode) []*GioImage {
	return node.FindSlot(Material3CardImageSlotID).UnwrapUnsafe().([]*GioImage)

}

func cardWidgetConstructor(args CardConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			cardWidget := card.Card{Kind: args.Kind}

			cardChildren := []*m3CardChild{}
			nodeChildren := node.Children()
			images := getCardImages(node)

			for _, indexedChild := range args.Contents.children {
				switch indexedChild.contentType {
				case CardContentCover:
					nodeIndex := indexedChild.childIndex
					m3Child := m3CardContentCover(nodeChildren[nodeIndex].(layoutnode.NodeCoordinator).Layout)
					cardChildren = append(cardChildren, m3Child)
				case CardContentImage:
					imageIndex := indexedChild.childIndex
					m3Child := m3CardImage(images[imageIndex])
					cardChildren = append(cardChildren, m3Child)
				default:
					// content
					nodeIndex := indexedChild.childIndex
					m3Child := m3CardContent(nodeChildren[nodeIndex].(layoutnode.NodeCoordinator).Layout)
					cardChildren = append(cardChildren, m3Child)
				}
			}

			return cardWidget.Layout(gtx,
				cardChildren...,
			)
		}
	})
}
