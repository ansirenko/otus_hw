package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	count int
	first *ListItem
	last *ListItem
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFront := ListItem{Value: v}
	if l.count == 0 {
		l.last = &newFront
	} else {
		l.first.Prev = &newFront
		newFront.Next = l.first
	}
	l.first = &newFront
	l.count++
	return &newFront
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := ListItem{Value: v}
	if l.count == 0 {
		l.first = &newBack
	} else {
		l.last.Next = &newBack
		newBack.Prev = l.last
	}
	l.last = &newBack
	l.count++
	return &newBack
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.last == i && l.first == i:
		l.last = nil
		l.first = nil
	case l.last == i:
		l.last = i.Prev
		l.last.Next = nil
	case l.first == i:
		l.first = i.Next
		l.first.Prev = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.first == i {
		return
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
	}
	i.Prev.Next = i.Next
	i.Next = l.first
	l.first.Prev = i
	l.first = i
}

func NewList() List {
	return &list{}
}
