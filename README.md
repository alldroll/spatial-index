# Spatial index (for fun)

Try to play with Google maps, quadkeys, quadtree, points clustering.

So the idea is very simple:
We build a simple trie, where key is quadkey.
Each node represents an appropriate tile and zoom.
Points are contained in the leaves of the trie.
For clustering, each node (non leaf) is a cluster. It is all:)

## Demo

https://shrouded-cove-41776.herokuapp.com/

