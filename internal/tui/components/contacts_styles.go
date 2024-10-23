package components

import (
	"github.com/charmbracelet/bubbles/list"
)

type contactsStyles struct {
	List     list.Styles
	ListItem list.DefaultItemStyles
}

func enabledContactsStyles() *contactsStyles {
	return &contactsStyles{
		List:     list.DefaultStyles(),
		ListItem: list.NewDefaultItemStyles(),
	}
}

func disabledContactsStyles() *contactsStyles {
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

	return &contactsStyles{
		List:     listStyles,
		ListItem: itemStyles,
	}
}
