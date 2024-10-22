package components

import (
	"github.com/charmbracelet/bubbles/list"
)

type ContactsStyles struct {
	List     list.Styles
	ListItem list.DefaultItemStyles
}

func EnabledContactsStyles() *ContactsStyles {
	return &ContactsStyles{
		List:     list.DefaultStyles(),
		ListItem: list.NewDefaultItemStyles(),
	}
}

func DisabledContactsStyles() *ContactsStyles {
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

	return &ContactsStyles{
		List:     listStyles,
		ListItem: itemStyles,
	}
}
