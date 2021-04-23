module Gee

go 1.15

require (
	gee v0.0.0
	geeorm v0.0.0
	geecache v0.0.0
	geerpc v0.0.0
)

replace (
	gee => ./gee
	geeorm => ./geeorm
	geecache => ./geecache
	geerpc => ./geerpc
)
