package components

import (
	"github.com/charmbracelet/bubbles/list"
)

type conversationsStyles struct {
	List     list.Styles
	ListItem list.DefaultItemStyles
}

func enabledConversationsStyles() *conversationsStyles {
	return &conversationsStyles{
		List:     list.DefaultStyles(),
		ListItem: list.NewDefaultItemStyles(),
	}
}

func disabledConversationsStyles() *conversationsStyles {
	listStyles := list.DefaultStyles()
	listStyles.Title = listStyles.Title.UnsetBackground()

	itemStyles := list.NewDefaultItemStyles()

	itemStyles.SelectedTitle = itemStyles.SelectedTitle.
		Foreground(itemStyles.NormalTitle.GetForeground()).
		BorderForeground(itemStyles.NormalTitle.GetForeground())
	itemStyles.SelectedDesc = itemStyles.SelectedDesc.
		Foreground(itemStyles.NormalTitle.GetForeground()).
		BorderForeground(itemStyles.NormalTitle.GetForeground())

	itemStyles.NormalTitle = itemStyles.DimmedTitle
	itemStyles.NormalDesc = itemStyles.DimmedDesc

	return &conversationsStyles{
		List:     listStyles,
		ListItem: itemStyles,
	}
}
