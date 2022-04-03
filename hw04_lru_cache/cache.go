package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	//Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

// Элемент кэша хранит в себе ключ, по которому он лежит в словаре, и само значение.
type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

//если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди;

//если элемента нет в словаре, то добавить в словарь и в начало очереди (при этом, если размер очереди больше ёмкости кэша,
//то необходимо удалить последний элемент из очереди и его значение из словаря);

//возвращаемое значение - флаг, присутствовал ли элемент в кэше.

func (l lruCache) Set(key Key, value interface{}) bool {
	existedItem, exists := l.items[key]
	_ = existedItem
	if exists {
		// Замена значения на новое
		e := existedItem.Value.(*cacheItem)
		e.value = value
	} else {
		// Создадим новое значение
		newItem := &cacheItem{
			key:   key,
			value: value,
		}
		l.items[key] = l.queue.PushFront(newItem)

		// Удаление самого редко используемого элемента
		if l.queue.Len() > l.capacity {
			oldest := l.queue.Back()
			delete(l.items, oldest.Value.(*cacheItem).key)
			l.queue.Remove(oldest)
		}
	}

	return exists
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	existedItem, exists := l.items[key]
	if !exists {
		return nil, false
	}

	e := existedItem.Value.(*cacheItem)
	l.queue.PushFront(existedItem)
	return e.value, exists
}

func (l lruCache) Clear() {
	panic("implement me")
}
