Whenever I say "assigned an id trivially", I'm referring a specific data type. See #9.

# block_id\<T>
The id system requires a structure of blocks. Each block has a starting index and a length. Ideally stored in order. Blocks do not overlap.
1. `vector<block> blocks`
The id of an added element is then the lowest id that is either higher than the end of all blocks, or lower than their starts. This is necessarily adjacent to some block. When an object is added, the block is extended to include the id, and the id is returned.
2. `sn_id add(T object)`
Objects can be associated with their ids via linear search.
3. `T get_object(sn_id id)`
And ids can be associated with objects via a mapping.
4. `sn_id get_id(T object)`