package mongu

import "go.mongodb.org/mongo-driver/bson"

//MergeBsonM ignores non bson.M type
func MergeBsonM(ms ...interface{}) bson.M {
	var result = bson.M{}
	for _, one := range ms {
		if one, ok := one.(bson.M); ok {
			for k, v := range one {
				result[k] = v
			}
		}
	}
	return result
}
