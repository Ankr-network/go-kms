package kvdb

/*
kvdb api constraint: we store anything, the data must format map, thank you know.
we thought to replace map structure with []byte structure, but for extensions, we
think the map is well,the map is a stable structure. think about some scenario of
only take part of the secrets, we can easy handle with it.
*/
type KVer interface {
	// take secrets by specify path
	Get(path string) (map[string]interface{}, error)
	// put secrets by specify path
	Put(path string, value map[string]string) error
}
