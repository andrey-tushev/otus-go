package hw04lrucache

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
	front, back *ListItem
	len         int
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.front,
	}

	if l.front != nil {
		l.front.Prev = newItem
	}
	l.front = newItem

	if l.back == nil {
		l.back = newItem
	}

	l.len++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  l.back,
		Next:  nil,
	}

	if l.back != nil {
		l.back.Next = newItem
	}
	l.back = newItem

	if l.front == nil {
		l.front = newItem
	}

	l.len++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	next := i.Next
	prev := i.Prev

	if prev != nil {
		prev.Next = next
	} else {
		l.front = next
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.back = prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
