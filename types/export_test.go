package types

func GetKVMapFromItemList(il *ItemList) map[string]*Item {
	if il == nil {
		return nil
	}

	return il.kvMap
}
