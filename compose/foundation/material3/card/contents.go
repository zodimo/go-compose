package card

type CardContentContainer struct {
	composableCount int
	imageCount      int
	children        []*indexedCardChild
}

func newCardContentContainer() *CardContentContainer {
	return &CardContentContainer{
		composableCount: 0,
		imageCount:      0,
		children:        []*indexedCardChild{},
	}
}

func (c *CardContentContainer) addChild(index int, child *CardChild) {
	var childIndex int
	if child.contentType == CardContentImage {
		childIndex = c.imageCount
		c.imageCount++
	} else {
		childIndex = c.composableCount
		c.composableCount++
	}

	c.children = append(c.children, &indexedCardChild{
		index:       index,
		contentType: child.contentType,
		childIndex:  childIndex,
		image:       child.image,
		composable:  child.composable,
	})
}

func CardContents(children ...*CardChild) CardContentContainer {

	container := newCardContentContainer()
	for index, child := range children {
		container.addChild(index, child)
	}
	return *container
}

func Image(image *GioImage) *CardChild {
	return &CardChild{
		image:       image,
		contentType: CardContentImage,
	}
}

func Content(composable Composable) *CardChild {
	return &CardChild{
		composable:  composable,
		contentType: CardContent,
	}
}

func ContentCover(composable Composable) *CardChild {
	return &CardChild{
		cover:       true,
		composable:  composable,
		contentType: CardContentCover,
	}
}

type indexedCardChild struct {
	index       int
	contentType CardContentType
	childIndex  int
	image       *GioImage
	composable  Composable
}

type CardChild struct {
	cover       bool
	image       *GioImage
	composable  Composable
	contentType CardContentType
}
